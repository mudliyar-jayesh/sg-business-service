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

    // outstanding endpoints
    http.HandleFunc("/os/get/report", endpoints.GetOutstandingReport)
    fmt.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
