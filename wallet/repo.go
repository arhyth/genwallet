package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type CreateTransferRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type Repository interface {
	ListAccounts(ListAccountsRequest) ([]Account, error)
	GetAccount(GetAccountRequest) (Account, error)
	CreateAccount(CreateAccountRequest) (Account, error)
	CreateTransfer(CreateTransferRequest) (Transfer, error)
	ListTransfers(ListTransfersRequest) ([]Transfer, error)
}

var _ Repository = (*Repo)(nil)

type Repo struct {
	DB *sql.DB

	// Note: here we make prepared statements for each repository method
	// and use sync.Once/s to lazily initialize the statements.
	// `Transfer` related queries are purposely omitted since those
	// are tricky enough even without these.
	//
	// To be honest, this is an "optimization" that I have not done on any
	// production services because it easily becomes tedious and burdensome
	// to maintain especially if the API is under heavy development;
	// but since this is only a test repo/service, might as well experiment ;)
	createAcctOnce *sync.Once
	createAcctStmt *sql.Stmt

	listAcctsOnce    *sync.Once
	listAcctsStmt    *sql.Stmt
	listAcctsCurStmt *sql.Stmt

	getAcctOnce *sync.Once
	getAcctStmt *sql.Stmt
}

func NewRepo(dsn string) (*Repo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &Repo{
		DB: db,
	}

	repo.createAcctOnce = &sync.Once{}
	repo.listAcctsOnce = &sync.Once{}
	repo.getAcctOnce = &sync.Once{}

	return repo, nil
}

func (r *Repo) ListAccounts(req ListAccountsRequest) ([]Account, error) {
	// TODO: add pagination

	r.listAcctsOnce.Do(func() {
		var err error
		listAccts := `SELECT id, balance, currency, created_at, updated_at FROM accounts;`
		r.listAcctsStmt, err = r.DB.Prepare(listAccts)
		if err != nil {
			panic(err.Error())
		}
		listAcctsWithCur := `SELECT id, balance, currency, created_at, updated_at
			FROM accounts WHERE currency = $1;`
		r.listAcctsCurStmt, err = r.DB.Prepare(listAcctsWithCur)
		if err != nil {
			panic(err)
		}
	})

	var (
		rows *sql.Rows
		err  error
	)
	if req.Currency != nil {
		rows, err = r.listAcctsCurStmt.Query(req.Currency)
	} else {
		rows, err = r.listAcctsStmt.Query()
	}
	if err != nil {
		return nil, err
	}

	var accts []Account
	for rows.Next() {
		var acct Account
		if err := rows.Scan(&acct.ID, &acct.Balance, &acct.Currency,
			&acct.CreatedAt, &acct.UpdatedAt); err != nil {
			return nil, err
		}

		accts = append(accts, acct)
	}

	return accts, err
}

func (r *Repo) GetAccount(req GetAccountRequest) (Account, error) {
	r.getAcctOnce.Do(func() {
		var err error
		getAcct := `SELECT id, balance, currency, created_at, updated_at
		FROM accounts WHERE id = $1;`
		r.getAcctStmt, err = r.DB.Prepare(getAcct)
		if err != nil {
			panic(err.Error())
		}
	})

	var acct Account
	err := r.getAcctStmt.QueryRow(req.ID).
		Scan(&acct.ID, &acct.Balance, &acct.Currency, &acct.CreatedAt, &acct.UpdatedAt)
	if err != nil {
		return acct, err
	}

	return acct, nil
}

func (r *Repo) CreateAccount(req CreateAccountRequest) (Account, error) {
	r.createAcctOnce.Do(func() {
		var err error
		createAcct := `INSERT INTO accounts (id, balance, currency)
		VALUES ($1, $2, $3)
		RETURNING id, balance, currency, created_at, updated_at;`
		r.createAcctStmt, err = r.DB.Prepare(createAcct)
		if err != nil {
			panic(err.Error())
		}
	})

	var acct Account
	err := r.createAcctStmt.QueryRow(req.ID, req.InitAmt, req.Currency).
		Scan(&acct.ID, &acct.Balance, &acct.Currency, &acct.CreatedAt, &acct.UpdatedAt)
	if err != nil {
		return acct, err
	}

	return acct, nil
}

