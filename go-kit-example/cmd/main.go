package main

import (
	"github.com/go-kit/kit/log"
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

	var (
		service     = aggservice.New()
		endpoints   = aggendpoint.New(service, logger)
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
