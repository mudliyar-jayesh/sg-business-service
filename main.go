package main 

import (
    "log"
    "fmt"
    "net/http"
    "sg-business-service/config" 
    "sg-business-service/handlers" 
    "sg-business-service/endpoints" 
)


func main() {
    mongoConfig := config.LoadMongoConfig()
    handlers.ConnectToMongo(mongoConfig)
    handlers.MakeGroupCache()

    // outstanding endpoints
    http.HandleFunc("/os/get/groups", endpoints.GetCachedGroups)
    http.HandleFunc("/os/get/report", endpoints.GetOutstandingReport)

    // inventory endpoints
    http.HandleFunc("/stock-items/get/report", endpoints.GetStockItemReport)

    fmt.Println("Server starting on port 35001...")
    log.Fatal(http.ListenAndServe(":35001", nil))
}
