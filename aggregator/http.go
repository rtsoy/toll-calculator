package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rtsoy/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, msg string) APIError {
	return APIError{
		Code:    code,
		Message: msg,
	}
}

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})

	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
		Name:      "aggregator",
	})

	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &HTTPMetricHandler{
		reqCounter: reqCounter,
		errCounter: errCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()

			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
			}).Info()

			h.reqLatency.Observe(latency)
		}(time.Now())

		h.reqCounter.Inc()

		return next(w, r)
	}
}

func makeHTTPHandler(next HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			if apiError, ok := err.(APIError); ok {
				writeJSON(w, apiError.Code, apiError)
				return
			}

			someError := NewAPIError(http.StatusInternalServerError, err.Error())
			writeJSON(w, someError.Code, someError)
			return
		}
	}
}

func handleGetInvoice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return NewAPIError(http.StatusMethodNotAllowed, "Method Not Allowed")
		}

		values, ok := r.URL.Query()["obuID"]
		if !ok {
			return NewAPIError(http.StatusBadRequest, "Missing obuID")
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			return NewAPIError(http.StatusBadRequest, "Invalid obuID")
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			return NewAPIError(http.StatusNotFound, fmt.Sprintf("Resource with id=%d does not exist", obuID))
		}

		return writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return NewAPIError(http.StatusMethodNotAllowed, "Method Not Allowed")
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		if err := svc.AggregateDistance(distance); err != nil {
			return NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		return writeJSON(w, http.StatusOK, distance)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
