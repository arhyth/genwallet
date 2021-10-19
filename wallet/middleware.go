package wallet

import "github.com/rs/zerolog"

var _ Service = (*ValidationMiddleware)(nil)

// ValidationMiddleware is as name suggests validation middleware for wallet
// service. Validation is done in service layer so that it concerns only business/domain
// and does not need change whatever transport/protocol is used to expose the API
type ValidationMiddleware struct {
	Next   Service
	Logger *zerolog.Logger
}

// TODO: implement validation :D

func (vm *ValidationMiddleware) ListAccounts(req ListAccountsRequest) ([]Account, error) {
	return vm.Next.ListAccounts(req)
}

func (vm *ValidationMiddleware) GetAccount(req GetAccountRequest) (Account, error) {
	return vm.Next.GetAccount(req)
}

func (vm *ValidationMiddleware) CreateAccount(req CreateAccountRequest) (Account, error) {
	return vm.Next.CreateAccount(req)
}

func (vm *ValidationMiddleware) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	return vm.Next.ListPayments(req)
}

func (vm *ValidationMiddleware) CreatePayment(req CreatePaymentRequest) (Payment, error) {
	return vm.Next.CreatePayment(req)
}

func (vm *ValidationMiddleware) ListTransfers(req ListTransfersRequest) ([]Transfer, error) {
	return vm.Next.ListTransfers(req)
}
