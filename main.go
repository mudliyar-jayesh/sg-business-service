package main 

import (
    "log"
    "fmt"
    "net/http"
    "sg-business-service/config" 
    "sg-business-service/handlers" 
    "sg-business-service/endpoints" 
)

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            return
        }

        // Call the next handler
        next.ServeHTTP(w, r)
    })
}


func main() {
    mongoConfig := config.LoadMongoConfig()
    handlers.ConnectToMongo(mongoConfig)
    handlers.MakeGroupCache()

    // outstanding endpoints
    http.Handle("/os/get/groups", corsMiddleware(http.HandlerFunc(endpoints.GetCachedGroups)))
    http.Handle("/os/get/report", corsMiddleware(http.HandlerFunc(endpoints.GetOutstandingReport)))

    // inventory endpoints
    http.Handle("/stock-items/get/report", corsMiddleware(http.HandlerFunc(endpoints.GetStockItemReport)))

    fmt.Println("Server starting on port 35001...")
    log.Fatal(http.ListenAndServe(":35001", nil))
}
