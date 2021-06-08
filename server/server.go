package server //controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iceokoli/get-crypto-balance/broker"
	"github.com/iceokoli/get-crypto-balance/portfolio"
)

type authorised struct {
	user   string
	secret []byte
}

var perms *authorised

func New(pfolio portfolio.Portfolio, env map[string]string) *Server {
	perms = &authorised{user: env["SERVER_AUTH_KEY"], secret: []byte(env["SERVER_AUTH_SECRET"])}
	s := &Server{Router: mux.NewRouter(), Portfolio: pfolio}
	s.setUpMiddleWare()
	s.setUpRoutes()
	return s
}

type Server struct {
	*mux.Router
	Portfolio portfolio.Portfolio
}

func (s *Server) setUpMiddleWare() {
	s.Use(loggingMiddleWare)
	s.Use(authMiddleWare)
}

func (s *Server) setUpRoutes() {
	s.HandleFunc("/balance/local", s.GetLocalBalance()).Methods("GET")
	s.HandleFunc("/balance/local/total", s.GetTotalLocalBalance()).Methods("GET")
	s.HandleFunc("/balance/local/{broker}", s.GetLocalBalanceByBroker()).Methods("GET")
}

func (s *Server) handleError(cErr ClientError, w http.ResponseWriter, r *http.Request) {
	body, _ := cErr.ResponseBody()
	status, headers := cErr.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)

}

func (s *Server) GetLocalBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalBalance := s.Portfolio.GetSegregatedBalance()
		w.Header().Set("Content-Type", "application/json")

		balanceJSON, err := json.Marshal(totalBalance)
		if err != nil {
			log.Println(err)
		}
		w.Write(balanceJSON)
	}
}

func (s *Server) GetTotalLocalBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalBalance := s.Portfolio.GetAggregatedBalance()
		w.Header().Set("Content-Type", "application/json")

		balanceJSON, err := json.Marshal(totalBalance)
		if err != nil {
			log.Println(err)
		}
		w.Write(balanceJSON)
	}
}

func (s *Server) GetLocalBalanceByBroker() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		var balanceApi broker.CryptoAccount
		var ok bool
		if balanceApi, ok = s.Portfolio.Accounts[params["broker"]]; !ok {
			clientError := NewHTTPError(nil, 400, "Invalid broker: Choose a valid crypto broker")
			s.handleError(clientError, w, r)
			log.Println(clientError.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		balance := balanceApi.GetBalance()

		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			log.Println(err)
		}
		w.Write(balanceJSON)
	}
}
