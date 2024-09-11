package models

type RequestFilter struct {
	Batch      Pagination
	SearchText string
	SearchKey  string
	SortKey    string
	SortOrder  string
}
