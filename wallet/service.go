package wallet

import (
	"errors"
	"time"
)

type EntryType int

// EntryType is accounting ledger entry type. It is relative
// to an account: an incoming entry to one account is an outgoing
// entry to another
const (
	Incoming EntryType = iota + 1
	Outgoing
)

var (
	ErrNotFound = errors.New("not found")
)

type Account struct {
	ID        string    `json:"id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Note: The API seems a bit unintuitive since the business/domain model
// is a bit confusing. Nevertheless, I tried to make it as simple as possible.
type Service interface {
	List(ListRequest) ([]Account, error)
	Get(GetRequest) (Account, error)
	Create(CreateRequest) (Account, error)
	ListPayments(ListPaymentsRequest) ([]Payment, error)
	Transfer(PaymentRequest) (Payment, error)
	TransferLedger(TransferLedgerRequest) ([]Transfer, error)
}

type GetRequest struct {
	ID string `json:"id"`
}

type PaymentRequest struct {
	From   string
	To     string
	Amount float64
}

type Payment struct {
	From   string
	To     string
	Amount float64
	// to be honest this seems redundant as payment is always
	// in outgoing fashion with respect to a wallet/account;
	// also, `From` and `To` fields already indicate direction
	Direction EntryType
}

type Transfer struct {
	From      string
	To        string
	Amount    float64
	CreatedAt time.Time
}

type TransferLedgerRequest struct {
	// wallet/account IDs
	From *string
	To   *string
	// format: `2006-01-02`
	// TODO: add validation in service middleware
	Since *string
	Upto  *string
}

func NewSimpleWalletService() *SimpleService {
	return &SimpleService{}
}

type SimpleService struct{}

var _ Service = (*SimpleService)(nil)

func (ws *SimpleService) Get(req GetRequest) (Account, error) {
	return Account{
		ID:       req.ID,
		Balance:  100.0,
		Currency: "USD",
	}, nil
}

type ListRequest struct {
	Currency *string
}

func (ws *SimpleService) List(req ListRequest) ([]Account, error) {
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

type CreateRequest struct {
	ID       string  `json:"id"`
	InitAmt  float64 `json:"init_amt"`
	Currency string  `json:"currency" `
}

func (ws *SimpleService) Create(req CreateRequest) (Account, error) {
	now := time.Now().UTC()
	return Account{
		ID:        req.ID,
		Balance:   req.InitAmt,
		Currency:  req.Currency,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (ws *SimpleService) Transfer(req PaymentRequest) (Payment, error) {
	return Payment{
		From:      req.From,
		To:        req.To,
		Amount:    req.Amount,
		Direction: Outgoing,
	}, nil
}

type ListPaymentsRequest struct {
	ID string `json:"id"`
}

func (ws *SimpleService) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	toother := "toOther123"
	fromother := "fromOther123"
	payments := []Payment{
		{
			From:      req.ID,
			To:        toother,
			Amount:    10.0,
			Direction: Outgoing,
		},
		{
			From:      fromother,
			To:        req.ID,
			Amount:    80.0,
			Direction: Incoming,
		},
	}

	return payments, nil
}

func (ws *SimpleService) TransferLedger(req TransferLedgerRequest) ([]Transfer, error) {
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
