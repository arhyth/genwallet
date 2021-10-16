package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func main() {
	var (
		httpAddr = flag.String("ADDR_PORT", ":8000", "Address for HTTP (JSON) server")
	)
	flag.Parse()

	// Logging
	logger := zerolog.New(os.Stderr)

	// Transport
	r := chi.NewMux()

	// Interrupt
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Run
	go func() {
		logger.Info().
			Str("transport", "HTTP").
			Str("addr", *httpAddr).
			Msg("genwallet server start")
		errc <- http.ListenAndServe(*httpAddr, r)
	}()

	logger.Err(<-errc).Msg("genwallet server exit")
}
