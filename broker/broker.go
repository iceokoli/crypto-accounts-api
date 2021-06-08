package broker //model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/iceokoli/get-crypto-balance/config"
)

type Crypto struct {
	Asset  string
	Amount float64
}

func (c *Crypto) Add(other Crypto) error {

	if other.Asset != c.Asset {
		return errors.New("Cannot add together different crypto assets/currencies")
	}
	c.Amount += other.Amount

	return nil
}

type CryptoAccount interface {
	GetBalance() []Crypto
}

func GenerateSignature(msg string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

func New(cfg config.Config, env map[string]string) map[string]CryptoAccount {

	bitstamp := BitstampAccount{
		CustomerID: env["BITSTAMP_ID"],
		Key:        env["BITSTAMP_KEY"],
		Secret:     []byte(env["BITSTAMP_SECRET"]),
		URL:        cfg.Bitstamp.URL,
		Endpoints:  cfg.Bitstamp.Endpoints,
	}

	binance := BinanceAccount{
		Key:       env["BINANCE_KEY"],
		Secret:    []byte(env["BINANCE_SECRET"]),
		URL:       cfg.Binance.URL,
		Endpoints: cfg.Binance.Endpoints,
	}
	accounts := map[string]CryptoAccount{
		"bitstamp": bitstamp,
		"binance":  binance,
	}

	return accounts
}
