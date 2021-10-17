package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arhyth/genwallet/wallet"
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

	ok := []byte("OK")
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write(ok)
	})
	r.Get("/healthcheck", okHandler)

	walletSvc := wallet.NewSimpleWalletService()
	walletsIndexHandler := wallet.NewWalletListHandler(walletSvc)
	walletGetHandler := wallet.NewWalletGetHandler(walletSvc)
	walletCreateHandler := wallet.NewWalletCreateHandler(walletSvc)
	walletPaymentsIndexHandler := wallet.NewWalletPaymentsIndexHandler(walletSvc)
	walletPostPaymentHandler := wallet.NewWalletPostPaymentHandler(walletSvc)
	ledgerHandler := wallet.NewWalletLedgerHandler(walletSvc)
	r.Method("GET", "/wallets", walletsIndexHandler)
	r.Method("POST", "/wallets", walletCreateHandler)
	r.Method("GET", "/wallets/{id}", walletGetHandler)
	r.Method("GET", "/wallets/{id}/payments", walletPaymentsIndexHandler)
	r.Method("POST", "/wallets/{id}/payments", walletPostPaymentHandler)
	r.Method("GET", "/transfers", ledgerHandler)

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
