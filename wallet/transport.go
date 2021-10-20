package wallet

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/arhyth/genwallet/errorrrs"
	"github.com/go-kit/kit/endpoint"
)

var (
	rgxpWalletsIDPayments = regexp.MustCompile(`/wallets/([\w-]+)/payments`)
	rgxpWalletsID         = regexp.MustCompile(`/wallets/([\w-]+)`)
)

// Go-kit http transport signature funcs

func MakeWalletGetEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountRequest)
		return svc.GetAccount(req)
	}
}

func DecodeHTTPGetAccountReq(_ context.Context, req *http.Request) (interface{}, error) {
	var getReq GetAccountRequest
	match := rgxpWalletsID.FindStringSubmatch(req.URL.Path)
	if len(match) < 2 {
		return nil, &errorrrs.E{
			ID:  errorrrs.BadRequest,
			Msg: "malformed path: should be of `/wallets/{id}` format",
		}
	}
	getReq.ID = match[1]

	return getReq, nil
}

func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

func MakeWalletListEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ListAccountsRequest)
		return svc.ListAccounts(req)
	}
}

func DecodeHTTPListAccountsReq(_ context.Context, req *http.Request) (interface{}, error) {
	var listReq ListAccountsRequest
	cur := req.URL.Query().Get("currency")
	if cur != "" {
		listReq.Currency = &cur
	}

	return listReq, nil
}

func MakeWalletCreateEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)
		return svc.CreateAccount(req)
	}
}

func DecodeHTTPCreateAccountReq(_ context.Context, req *http.Request) (interface{}, error) {
	var createReq CreateAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		return nil, err
	}

	return createReq, nil
}

func MakePaymentsIndexEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ListPaymentsRequest)
		return svc.ListPayments(req)
	}
}

func DecodeHTTPListPaymentsReq(_ context.Context, req *http.Request) (interface{}, error) {
	var listPayments ListPaymentsRequest
	match := rgxpWalletsIDPayments.FindStringSubmatch(req.URL.Path)
	if len(match) < 2 {
		return nil, &errorrrs.E{
			ID:  errorrrs.BadRequest,
			Msg: "malformed path: should be of `/wallets/{id}/payments` format",
		}
	}
	listPayments.ID = match[1]

	return listPayments, nil
}

func MakePaymentsPostEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreatePaymentRequest)
		return svc.CreatePayment(req)
	}
}

func DecodeHTTPPostPaymentsReq(_ context.Context, req *http.Request) (interface{}, error) {
	var paymentReq CreatePaymentRequest
	if err := json.NewDecoder(req.Body).Decode(&paymentReq); err != nil {
		return nil, &errorrrs.E{
			ID:  errorrrs.InternalServerError,
			Msg: err.Error(),
		}
	}
	match := rgxpWalletsIDPayments.FindStringSubmatch(req.URL.Path)
	if len(match) < 2 {
		return nil, &errorrrs.E{
			ID:  errorrrs.BadRequest,
			Msg: "malformed path: should be of `/wallets/{id}/payments` format",
		}
	}
	paymentReq.Self = match[1]

	return paymentReq, nil
}

func MakeListTransfersEndpt(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ListTransfersRequest)
		return svc.ListTransfers(req)
	}
}

func DecodeHTTPListTransfersReq(_ context.Context, req *http.Request) (interface{}, error) {
	var listReq ListTransfersRequest
	currency := req.URL.Query().Get("currency")
	if currency != "" {
		listReq.Currency = &currency
	}
	from := req.URL.Query().Get("from")
	if from != "" {
		listReq.From = &from
	}
	to := req.URL.Query().Get("to")
	if to != "" {
		listReq.To = &to
	}

	return listReq, nil
}
