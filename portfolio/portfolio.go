package portfolio

import (
	"sync"

	"github.com/iceokoli/get-crypto-balance/broker"
)

type CryptoPortfolio interface {
	GetBalanceByBroker(string) ([]broker.Crypto, bool)
	GetSegregatedBalance() map[string][]broker.Crypto
	GetAggregatedBalance() []broker.Crypto
}

type MyCryptoPortfolio struct {
	Accounts map[string]broker.CryptoAccount
}

func (p MyCryptoPortfolio) GetBalanceByBroker(broker string) ([]broker.Crypto, bool) {
	api, ok := p.Accounts[broker]
	if !ok {
		return nil, ok
	}
	return api.GetBalance(), true
}

func (p MyCryptoPortfolio) GetSegregatedBalance() map[string][]broker.Crypto {

	balance := map[string][]broker.Crypto{}
	var wg sync.WaitGroup

	for brk, api := range p.Accounts {
		wg.Add(1)
		go func(brk string, api broker.CryptoAccount) {
			defer wg.Done()
			balance[brk] = api.GetBalance()
		}(brk, api)
	}

	wg.Wait()

	return balance
}

func (p MyCryptoPortfolio) GetAggregatedBalance() []broker.Crypto {

	balance := p.GetSegregatedBalance()

	library := map[string]int{}
	totalBalance := []broker.Crypto{}

	numAssets := 0
	for _, assets := range balance {

		for _, crypto := range assets {
			location, ok := library[crypto.Asset]

			if ok {
				totalBalance[location].Add(crypto)
				continue
			}

			totalBalance = append(totalBalance, crypto)

			library[crypto.Asset] = numAssets
			numAssets++
		}
	}

	return totalBalance
}

func New(accounts map[string]broker.CryptoAccount) MyCryptoPortfolio {
	return MyCryptoPortfolio{Accounts: accounts}
}
