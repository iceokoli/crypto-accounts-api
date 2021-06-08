package portfolio

import (
	"sync"

	"github.com/iceokoli/get-crypto-balance/broker"
)

type Portfolio struct {
	Accounts map[string]broker.CryptoAccount
}

func (p Portfolio) GetSegregatedBalance() map[string][]broker.Crypto {

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

func (p Portfolio) GetAggregatedBalance() []broker.Crypto {

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

func New(accounts map[string]broker.CryptoAccount) Portfolio {
	return Portfolio{Accounts: accounts}
}
