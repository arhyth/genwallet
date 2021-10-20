package wallet

import (
	"github.com/arhyth/genwallet/errorrrs"
	"github.com/rs/zerolog"
)

var _ Service = (*ValidationMiddleware)(nil)

// ValidationMiddleware is as name suggests validation middleware for wallet
// service. Validation is done in service layer so that it concerns only business/domain
// and does not need change whatever transport/protocol is used to expose the API
type ValidationMiddleware struct {
	Next   Service
	Logger *zerolog.Logger
}

func (vm *ValidationMiddleware) ListAccounts(req ListAccountsRequest) ([]Account, error) {
	return vm.Next.ListAccounts(req)
}

func (vm *ValidationMiddleware) GetAccount(req GetAccountRequest) (Account, error) {
	return vm.Next.GetAccount(req)
}

func (vm *ValidationMiddleware) CreateAccount(req CreateAccountRequest) (Account, error) {
	if _, exist := ValidCurrencies[req.Currency]; !exist {
		return Account{}, &errorrrs.E{
			ID:  errorrrs.BadRequest,
			Msg: "invalid currency",
		}
	}

	return vm.Next.CreateAccount(req)
}

func (vm *ValidationMiddleware) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	return vm.Next.ListPayments(req)
}

func (vm *ValidationMiddleware) CreatePayment(req CreatePaymentRequest) (Payment, error) {
	// Note: we can check here if payee and payer wallet currencies match by adding
	// a dependency to wallet.Repository. However, since balance access and updates still need
	// to be serialized we just piggyback currency matching validation on the transaction.

	if req.Self == req.To {
		return Payment{}, &errorrrs.E{
			ID:  errorrrs.BadRequest,
			Msg: "transfer recipient is same wallet",
		}
	}

	if req.Amount == 0 {
		return Payment{}, &errorrrs.E{
			ID:  errorrrs.BadRequest,
			Msg: "transfer amount is `0`",
		}
	}

	return vm.Next.CreatePayment(req)
}

func (vm *ValidationMiddleware) ListTransfers(req ListTransfersRequest) ([]Transfer, error) {
	return vm.Next.ListTransfers(req)
}
