package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func New() Server {
	return Server{mux.NewRouter()}
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

func (s Server) deviceFeature() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		featureName := mux.Vars(r)["featureName"]
		w.Write([]byte(fmt.Sprintf("received request for feature %v", featureName)))
	}
}

func (s Server) routes() {
	s.router.Use(s.setJSONContentType)

	s.router.HandleFunc("/v1", s.info()).Methods("GET")

	gateways := s.router.PathPrefix("/v1/gateways/{gatewayID}").Subrouter()
	gateways.HandleFunc("/devices/{deviceID}/features/{featureName}", s.deviceFeature()).Methods("GET")
}

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s Server) Run(port int) error {
	s.routes()
	return http.ListenAndServe(fmt.Sprintf(":%v", port), s)
}
