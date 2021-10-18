package wallet

import (
	"database/sql"
)

type Repository interface {
	ListAccounts(ListRequest) ([]Account, error)
	GetAccount(GetRequest) (Account, error)
	CreateAccount(CreateRequest) (Account, error)
	CreateTransfer(PaymentRequest) (Payment, error)
	ListTransfers(TransferLedgerRequest) ([]Transfer, error)
}

var _ Repository = (*Repo)(nil)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) ListAccounts(req ListRequest) ([]Account, error) {
	return nil, nil
}

func (r *Repo) GetAccount(req GetRequest) (Account, error) {
	return Account{}, nil
}

func (r *Repo) CreateAccount(req CreateRequest) (Account, error) {
	return Account{}, nil
}

func (r *Repo) CreateTransfer(req PaymentRequest) (Payment, error) {
	return Payment{}, nil
}

func (r *Repo) ListTransfers(TransferLedgerRequest) ([]Transfer, error) {
	return nil, nil
}
