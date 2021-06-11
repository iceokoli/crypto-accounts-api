package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/iceokoli/get-crypto-balance/broker"
)

type TestCryptoPortfolio struct{}

var RawBalanceMockOne func(string) ([]broker.Crypto, bool)

func (p TestCryptoPortfolio) GetBalanceByBroker(broker string) ([]broker.Crypto, bool) {
	return RawBalanceMockOne(broker)
}

var RawBalanceMockTwo func() map[string][]broker.Crypto

func (p TestCryptoPortfolio) GetSegregatedBalance() map[string][]broker.Crypto {
	return RawBalanceMockTwo()
}

var RawBalanceMockThree func() []broker.Crypto

func (p TestCryptoPortfolio) GetAggregatedBalance() []broker.Crypto {
	return RawBalanceMockThree()
}

func TestHandleError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()

	cErr := NewHTTPError(nil, 400, "Error")
	TestServer := Server{}
	TestServer.handleError(cErr, w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Could not read body")
	}

	if string(data) != `{"detail":"Error"}` {
		t.Error("Incorrect response body", string(data))
	}
}

func TestGetLocalBalanceByBrokerEndpoint(t *testing.T) {

	RawBalanceMockOne = func(b string) ([]broker.Crypto, bool) {
		result := []broker.Crypto{
			{Asset: "BTC", Amount: 4723846.89208129},
			{Asset: "LTC", Amount: 4763368.68006011},
		}
		return result, true
	}

	req := httptest.NewRequest(http.MethodGet, "/balance/local/bitstamp", nil)
	w := httptest.NewRecorder()

	testPfolio := TestCryptoPortfolio{}
	TestServer := Server{Router: mux.NewRouter(), Portfolio: testPfolio}
	TestServer.HandleFunc("/balance/local/{broker}", TestServer.GetLocalBalanceByBroker()).Methods("GET")
	TestServer.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Error("Status Code:", res.StatusCode)
	}
}
