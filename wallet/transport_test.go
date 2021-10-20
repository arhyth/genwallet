package wallet_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arhyth/genwallet/errorrrs"
	"github.com/arhyth/genwallet/wallet"
	MOCKWALLET "github.com/arhyth/genwallet/wallet/mock"
	httptransport "github.com/go-kit/kit/transport/http"
)

func TestHTTPGetWallet(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		reqrd := require.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		walletSvc := &wallet.ValidationMiddleware{
			Next: &wallet.ServiceImpl{
				Repo: repo,
			},
		}

		serverErrcoder := httptransport.ServerErrorEncoder(errorrrs.GokitErrorEncoder)
		walletGetHandler := httptransport.NewServer(
			wallet.MakeWalletGetEndpt(walletSvc),
			wallet.DecodeHTTPGetAccountReq,
			wallet.EncodeJSONResponse,
			serverErrcoder)
		w := httptest.NewRecorder()

		account := wallet.Account{
			ID:       "sato-91011",
			Balance:  1000.0,
			Currency: "CNY",
		}
		getReq := wallet.GetAccountRequest{ID: account.ID}
		req, err := http.NewRequest("GET", fmt.Sprintf(`/wallets/%v`, account.ID), nil)
		reqrd.Nil(err)

		repo.EXPECT().
			GetAccount(gomock.AssignableToTypeOf(getReq)).
			Return(account, nil).
			Times(1)

		walletGetHandler.ServeHTTP(w, req)

		bits, err := io.ReadAll(w.Result().Body)
		reqrd.Nil(err)
		var resp wallet.Account
		err = json.Unmarshal(bits, &resp)
		reqrd.Nil(err)

		as.Equal(account.ID, resp.ID)
		as.Equal(account.Balance, resp.Balance)
		as.Equal(account.Currency, resp.Currency)
	})
}

func TestHTTPListWallets(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		reqrd := require.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		walletSvc := &wallet.ValidationMiddleware{
			Next: &wallet.ServiceImpl{
				Repo: repo,
			},
		}

		serverErrcoder := httptransport.ServerErrorEncoder(errorrrs.GokitErrorEncoder)
		walletListHandler := httptransport.NewServer(
			wallet.MakeWalletListEndpt(walletSvc),
			wallet.DecodeHTTPListAccountsReq,
			wallet.EncodeJSONResponse,
			serverErrcoder)
		w := httptest.NewRecorder()

		accounts := []wallet.Account{
			{
				ID:       "sato-91011",
				Balance:  6000.0,
				Currency: "CNY",
			},
			{
				ID:       "fan-1234",
				Balance:  3000.0,
				Currency: "CNY",
			},
			{
				ID:       "hao-91011",
				Balance:  5000.0,
				Currency: "CNY",
			},
		}
		cur := "CNY"
		listReq := wallet.ListAccountsRequest{Currency: &cur}
		req, err := http.NewRequest("GET", fmt.Sprintf(`/wallets?%v`, cur), nil)
		reqrd.Nil(err)

		repo.EXPECT().
			ListAccounts(gomock.AssignableToTypeOf(listReq)).
			Return(accounts, nil).
			Times(1)

		walletListHandler.ServeHTTP(w, req)

		bits, err := io.ReadAll(w.Result().Body)
		reqrd.Nil(err)
		var resp []wallet.Account
		err = json.Unmarshal(bits, &resp)
		reqrd.Nil(err)

		as.Len(resp, len(accounts))
		as.Equal(accounts[0].ID, resp[0].ID)
		as.Equal(accounts[0].Balance, resp[0].Balance)
		as.Equal(accounts[0].Currency, resp[0].Currency)
	})
}

