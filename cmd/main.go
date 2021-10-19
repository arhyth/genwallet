package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arhyth/genwallet/config"
	"github.com/arhyth/genwallet/errorrrs"
	"github.com/arhyth/genwallet/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	// Logging
	logger := zerolog.New(os.Stderr)

	// Config
	cfg, err := config.GetAPIConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("genwallet server start: config parse fail")
	}
	httpAddr := cfg.AddrPort

	// Transport
	r := chi.NewMux()

	ok := []byte("OK")
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write(ok)
	})
	r.Get("/healthcheck", okHandler)

	repo, err := wallet.NewRepo(cfg.DBConnStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("genwallet server start: wallet.NewRepo")
	}
	walletSvc := &wallet.ServiceImpl{
		Repo: repo,
	}

	serverErrcoder := httptransport.ServerErrorEncoder(errorrrs.GokitErrorEncoder)
	walletsIndexHandler := httptransport.NewServer(
		wallet.MakeWalletListEndpt(walletSvc),
		wallet.DecodeHTTPListAccountsReq,
		wallet.EncodeJSONResponse,
		serverErrcoder,
	)
	walletGetHandler := httptransport.NewServer(
		wallet.MakeWalletGetEndpt(walletSvc),
		wallet.DecodeHTTPGetAccountReq,
		wallet.EncodeJSONResponse,
		serverErrcoder,
	)
	walletCreateHandler := httptransport.NewServer(
		wallet.MakeWalletCreateEndpt(walletSvc),
		wallet.DecodeHTTPCreateAccountReq,
		wallet.EncodeJSONResponse,
		serverErrcoder,
	)
	walletPaymentsIndexHandler := httptransport.NewServer(
		wallet.MakePaymentsIndexEndpt(walletSvc),
		wallet.DecodeHTTPListPaymentsReq,
		wallet.EncodeJSONResponse,
		serverErrcoder,
	)
	walletPostPaymentHandler := httptransport.NewServer(
		wallet.MakePaymentsPostEndpt(walletSvc),
		wallet.DecodeHTTPPostPaymentsReq,
		wallet.EncodeJSONResponse,
		serverErrcoder,
	)
	ledgerHandler := httptransport.NewServer(
		wallet.MakeListTransfersEndpt(walletSvc),
		wallet.DecodeHTTPListTransfersReq,
		wallet.EncodeJSONResponse,
		serverErrcoder,
	)
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
			Str("addr", httpAddr).
			Msg("genwallet server start")
		errc <- http.ListenAndServe(httpAddr, r)
	}()

	logger.Err(<-errc).Msg("genwallet server exit")
}
