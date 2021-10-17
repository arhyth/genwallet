package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	rgxpWalletsIDPayments = regexp.MustCompile(`/wallets/([\w-]+)/payments`)
	rgxpWalletsID         = regexp.MustCompile(`/wallets/([\w-]+)`)

	errBadRequest = errors.New("bad request")
)

func NewWalletListHandler(walletSvc Service) *httptransport.Server {
	return httptransport.NewServer(
		makeWalletListEndpt(walletSvc),
		decodeListReqFunc,
		encodeResFunc,
	)
}

func NewWalletGetHandler(walletSvc Service) *httptransport.Server {
	return httptransport.NewServer(
		makeWalletGetEndpt(walletSvc),
		decodeGetReqFunc,
		encodeResFunc,
	)
}

func NewWalletCreateHandler(walletSvc Service) *httptransport.Server {
	return httptransport.NewServer(
		makeWalletCreateEndpt(walletSvc),
		decodeCreateReqFunc,
		encodeResFunc,
	)
}

func NewWalletPaymentsIndexHandler(walletSvc Service) *httptransport.Server {
	return httptransport.NewServer(
		makePaymentsIndexEndpt(walletSvc),
		decodeListPaymentsReqFunc,
		encodeResFunc,
	)
}

func NewWalletPostPaymentHandler(walletSvc Service) *httptransport.Server {
	return httptransport.NewServer(
		makePaymentsPostEndpt(walletSvc),
		decodePostPaymentsReqFunc,
		encodeResFunc,
	)
}

func NewWalletLedgerHandler(walletSvc Service) *httptransport.Server {
	return httptransport.NewServer(
		makeLedgerEndpt(walletSvc),
		decodeLedgerReqFunc,
		encodeResFunc,
	)
}

// go-kit helper funcs

func makeWalletGetEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		return svc.Get(req)
	}
}

func decodeGetReqFunc(_ context.Context, req *http.Request) (interface{}, error) {
	var getReq GetRequest
	match := rgxpWalletsID.FindStringSubmatch(req.URL.Path)
	if len(match) < 2 {
		return nil, errBadRequest
	}
	getReq.ID = match[1]

	return getReq, nil
}

func encodeResFunc(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

func makeWalletListEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		return svc.List(req)
	}
}

func decodeListReqFunc(_ context.Context, req *http.Request) (interface{}, error) {
	var listReq ListRequest
	cur := req.URL.Query().Get("currency")
	if cur != "" {
		listReq.Currency = &cur
	}

	return listReq, nil
}

func makeWalletCreateEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		return svc.Create(req)
	}
}

func decodeCreateReqFunc(_ context.Context, req *http.Request) (interface{}, error) {
	var createReq CreateRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		return nil, err
	}

	return createReq, nil
}

func makePaymentsIndexEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ListPaymentsRequest)
		return svc.ListPayments(req)
	}
}

func decodeListPaymentsReqFunc(_ context.Context, req *http.Request) (interface{}, error) {
	var listPayments ListPaymentsRequest
	match := rgxpWalletsIDPayments.FindStringSubmatch(req.URL.Path)
	if len(match) < 2 {
		return nil, errBadRequest
	}
	listPayments.ID = match[1]

	return listPayments, nil
}

func makePaymentsPostEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(PaymentRequest)
		return svc.Transfer(req)
	}
}

func decodePostPaymentsReqFunc(_ context.Context, req *http.Request) (interface{}, error) {
	var paymentReq PaymentRequest
	if err := json.NewDecoder(req.Body).Decode(&paymentReq); err != nil {
		return nil, err
	}
	match := rgxpWalletsIDPayments.FindStringSubmatch(req.URL.Path)
	if len(match) < 2 {
		return nil, errBadRequest
	}
	paymentReq.From = match[1]

	return paymentReq, nil
}

func makeLedgerEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(TransferLedgerRequest)
		return svc.TransferLedger(req)
	}
}

func decodeLedgerReqFunc(_ context.Context, req *http.Request) (interface{}, error) {
	var ledgerReq TransferLedgerRequest
	from := req.URL.Query().Get("from_id")
	if from != "" {
		ledgerReq.From = &from
	}
	to := req.URL.Query().Get("to_id")
	if to != "" {
		ledgerReq.To = &to
	}
	since := req.URL.Query().Get("since")
	if since != "" {
		ledgerReq.Since = &since
	}
	upto := req.URL.Query().Get("upto")
	if from != "" {
		ledgerReq.Upto = &upto
	}

	return ledgerReq, nil
}
