package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/iceokoli/get-crypto-balance/broker"
	"github.com/iceokoli/get-crypto-balance/config"
	"github.com/iceokoli/get-crypto-balance/portfolio"
	"github.com/iceokoli/get-crypto-balance/server"
	"github.com/joho/godotenv"
)

const portNumber = ":8080"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Failed to load config", err)
	}

	envVariables, err := loadEnvVariables()
	if err != nil {
		log.Println("Failed to load environment variables", err)
	}

	accounts := broker.New(cfg, envVariables)
	pfolio := portfolio.New(accounts)
	srv := server.New(pfolio, envVariables)

	log.Printf("Starting Server on port %s\n", portNumber)
	if err := http.ListenAndServe(portNumber, srv); err != nil {
		log.Println(err)
	}

}

func loadEnvVariables() (map[string]string, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("Could not find .env file")
		log.Println("Searching for environment variables")
	}

	var ok bool

	result := map[string]string{}
	result["BITSTAMP_ID"], ok = os.LookupEnv("BITSTAMP_ID")
	result["BITSTAMP_KEY"], ok = os.LookupEnv("BITSTAMP_KEY")
	result["BITSTAMP_SECRET"], ok = os.LookupEnv("BITSTAMP_SECRET")

	result["BINANCE_ID"], ok = os.LookupEnv("BINANCE_ID")
	result["BINANCE_KEY"], ok = os.LookupEnv("BINANCE_KEY")
	result["BINANCE_SECRET"], ok = os.LookupEnv("BINANCE_SECRET")

	result["SERVER_AUTH_KEY"], ok = os.LookupEnv("SERVER_AUTH_KEY")
	result["SERVER_AUTH_SECRET"], ok = os.LookupEnv("SERVER_AUTH_SECRET")

	if !ok {
		return nil, errors.New("Missing environment variable")
	}
	return result, nil
}
