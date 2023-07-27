package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rtsoy/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listen address of the HTTP server")
	flag.Parse()

	var (
		clnt       = client.NewHTTPClient("http://127.0.0.1:3000")
		invHandler = NewInvoiceHandler(clnt)
	)

	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))

	logrus.Infof("gateway HTTP server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	values, ok := r.URL.Query()["obuID"]
	if !ok {
		return writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing obuID",
		})
	}

	obuID, err := strconv.Atoi(values[0])
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid obuID",
		})
	}

	inv, err := h.client.GetInvoice(r.Context(), obuID)
	if err != nil {
		return writeJSON(w, http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("could not find a distance for obuID=%d", obuID),
		})
	}

	return writeJSON(w, http.StatusOK, inv)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
