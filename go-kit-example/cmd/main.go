package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rtsoy/toll-calculator/go-kit-example/pkg/aggendpoint"
	"github.com/rtsoy/toll-calculator/go-kit-example/pkg/aggservice"
	"github.com/rtsoy/toll-calculator/go-kit-example/pkg/aggtransport"
	"net"
	"net/http"
	"os"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var duration metrics.Histogram
	{
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "tollCalculator",
			Subsystem: "aggService",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	var (
		service     = aggservice.New(logger)
		endpoints   = aggendpoint.New(service, logger, duration)
		httpHandler = aggtransport.NewHTTPHandler(endpoints, logger)
	)

	httpListener, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "HTTP", "addr", ":3000")

	if err := http.Serve(httpListener, httpHandler); err != nil {
		panic(err)
	}
}