func TestHTTPListPayments(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		reqrd := require.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		walletSvc := &wallet.ValidationMiddleware{
			Next: &wallet.ServiceImpl{
				Repo: repo,
			},
		}

		serverErrcoder := httptransport.ServerErrorEncoder(errorrrs.GokitErrorEncoder)
		walletPaymentsIndexHandler := httptransport.NewServer(
			wallet.MakePaymentsIndexEndpt(walletSvc),
			wallet.DecodeHTTPListPaymentsReq,
			wallet.EncodeJSONResponse,
			serverErrcoder)
		w := httptest.NewRecorder()

		acctID := "bob-888"
		transfers := []wallet.Transfer{
			{
				From:     "sato-91011",
				To:       acctID,
				Amount:   50.0,
				Currency: "USD",
			},
			{
				From:     acctID,
				To:       "fan-1234",
				Amount:   300.0,
				Currency: "USD",
			},
			{
				From:     acctID,
				To:       "hao-91011",
				Amount:   100.0,
				Currency: "USD",
			},
		}

		listReq := wallet.ListTransfersRequest{
			From: &acctID,
			To:   &acctID,
		}
		req, err := http.NewRequest("GET", fmt.Sprintf(`/wallets/%v/payments`, acctID), nil)
		reqrd.Nil(err)

		repo.EXPECT().
			ListTransfers(gomock.AssignableToTypeOf(listReq)).
			Return(transfers, nil).
			Times(1)

		walletPaymentsIndexHandler.ServeHTTP(w, req)

		bits, err := io.ReadAll(w.Result().Body)
		reqrd.Nil(err)

		var resp []wallet.Payment
		err = json.Unmarshal(bits, &resp)
		reqrd.Nil(err)

		as.Len(resp, len(transfers))
		as.Equal(transfers[0].From, *resp[0].From)
		as.Equal(transfers[0].To, resp[0].Self)
		as.Equal(wallet.Incoming, resp[0].Direction)
		as.Equal(transfers[1].To, *resp[1].To)
		as.Equal(transfers[1].From, resp[1].Self)
		as.Equal(wallet.Outgoing, resp[1].Direction)
		as.Equal(transfers[2].To, *resp[2].To)
		as.Equal(transfers[2].From, resp[2].Self)
		as.Equal(wallet.Outgoing, resp[2].Direction)
	})
}

func TestHTTPCreatePayment(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		as := assert.New(tt)
		reqrd := require.New(tt)
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()
		repo := MOCKWALLET.NewMockRepository(ctrl)

		walletSvc := &wallet.ValidationMiddleware{
			Next: &wallet.ServiceImpl{
				Repo: repo,
			},
		}

		serverErrcoder := httptransport.ServerErrorEncoder(errorrrs.GokitErrorEncoder)
		walletCreatePaymentsHandler := httptransport.NewServer(
			wallet.MakePaymentsPostEndpt(walletSvc),
			wallet.DecodeHTTPPostPaymentsReq,
			wallet.EncodeJSONResponse,
			serverErrcoder)
		w := httptest.NewRecorder()

		acctID := "bob-888"

		payReq := wallet.CreatePaymentRequest{
			Self:   acctID,
			To:     "hao-91011",
			Amount: 100.0,
		}
		create := wallet.CreateTransferRequest{
			From:   acctID,
			To:     payReq.To,
			Amount: payReq.Amount,
		}

		trnsfr := wallet.Transfer{
			From:   acctID,
			To:     payReq.To,
			Amount: payReq.Amount,
		}

		reqBits, err := json.Marshal(payReq)
		reqrd.Nil(err)
		req, err := http.NewRequest("POST", fmt.Sprintf(`/wallets/%v/payments`, acctID), bytes.NewReader(reqBits))
		reqrd.Nil(err)

		repo.EXPECT().
			CreateTransfer(create).
			Return(trnsfr, nil).
			Times(1)

		walletCreatePaymentsHandler.ServeHTTP(w, req)

		bits, err := io.ReadAll(w.Result().Body)
		reqrd.Nil(err)

		var resp wallet.Payment
		err = json.Unmarshal(bits, &resp)
		reqrd.Nil(err)

		as.Equal(payReq.Self, resp.Self)
		as.Equal(payReq.To, *resp.To)
		as.Equal(payReq.Amount, resp.Amount)
		as.Equal(wallet.Outgoing, resp.Direction)
	})
}
