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

func (vm *ValidationMiddleware) List(req ListRequest) ([]Account, error) {
	return vm.Next.List(req)
}

func (vm *ValidationMiddleware) Get(req GetRequest) (Account, error) {
	return vm.Next.Get(req)
}

func (vm *ValidationMiddleware) Create(req CreateRequest) (Account, error) {
	return vm.Next.Create(req)
}

func (vm *ValidationMiddleware) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	return vm.Next.ListPayments(req)
}

func (vm *ValidationMiddleware) Transfer(req PaymentRequest) (Payment, error) {
	return vm.Next.Transfer(req)
}

func (vm *ValidationMiddleware) TransferLedger(req TransferLedgerRequest) ([]Transfer, error) {
	return vm.Next.TransferLedger(req)
}
