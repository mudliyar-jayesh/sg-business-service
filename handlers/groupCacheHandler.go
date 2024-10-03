package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
)

// Group represents the model equivalent to the C# Group model
type Group struct {
	GUID   string
	Name   string
	Parent string
}

// GroupNode represents the structure of a group node in the tree
type GroupNode struct {
	Parent   *GroupNode
	Value    string
	Children []*GroupNode
}

// GroupCacheManager manages the group cache by company ID
type GroupCacheManager struct {
	rootByCompanyId map[string]*GroupNode
}

// NewGroupCacheManager creates a new instance of GroupCacheManager
func NewGroupCacheManager() *GroupCacheManager {
	return &GroupCacheManager{
		rootByCompanyId: make(map[string]*GroupNode),
	}
}

var CachedGroups *GroupCacheManager

func MakeGroupCache() {
	CachedGroups = NewGroupCacheManager()
	CachedGroups.BuildCache()
}

// FindNode recursively finds a node by its value
func (gcm *GroupCacheManager) FindNode(node *GroupNode, value string) *GroupNode {
	if node.Value == value {
		return node
	}

	for _, child := range node.Children {
		found := gcm.FindNode(child, value)
		if found != nil {
			return found
		}
	}
	return nil
}

// SearchChildrenByParent recursively collects all child values of a given node
func (gcm *GroupCacheManager) SearchChildrenByParent(node *GroupNode) []string {
	children := []string{}
	for _, child := range node.Children {
		children = append(children, child.Value)
		children = append(children, gcm.SearchChildrenByParent(child)...)
	}
	return children
}

// Build constructs the group tree for a given company ID
func (gcm *GroupCacheManager) Build(companyId string, groups []Group) {
	nodeByGroupName := make(map[string]*GroupNode)

	// Initialize the map with group names
	for _, group := range groups {
		nodeByGroupName[group.Name] = &GroupNode{Value: group.Name}
	}

	// Create the root node
	root := &GroupNode{Value: "Primary"}

	// Build the tree
	for _, group := range groups {
		parent := root
		if !strings.EqualFold(group.Parent, "Primary") {
			parent = nodeByGroupName[group.Parent]
		}

		child := nodeByGroupName[group.Name]
		child.Parent = parent

		// Check if the child already exists in the parent's children
		exists := false
		for _, existingChild := range parent.Children {
			if existingChild.Value == child.Value {
				exists = true
				break
			}
		}
		if !exists {
			parent.Children = append(parent.Children, child)
		}
	}

	// Store the root node in the map by company ID
	gcm.rootByCompanyId[companyId] = root
}

// GetChildrenNames retrieves the names of the children for a given parent node in a company
func (gcm *GroupCacheManager) GetChildrenNames(companyId, parent string) []string {
	root, exists := gcm.rootByCompanyId[companyId]
	if !exists || root == nil {
		return []string{}
	}

	node := gcm.FindNode(root, parent)
	if node == nil {
		return []string{}
	}

	parentGroups := []string{parent}
	parentGroups = append(parentGroups, gcm.SearchChildrenByParent(node)...)
	return parentGroups
}

func (gcm *GroupCacheManager) BuildCache() {
	collection := GetCollection("NewTallyDesktopSync", "Groups")
	mongoHandler := NewMongoHandler(collection)
	results := mongoHandler.FindDocuments(DocumentFilter{
		Filter:        bson.M{},
		Limit:         0,
		Offset:        0,
		UsePagination: false,
		Ctx:           context.TODO(),
	})
	if results.Err != nil {
		log.Printf("[+] Coult not get companies")
		panic("No Group Cache Built")
	}
	log.Printf("[+] Building Cache for %v companies \n", len(results.Data))

	companyGroups := make(map[string][]Group)
	for _, group := range results.Data {
		if group["GUID"] != nil {
			newGroup := Group{
				GUID:   group["GUID"].(string),
				Parent: group["Parent"].(string),
				Name:   group["Name"].(string),
			}
			key := (newGroup.GUID)[:36]
			companyGroups[key] = append(companyGroups[key], newGroup)
		}
	}

	for comanyId, groupList := range companyGroups {
		log.Printf("[+] Building Cache for %s  \n", comanyId)
		gcm.Build(comanyId, groupList)
		log.Printf("[+] Cache Built for %s  \n", comanyId)
	}

}
