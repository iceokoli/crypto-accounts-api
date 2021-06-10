package broker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockBitstampAccount struct {
	*BitstampAccount
}

var BitstampRawBalanceMock func() []byte

func (m *MockBitstampAccount) retrieveRawBalance() []byte {
	return BitstampRawBalanceMock()
}

var stringBitstampResponse = `{
	"btc_balance": "4723846.89208129",
	"ltc_balance": "4763368.68006011",
	"bch_balance": "0.00000000", 
	"bch_reserved": "0.00000000"
}`

func TestRetrieveRawBitstampBalance(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(stringBitstampResponse))
	}))
	b := BitstampAccount{
		CustomerID: "test",
		Key:        "test",
		Secret:     []byte("secret"),
		URL:        srv.URL,
		Endpoints:  map[string]string{"bitstamp": "/balance"},
	}
	_, err := b.retrieveRawBalance()
	if err != nil {
		t.Error(err)
	}
}

func TestFormatBitstampBalance(t *testing.T) {

	BitstampRawBalanceMock = func() []byte {
		return []byte(stringBitstampResponse)
	}

	expected := []Crypto{
		{Asset: "BTC", Amount: 4723846.89208129},
		{Asset: "LTC", Amount: 4763368.68006011},
	}

	m := MockBitstampAccount{BitstampAccount: &BitstampAccount{}}
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
