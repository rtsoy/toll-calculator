package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rtsoy/toll-calculator/types"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
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

	aggMetricHandler := NewHTTPMetricHandler("aggregate")
	invMetricHandler := NewHTTPMetricHandler("calculate")

	http.HandleFunc("/aggregate", makeHTTPHandler(
		aggMetricHandler.instrument(handleAggregate(svc))),
	)
	http.HandleFunc("/invoice", makeHTTPHandler(
		invMetricHandler.instrument(handleGetInvoice(svc))),
	)

	// No metrics
	//http.HandleFunc("/invoice", makeHTTPHandler(handleGetInvoice(svc)))
	//http.HandleFunc("/aggregate", makeHTTPHandler(handleAggregate(svc)))

	// Prometheus
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(listenAddr, nil))
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
