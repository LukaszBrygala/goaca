package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academy/gateways"

	"github.com/gorilla/mux"
)

type GatawaysRepository interface {
	Find(gatewayID, brokerType string) (gateways.Gateway, error)
}

type Server struct {
	gateways GatawaysRepository
	router   *mux.Router
}

func New(gateways GatawaysRepository) Server {
	return Server{gateways, mux.NewRouter()}
}

func (s Server) setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s Server) info() http.HandlerFunc {
	info := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{"goaca-server", "0.0.1"}

	bs, _ := json.Marshal(info)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(bs)
	}
}

func (s Server) getGateway() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gatewayID := mux.Vars(r)["gatewayID"]
		brokerType := r.FormValue("brokerType")

		gateway, err := s.gateways.Find(gatewayID, brokerType)
		if err == gateways.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(gateway)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(json)
	}
}

func (s Server) routes() {
	s.router.Use(s.setJSONContentType)

	s.router.HandleFunc("/v1", s.info()).Methods("GET")

	gateways := s.router.PathPrefix("/v1/gateways/{gatewayID}").Subrouter()
	gateways.HandleFunc("", s.getGateway()).Methods("GET")
}

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s Server) Run(port int) error {
	s.routes()
	return http.ListenAndServe(fmt.Sprintf(":%v", port), s)
}
