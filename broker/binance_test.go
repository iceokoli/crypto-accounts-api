package broker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockBinanceAccount struct {
	*BinanceAccount
}

var RawBalanceMock func() []byte

func (m *MockBinanceAccount) retrieveRawBalance() []byte {
	return RawBalanceMock()
}

var stringResponse = `{
	"makerCommission": 15,
	"takerCommission": 15,
	"buyerCommission": 0,
	"sellerCommission": 0,
	"canTrade": true,
	"canWithdraw": true,
	"canDeposit": true,
	"updateTime": 123456789,
	"accountType": "SPOT",
	"balances": [
	  {
		"asset": "BTC",
		"free": "4723846.89208129",
		"locked": "0.00000000"
	  },
	  {
		"asset": "LTC",
		"free": "4763368.68006011",
		"locked": "0.00000000"
	  }
	],
	"permissions": [
	  "SPOT"
	]
  }`

func TestRetrieveRawBalance(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(stringResponse))
	}))
	defer svr.Close()

	testEndpoints := map[string]string{
		"balance": "/balance",
	}
	api := BinanceAccount{Key: "test", Secret: []byte("test"), URL: svr.URL, Endpoints: testEndpoints}
	_, err := api.retrieveRawBalance()
	if err != nil {
		t.Error(err)
	}
}

func TestFormatBinanceBalance(t *testing.T) {

	RawBalanceMock = func() []byte {
		return []byte(stringResponse)
	}

	expected := []Crypto{
		{Asset: "BTC", Amount: 4723846.89208129},
		{Asset: "LTC", Amount: 4763368.68006011},
	}

	m := MockBinanceAccount{BinanceAccount: &BinanceAccount{}}
	raw := m.retrieveRawBalance()
	actual := m.formatBalance(raw)

	lenActual := len(actual)
	lenExpected := len(expected)
	if lenActual != lenExpected {
		t.Errorf("Expected %d assets got %d", lenActual, lenExpected)
	}

	for i := 0; i < 2; i++ {
		assetCorrect := actual[i].Asset == expected[i].Asset
		amountCorrect := actual[i].Amount == expected[i].Amount

		if !assetCorrect || !amountCorrect {
			t.Error("Failing to parse crypto assets correctly")
		}
	}
}
