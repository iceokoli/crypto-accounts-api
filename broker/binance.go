package broker //model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BinanceAccount struct {
	Key       string
	Secret    []byte
	URL       string
	Endpoints map[string]string
}

func (b BinanceAccount) generateURL(endpoint, msg, signature string) string {
	return fmt.Sprintf("%s%s?%s&signature=%s", b.URL, endpoint, msg, signature)
}

func (b BinanceAccount) retrieveRawBalance() ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)

	params := url.Values{}
	params.Add("timestamp", timestamp)

	msg := params.Encode()
	signature := GenerateSignature(msg, b.Secret)

	apiUrl := b.generateURL(b.Endpoints["balance"], msg, signature)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-MBX-APIKEY", b.Key)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responceBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responceBody, nil
}

func (b BinanceAccount) formatBalance(raw []byte) []Crypto {
	var staging1 map[string]json.RawMessage
	json.Unmarshal(raw, &staging1)

	var staging2 []map[string]string
	json.Unmarshal(staging1["balances"], &staging2)

	formattedBalance := []Crypto{}
	for _, value := range staging2 {
		assetBalance, err := strconv.ParseFloat(value["free"], 10)
		if err != nil {
			log.Println(err)
		}
		if assetBalance == 0 {
			continue
		}
		newAsset := Crypto{Asset: value["asset"], Amount: assetBalance}
		formattedBalance = append(formattedBalance, newAsset)
	}
	return formattedBalance
}

func (b BinanceAccount) GetBalance() []Crypto {
	rawBalance, err := b.retrieveRawBalance()
	if err != nil {
		log.Printf("Failed to retrieve raw binance balance %s\n", err)
	}
	cleanBalance := b.formatBalance(rawBalance)

	return cleanBalance
}
