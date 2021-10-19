package wallet

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/arhyth/genwallet/errorrrs"
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
// be making use of cancellation or deadline, seems currently unnecessary.
type Service interface {
	ListAccounts(ListAccountsRequest) ([]Account, error)
	GetAccount(GetAccountRequest) (Account, error)
	CreateAccount(CreateAccountRequest) (Account, error)
	ListPayments(ListPaymentsRequest) ([]Payment, error)
	CreatePayment(CreatePaymentRequest) (Payment, error)
	ListTransfers(ListTransfersRequest) ([]Transfer, error)
}

type GetAccountRequest struct {
	ID string `json:"id"`
}

type EntryType int

// EntryType is accounting ledger entry type. It is relative
// to an account: an incoming entry to one account is an outgoing
// entry to another
const (
	Incoming EntryType = iota + 1
	Outgoing
)

var (
	incomingJSON, _ = json.Marshal("Incoming")
	outgoingJSON, _ = json.Marshal("Outgoing")
)

func (et *EntryType) MarshalJSON() ([]byte, error) {
	switch *et {
	case Incoming:
		return incomingJSON, nil
	case Outgoing:
		return outgoingJSON, nil
	default:
		return nil, errors.New("undefined EntryType(Direction)")
	}
}

type Payment struct {
	Self      string    `json:"account"`
	From      *string   `json:"from_account,omitempty"`
	To        *string   `json:"to_account,omitempty"`
	Amount    float64   `json:"amount"`
	Direction EntryType `json:"direction"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID        int
	From      string
	To        string
	Amount    float64
	CreatedAt time.Time
}

type ListTransfersRequest struct {
	// wallet/account IDs
	From *string
	To   *string
}

type ListAccountsRequest struct {
	Currency *string
}

type CreateAccountRequest struct {
	ID       string  `json:"id"`
	InitAmt  float64 `json:"init_amt"`
	Currency string  `json:"currency" `
}

type ListPaymentsRequest struct {
	ID string `json:"id"`
}

type CreatePaymentRequest struct {
	Self   string  `json:"account"`
	To     string  `json:"to_account"`
	Amount float64 `json:"amount"`
}

var _ Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	Repo Repository
}

func (ws *ServiceImpl) GetAccount(req GetAccountRequest) (Account, error) {
	acct, err := ws.Repo.GetAccount(req)
	if err != nil {
		// Yes, yes, reader, this seems unergonomic. I haven't found a better
		// way to classify errors without this hassle. If you find one, please
		// feel free to share and refactor.
		if err == sql.ErrNoRows {
			return acct, &errorrrs.E{
				ID:  errorrrs.NotFound,
				Msg: err.Error(),
			}
		} else {
			return acct, &errorrrs.E{
				ID:  errorrrs.InternalServerError,
				Msg: err.Error(),
			}
		}
	}

	return acct, err
}

func (ws *ServiceImpl) ListAccounts(req ListAccountsRequest) ([]Account, error) {
	accts, err := ws.Repo.ListAccounts(req)
	if err != nil {
		return accts, &errorrrs.E{
			ID:  errorrrs.InternalServerError,
			Msg: err.Error(),
		}
	}

	return accts, err
}

func (ws *ServiceImpl) CreateAccount(req CreateAccountRequest) (Account, error) {
	acct, err := ws.Repo.CreateAccount(req)
	if err != nil {
		return acct, &errorrrs.E{
			ID:  errorrrs.InternalServerError,
			Msg: err.Error(),
		}
	}

	return acct, err
}

func (ws *ServiceImpl) CreatePayment(req CreatePaymentRequest) (Payment, error) {
	transferReq := CreateTransferRequest{
		From:   req.Self,
		To:     req.To,
		Amount: req.Amount,
	}

	var pymt Payment
	transfer, err := ws.Repo.CreateTransfer(transferReq)
	if err != nil {
		return pymt, &errorrrs.E{
			ID:  errorrrs.InternalServerError,
			Msg: err.Error(),
		}
	}

	pymt.From = &transfer.From
	pymt.To = &transfer.To
	pymt.Amount = transfer.Amount
	pymt.Direction = Outgoing

	return pymt, nil
}

func (ws *ServiceImpl) ListPayments(req ListPaymentsRequest) ([]Payment, error) {
	// Note: we make use of same DB method as `ListTransfers` since `Payment`s
	// are only a `Service` "domain object" and exist in the DB also as `Transfer`s
	transferReq := ListTransfersRequest{
		From: &req.ID,
		To:   &req.ID,
	}

	transfers, err := ws.Repo.ListTransfers(transferReq)
	if err != nil {
		return nil, &errorrrs.E{
			ID:  errorrrs.InternalServerError,
			Msg: err.Error(),
		}
	}

	payments := make([]Payment, len(transfers))
	for i := range transfers {
		t := transfers[i]
		p := Payment{
			Self:   req.ID,
			Amount: t.Amount,
		}
		if req.ID == t.From {
			p.To = &t.To
			p.Direction = Outgoing
		} else {
			p.From = &t.From
			p.Direction = Incoming
		}
		payments[i] = p
	}

	return payments, nil
}

func (ws *ServiceImpl) ListTransfers(req ListTransfersRequest) ([]Transfer, error) {
	// Note: since `From` and `To` work together as a `where... OR` query,
	// it is up to the client component (e.g. web, mobile) to filter on their
	// end for cases where the user explicitly wants a `where... AND`

	trnsfrs, err := ws.Repo.ListTransfers(req)
	if err != nil {
		return nil, &errorrrs.E{
			ID:  errorrrs.InternalServerError,
			Msg: err.Error(),
		}
	}

	return trnsfrs, err
}
