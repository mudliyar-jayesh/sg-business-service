package main

import (
	"fmt"
	"log"
	"net/http"
	"sg-business-service/config"
	"sg-business-service/endpoints"
	"sg-business-service/handlers"
	dspEndpoints "sg-business-service/modules/dsp/endpoints"
	osEndpoints "sg-business-service/modules/outstanding/endpoints"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Priority, companyid")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
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

	http.Handle("/os/send-email", corsMiddleware(http.HandlerFunc(endpoints.SendLedgerEmail)))

	http.Handle("/os/aggr", corsMiddleware(http.HandlerFunc(endpoints.TempOS)))

	// outstanding settings endpoints
	http.Handle("/os-setting/create", corsMiddleware(http.HandlerFunc(endpoints.CreateOsSetting)))
	http.Handle("/os-setting/update", corsMiddleware(http.HandlerFunc(endpoints.UpdateOsSetting)))
	http.Handle("/os-setting/get", corsMiddleware(http.HandlerFunc(endpoints.GetSetting)))

	// outstanding endpoints
	http.Handle("/os/search/ledgers", corsMiddleware(http.HandlerFunc(endpoints.SearchLedgers)))
	http.Handle("/os/get/groups", corsMiddleware(http.HandlerFunc(endpoints.GetCachedGroups)))
	http.Handle("/os/get/report", corsMiddleware(http.HandlerFunc(endpoints.GetOutstandingReport)))

	// inventory endpoints
	http.Handle("/stock-items/get/report", corsMiddleware(http.HandlerFunc(endpoints.GetStockItemReport)))
	http.Handle("/stock-group/get/names", corsMiddleware(http.HandlerFunc(endpoints.GetItemGroupNames)))

	// sync info endpoints
	http.Handle("/sync-info/get", corsMiddleware(http.HandlerFunc(endpoints.GetLastSync)))

	// collection endpoints
	http.Handle("/collection/get", corsMiddleware(http.HandlerFunc(osEndpoints.GetCollectionOverview)))

<<<<<<< HEAD
    // followup endpoints
	http.Handle("/os/followup/sample", corsMiddleware(http.HandlerFunc(osEndpoints.SampleFollowUp)))
	http.Handle("/os/followup/status", corsMiddleware(http.HandlerFunc(osEndpoints.GetBillStatusList)))
    http.Handle("os/followup/create", corsMiddleware(http.HandlerFunc(osEndpoints.CreateFollowUp)))
=======
	// dsp endponts
	http.Handle("/dsp/upload", corsMiddleware(http.HandlerFunc(dspEndpoints.UploadCsvFile)))
>>>>>>> b4f9d60db7912f1e886e0d8a7e7d7ae1b27ceb80

	fmt.Println("Server starting on port 35001...")
	log.Fatal(http.ListenAndServe(":35001", nil))
}
