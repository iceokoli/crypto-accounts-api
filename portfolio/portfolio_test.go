package portfolio

import (
	"testing"

	"github.com/iceokoli/get-crypto-balance/broker"
)

type MockBroker struct{}

func (m MockBroker) GetBalance() []broker.Crypto {
	return []broker.Crypto{
		{Asset: "BTC", Amount: 10},
		{Asset: "ETH", Amount: 10},
	}
}

var accounts = map[string]broker.CryptoAccount{
	"broker1": MockBroker{},
	"broker2": MockBroker{},
}

func TestGetSegregatedBalance(t *testing.T) {

	var pfolio = MyCryptoPortfolio{}
	for k, v := range accounts {
		pfolio.AddAccount(k, v)
	}

	actual := pfolio.GetSegregatedBalance()
	expected := map[string][]broker.Crypto{
		"broker1": {
			{Asset: "BTC", Amount: 10},
			{Asset: "ETH", Amount: 10},
		},
		"broker2": {
			{Asset: "BTC", Amount: 10},
			{Asset: "ETH", Amount: 10},
		},
	}

	lenActual := len(actual)
	lenExpected := len(expected)
	if lenActual != lenExpected {
		t.Errorf("Expected %d assets got %d", lenActual, lenExpected)
	}

	for key := range expected {

		a, e := actual[key], expected[key]

		for i := 0; i < 2; i++ {
			assetCorrect := a[i].Asset == e[i].Asset
			amountCorrect := a[i].Amount == e[i].Amount

			if !assetCorrect || !amountCorrect {
				t.Error("Failing to parse crypto assets correctly")
			}
		}
	}

}

func TestGetAggregatedBalance(t *testing.T) {

	var pfolio = MyCryptoPortfolio{}
	for k, v := range accounts {
		pfolio.AddAccount(k, v)
	}

	actual := pfolio.GetAggregatedBalance()
	expected := []broker.Crypto{
		{Asset: "BTC", Amount: 20},
		{Asset: "ETH", Amount: 20},
	}

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
