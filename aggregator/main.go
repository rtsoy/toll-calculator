package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rtsoy/toll-calculator/aggregator/client"
	"github.com/rtsoy/toll-calculator/types"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var (
		httpAddr = os.Getenv("AGG_HTTP_ENDPOINT")
		grpcAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		store    = makeStore()
		svc      = NewInvoiceAggregator(store)
	)

	svc = NewLogMiddleware(svc)
	svc = NewMetricsMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcAddr, svc))
	}()
	time.Sleep(time.Second * 1)
	c, err := client.NewGRPCClient(grpcAddr)
	if err != nil {
		log.Fatal("NewGRPCClient:", err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 534.3,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal("Aggregate:", err)
	}

	makeHTTPTransport(httpAddr, svc)
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

	// Prometheus
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "Method Not Allowed",
			})
			return
		}

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
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "Method Not Allowed",
			})
			return
		}

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

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")

	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type: %s", storeType)
		return nil
	}
}
