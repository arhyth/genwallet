//go:build integration

package wallet_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arhyth/genwallet/config"
	"github.com/arhyth/genwallet/wallet"
)

var repo wallet.Repository

func TestMain(m *testing.M) {
	cfg, err := config.GetAPIConfig()
	if err != nil {
		panic(err.Error())
	}
	repo, err = wallet.NewRepo(cfg.DBConnStr)
	if err != nil {
		panic(err.Error())
	}
	os.Exit(m.Run())
}

func TestRepoCreateAccount(t *testing.T) {
	reqrd := require.New(t)
	as := assert.New(t)

	createReq := wallet.CreateAccountRequest{
		ID:       "alice-123",
		Currency: "USD",
		InitAmt:  800.0,
	}
	acct, err := repo.CreateAccount(createReq)
	reqrd.Nil(err)
	as.Equal(createReq.ID, acct.ID)
	as.Equal(createReq.InitAmt, acct.Balance)
	as.Equal(createReq.Currency, acct.Currency)
}
