package wallet

import (
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

type Account struct {
	ID        string    `json:"id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Note: The API seems a bit unintuitive since the business/domain model
// is a bit confusing; unsure if the `Service` interface should be broken
// into 2 interfaces and how to model `transfer`s/`payment`s.
// We could also propagate request context here but since we will not
// be making use of cancellation or deadline, seems curently unnecessary.
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

type ListRequest struct {
	Currency *string
}

type CreateRequest struct {
	ID       string  `json:"id"`
	InitAmt  float64 `json:"init_amt"`
	Currency string  `json:"currency" `
}

type ListPaymentsRequest struct {
	ID string `json:"id"`
}

var _ Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	// Note: the Repo should in most instances be a separate interface. However,
	// since genwallet operations are mostly CRUD and derive their concrete
	// implementation from the database, we can simply use the same interface.
	// Feel free to discuss and create a separate one when the need arises.
	Repo Repository
}

func (ws *ServiceImpl) Get(req GetRequest) (Account, error) {
	return ws.Repo.GetAccount(req)
}

func (ws *ServiceImpl) List(req ListRequest) ([]Account, error) {
	return ws.Repo.ListAccounts(req)
}

func (ws *ServiceImpl) Create(req CreateRequest) (Account, error) {
	return ws.Repo.CreateAccount(req)
}

func (ws *ServiceImpl) Transfer(req PaymentRequest) (Payment, error) {
	return ws.Repo.CreateTransfer(req)
}

func (ws *ServiceImpl) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	// Note: we make use of same method from the DB as for `TransferLedger`
	// since `Payment`s are only a `Service` "domain object" and exist in the DB
	// layer also as `Transfer`s
	transferReq := TransferLedgerRequest{
		From: &req.ID,
		To:   &req.ID,
	}

	transfers, err := ws.Repo.ListTransfers(transferReq)
	if err != nil {
		return nil, err
	}

	payments := make([]Payment, len(transfers))
	for i := range transfers {
		p := Payment{
			From:   transfers[i].From,
			To:     transfers[i].To,
			Amount: transfers[i].Amount,
		}
		if req.ID == transfers[i].From {
			p.Direction = Outgoing
		} else {
			p.Direction = Incoming
		}
		payments[i] = p
	}

	return payments, nil
}

func (ws *ServiceImpl) TransferLedger(req TransferLedgerRequest) ([]Transfer, error) {
	// Note: since `From` and `To` work together as a `where... OR` query,
	// it is up to the client component (e.g. web, mobile) to filter on their
	// end for cases where the user explicitly wants a `where... AND`

	return ws.Repo.ListTransfers(req)
}
