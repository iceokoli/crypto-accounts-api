package broker //model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type bitstampAuth struct {
	Timestamp int64
	Nonce     string
	Signature string
}

type BitstampAccount struct {
	CustomerID string
	Key        string
	Secret     []byte
	URL        string
	Endpoints  map[string]string
}

func (b BitstampAccount) authenticate(endpoint string, version string) bitstampAuth {

	out := bitstampAuth{}
	out.Timestamp = time.Now().Unix() * 1000
	out.Nonce = uuid.New().String()

	msg := fmt.Sprintf(
		"BITSTAMP %sPOSTwww.bitstamp.net/api%s%s%s%s",
		b.Key,
		endpoint,
		out.Nonce,
		strconv.FormatInt(out.Timestamp, 10),
		version,
	)
	out.Signature = GenerateSignature(msg, b.Secret)

	return out
}

func (b BitstampAccount) addHeaders(r *http.Request, auth bitstampAuth, version string) {

	r.Header.Add("X-Auth", "BITSTAMP "+b.Key)
	r.Header.Add("X-Auth-Signature", auth.Signature)
	r.Header.Add("X-Auth-Nonce", auth.Nonce)
	r.Header.Add("X-Auth-Timestamp", strconv.FormatInt(auth.Timestamp, 10))
	r.Header.Add("X-Auth-Version", version)
	r.Header.Add("Content-Type", "")
}

func (b BitstampAccount) retrieveRawBalance() []byte {
	version := "v2"
	auth := b.authenticate(b.Endpoints["balance"], version)
	apiUrl := b.URL + b.Endpoints["balance"]

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		log.Println(err)
	}
	b.addHeaders(req, auth, version)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	responceBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	return responceBody
}

func (b BitstampAccount) formatBalance(raw []byte) []Crypto {
	var staging map[string]string
	json.Unmarshal(raw, &staging)

	formattedBalance := []Crypto{}
	for key, value := range staging {
		if !strings.Contains(key, "balance") {
			continue
		}

		newValue, err := strconv.ParseFloat(value, 10)
		if err != nil {
			log.Println(err)
		}
		if newValue == 0 {
			continue
		}

		sym := strings.ToUpper(strings.Split(key, "_")[0])
		formattedBalance = append(formattedBalance, Crypto{Asset: sym, Amount: newValue})
	}

	return formattedBalance
}

func (b BitstampAccount) GetBalance() []Crypto {

	raw := b.retrieveRawBalance()
	cleanBalance := b.formatBalance(raw)

	return cleanBalance
}
