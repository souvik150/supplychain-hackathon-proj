package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

func main(){
	listenAddr := flag.String("listen-addr", ":3000", "server listen address")
	flag.Parse()

	store := NewMemoryStore()
	var (
		svc =  NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)

	makeHTTPTransport(*listenAddr, svc)

	http.ListenAndServe(*listenAddr, nil)
}

func makeHTTPTransport(listenlistenAddr string, svc Aggregator) {
	fmt.Println("HTTP Transport started at: ", listenlistenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenlistenAddr, nil)
}

func handleGetInvoice (svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		obuid, err := strconv.Atoi(r.URL.Query().Get("obuid"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "obuid must be an integer"})
			return
		}
		
		inv, err := svc.CalculateInvoice(obuid)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func writeJSON(rw http.ResponseWriter,status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}