package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/0xDarkXnight/Toll-Calculator-Microservice/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "the listen address of the HTTP server")
	flag.Parse()
	var (
		store   = NewMemoryStore()
		service = NewInvoiceAggregator(store)
	)
	service = NewLogMiddleware(service)
	makeHTTPTransport(*listenAddr, service)
}

func makeHTTPTransport(listenAddr string, service Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(service))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := service.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
