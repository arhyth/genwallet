package wallet_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/arhyth/genwallet/wallet"
	MOCKWALLET "github.com/arhyth/genwallet/wallet/mock"
)

// TODO: add test for specific failure cases
// These "happy path" tests serve only as base case that the service at least works :D

func TestListAccounts(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		svc := &wallet.ServiceImpl{
			Repo: repo,
		}
		listReq := wallet.ListAccountsRequest{}
		accounts := []wallet.Account{
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
		repo.EXPECT().
			ListAccounts(gomock.AssignableToTypeOf(listReq)).
			Return(accounts, nil).
			Times(1)

		result, err := svc.ListAccounts(listReq)
		as.Nil(err)
		as.Len(result, len(accounts))
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		svc := &wallet.ServiceImpl{
			Repo: repo,
		}
		getReq := wallet.GetAccountRequest{}
		account := wallet.Account{
			ID:       "sato-91011",
			Balance:  1000.0,
			Currency: "CNY",
		}
		repo.EXPECT().
			GetAccount(gomock.AssignableToTypeOf(getReq)).
			Return(account, nil).
			Times(1)

		result, err := svc.GetAccount(getReq)
		as.Nil(err)
		as.Equal(account.ID, result.ID)
	})
}

func TestCreateAccount(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		svc := &wallet.ServiceImpl{
			Repo: repo,
		}
		createReq := wallet.CreateAccountRequest{}
		now := time.Now().UTC()
		repo.EXPECT().
			CreateAccount(gomock.AssignableToTypeOf(createReq)).
			DoAndReturn(func(r wallet.CreateAccountRequest) (wallet.Account, error) {
				return wallet.Account{
					ID:        r.ID,
					Balance:   r.InitAmt,
					Currency:  r.Currency,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			}).
			Times(1)

		result, err := svc.CreateAccount(createReq)
		as.Nil(err)
		as.Equal(result.ID, createReq.ID)
		as.Equal(result.Balance, createReq.InitAmt)
		as.Equal(result.Currency, createReq.Currency)
		as.Equal(result.CreatedAt, now)
		as.Equal(result.UpdatedAt, now)
	})
}

func TestListPayments(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		svc := &wallet.ServiceImpl{
			Repo: repo,
		}
		listPReq := wallet.ListPaymentsRequest{
			ID: "alice123",
		}
		now := time.Now().UTC()
		listTransferReq := wallet.ListTransfersRequest{
			From: &listPReq.ID,
			To:   &listPReq.ID,
		}
		transfers := []wallet.Transfer{
			{
				From:      "alice123",
				To:        "bob456",
				Amount:    100.0,
				CreatedAt: now.AddDate(0, -1, 0),
			},
			{
				From:      "sato789",
				To:        "alice123",
				Amount:    150.0,
				CreatedAt: now.AddDate(0, 0, -10),
			},
		}
		repo.EXPECT().
			ListTransfers(gomock.AssignableToTypeOf(listTransferReq)).
			Return(transfers, nil).
			Times(1)

		result, err := svc.ListPayments(listPReq)
		as.Nil(err)
		as.Len(result, len(transfers))
		as.Equal(result[0].Self, transfers[0].From)
		as.Equal(*result[0].To, transfers[0].To)
		as.Equal(result[0].Direction, wallet.Outgoing)
		as.Equal(result[1].Direction, wallet.Incoming)
	})
}
