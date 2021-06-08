package broker

import "testing"

func TestCryptoAdd(t *testing.T) {

	a := Crypto{Asset: "BTC", Amount: 10}
	b := Crypto{Asset: "BTC", Amount: 20}

	a.Add(b)
	expected := 30.0
	if a.Amount != expected {
		t.Errorf("Expected %f got %f", expected, a.Amount)
	}

	c := Crypto{Asset: "BTC", Amount: 10}
	d := Crypto{Asset: "ETH", Amount: 20}

	if err := c.Add(d); err == nil {
		t.Error("Summing assets of different type together")
	}

}
