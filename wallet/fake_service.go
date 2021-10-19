package wallet

import (
	"time"
)

func NewSimpleWalletService() *SimpleService {
	return &SimpleService{}
}

type SimpleService struct{}

var _ Service = (*SimpleService)(nil)

func (ws *SimpleService) GetAccount(req GetAccountRequest) (Account, error) {
	return Account{
		ID:       req.ID,
		Balance:  100.0,
		Currency: "USD",
	}, nil
}

func (ws *SimpleService) ListAccounts(req ListAccountsRequest) ([]Account, error) {
	accounts := []Account{
		{
			ID:       "bob-1234",
			Balance:  10000.0,
			Currency: "JPY",
		},
		{
			ID:       "alice-5678",
			Balance:  100.0,
			Currency: "USD",
		},
		{
			ID:       "sato-91011",
			Balance:  1000.0,
			Currency: "CNY",
		},
	}

	if req.Currency != nil {
		for i := range accounts {
			accounts[i].Currency = *req.Currency
		}
	}

	return accounts, nil
}

func (ws *SimpleService) CreateAccount(req CreateAccountRequest) (Account, error) {
	now := time.Now().UTC()
	return Account{
		ID:        req.ID,
		Balance:   req.InitAmt,
		Currency:  req.Currency,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (ws *SimpleService) CreatePayment(req CreatePaymentRequest) (Payment, error) {
	return Payment{
		Self:      req.Self,
		To:        &req.To,
		Amount:    req.Amount,
		Direction: Outgoing,
	}, nil
}

func (ws *SimpleService) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	toother := "toOther123"
	fromother := "fromOther123"
	payments := []Payment{
		{
			Self:      req.ID,
			To:        &toother,
			Amount:    10.0,
			Direction: Outgoing,
		},
		{
			Self:      req.ID,
			From:      &fromother,
			Amount:    80.0,
			Direction: Incoming,
		},
	}

	return payments, nil
}

func (ws *SimpleService) ListTransfers(req ListTransfersRequest) ([]Transfer, error) {
	ben := "ben123"
	alice := "alice456"
	now := time.Now().UTC()

	transfers := []Transfer{
		{
			From:      ben,
			To:        alice,
			Amount:    10.0,
			CreatedAt: now.AddDate(0, 1, 0),
		},
		{
			From:      alice,
			To:        ben,
			Amount:    80.0,
			CreatedAt: now.AddDate(0, 0, 20),
		},
	}

	return transfers, nil
}
