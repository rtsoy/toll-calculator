package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rtsoy/toll-calculator/aggregator/client"
	"github.com/rtsoy/toll-calculator/types"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	httpAddr := flag.String("httpAddr", ":3000", "The listen address of the HTTP server")
	grpcAddr := flag.String("grpcAddr", ":3001", "The listen address of the HTTP server")

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)

	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcAddr, svc))
	}()
	time.Sleep(time.Second * 1)
	c, err := client.NewGRPCClient(*grpcAddr)
	if err != nil {
		log.Fatal("NewGRPCClient:", err)
	}
	if _, err := c.Client.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 534.3,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal("Aggregate:", err)
	}

	makeHTTPTransport(*httpAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))

	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obuID"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "missing obuID",
			})
			return
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid obuID",
			})
			return
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
