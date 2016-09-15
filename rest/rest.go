package rest

import (
	"fmt"
	"net/http"
	"os"

	rice "github.com/GeertJohan/go.rice"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/grindlemire/WellsFarGO/unifier"
)

// Service configures the rest endpoints
type Service struct {
	box  *rice.Box
	port int
	U    *unifier.Unifier
}

// NewRestService creates a new REST service object
func NewRestService(port int, u *unifier.Unifier) (s *Service) {
	box, err := rice.FindBox("build")
	if err != nil {
		log.Critical("Cannot find webserver build directory. Did you build the webserver (npm run buildDev)?")
		return
	}
	return &Service{
		box:  box,
		port: port,
		U:    u,
	}
}

// Start starts the webserver
func (s *Service) Start() {
	go s.run()
}

func (s *Service) run() {

	r := mux.NewRouter()

	files := http.FileServer(s.box.HTTPBox())
	r.Methods("GET").Path("/").HandlerFunc(s.GetHomeHandler)
	r.Methods("POST").Path("/transactions/range").HandlerFunc(s.TransactionsInRangeHandler)
	r.Methods("POST").Path("/transactions/day").HandlerFunc(s.TransactionsInDateHandler)
	r.Methods("POST").Path("/transactions/amount").HandlerFunc(s.TransactionsAmountHandler)
	r.Methods("POST").Path("/transactions/location").HandlerFunc(s.TransactionsLocationHandler)

	r.PathPrefix("/").Handler(files)

	log.Infof("Starting a new webserver on %d\n", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", s.port), r)
	// err := SSL.ListenAndServeTLS(fmt.Sprintf(":%d", s.port), "server.crt", "server.key", r)
	if nil != err {
		log.Critical("Failed to listen on port ", 3000, " err: ", err)
		os.Exit(2)
	}
}

// GetHomeHandler Handles the home page request
func (s *Service) GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	contentString, err := s.box.Bytes("static/index.html")
	if err != nil {
		log.Error("Error loading file: ", err)
	}
	w.Write(contentString)
	return
}

// Close satisifies the closer interface to allow for runnable instance
func (s *Service) Close() error {
	log.Info("Closed")
	return nil
}