func (r *Repo) CreateTransfer(req CreateTransferRequest) (Transfer, error) {
	var (
		trnsfr         Transfer
		fromCur, toCur string
		fromBal, toBal float64
		rbErr          error
	)
	ctx := context.Background()
	txOptns := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	tx, err := r.DB.BeginTx(ctx, txOptns)
	if err != nil {
		return trnsfr, err
	}

	defer func() {
		// catch if rollback fails
		if rbErr != nil {
			log.Err(rbErr).Msg("repo.CreateTransfer: txn rollback fail")
		}
	}()
	err = tx.QueryRow(`SELECT currency, balance FROM accounts where id = $1;`, req.From).Scan(&fromCur, &fromBal)
	if err != nil {
		rbErr = tx.Rollback()
		return trnsfr, err
	}

	err = tx.QueryRow(`SELECT currency, balance FROM accounts where id = $1;`, req.To).Scan(&toCur, &toBal)
	if err != nil {
		rbErr = tx.Rollback()
		return trnsfr, err
	}

	if toCur != fromCur {
		rbErr = tx.Rollback()
		return trnsfr, errors.New("wallet accounts are not of same currency")
	}

	if fromBal < req.Amount {
		rbErr = tx.Rollback()
		return trnsfr, errors.New("existing balance less than requested transfer amount")
	}

	_, err = tx.Exec(`UPDATE accounts
	SET (balance, updated_at) = ($1, now())
	WHERE id = $2;`, fromBal-req.Amount, req.From)
	if err != nil {
		rbErr = tx.Rollback()
		return trnsfr, err
	}

	_, err = tx.Exec(`UPDATE accounts
	SET (balance, updated_at) = ($1, now())
	WHERE id = $2;`, toBal+req.Amount, req.To)
	if err != nil {
		rbErr = tx.Rollback()
		return trnsfr, err
	}

	err = tx.QueryRow(`INSERT INTO transfers ("from", "to", currency, amount)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at;`, req.From, req.To, fromCur, req.Amount).
		Scan(&trnsfr.ID, &trnsfr.CreatedAt)
	if err != nil {
		rbErr = tx.Rollback()
		return trnsfr, err
	}
	if err = tx.Commit(); err != nil {
		rbErr = tx.Rollback()
		return trnsfr, err
	}
	trnsfr.Amount = req.Amount
	trnsfr.From = req.From
	trnsfr.To = req.To
	trnsfr.Currency = fromCur

	return trnsfr, nil
}

func (r *Repo) ListTransfers(req ListTransfersRequest) ([]Transfer, error) {
	// TODO: add pagination

	listTrnsfrBase := `SELECT id, "from", "to", amount, currency, created_at
	FROM transfers %v;`
	var whereClause string

	// TODO: use string builder, if possible, to optimize
	if req.Currency != nil {
		whereClause = fmt.Sprintf(`WHERE currency = '%v'`, req.Currency)

		if req.From == nil || req.To == nil {
			if req.From != nil {
				whereClause = fmt.Sprintf(` AND WHERE "from" = '%v'`, *req.From)
			} else if req.To != nil {
				whereClause = fmt.Sprintf(` AND WHERE "to" = '%v'`, *req.From)
			}
		} else {
			whereClause = fmt.Sprintf(` AND WHERE "from" = '%v' OR "to" = '%v'`, *req.From, *req.To)
		}
	} else {
		if req.From == nil || req.To == nil {
			if req.From != nil {
				whereClause = fmt.Sprintf(`WHERE "from" = '%v'`, *req.From)
			} else if req.To != nil {
				whereClause = fmt.Sprintf(`WHERE "to" = '%v'`, *req.From)
			}
		} else {
			whereClause = fmt.Sprintf(`WHERE "from" = '%v' OR "to" = '%v'`, *req.From, *req.To)
		}
	}

	query := fmt.Sprintf(listTrnsfrBase, whereClause)
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var transfers []Transfer
	for rows.Next() {
		var trnsfr Transfer
		if err := rows.Scan(&trnsfr.ID,
			&trnsfr.From,
			&trnsfr.To,
			&trnsfr.Amount,
			&trnsfr.Currency,
			&trnsfr.CreatedAt); err != nil {

			return nil, err
		}

		transfers = append(transfers, trnsfr)
	}

	return transfers, nil
}
