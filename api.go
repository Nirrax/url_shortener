package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func writeJson(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// Defers error for handler functions
func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If handler function has an error, handle the error
		if err := f(w, r); err != nil {
			writeJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type ApiServer struct {
	address string
}

// Constructor
func NewApiServer(address string) *ApiServer {
	return &ApiServer{
		address: address,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/{url}", makeHttpHandleFunc(s.handleGetUrl)).Methods("GET")
	router.HandleFunc("/{url}", makeHttpHandleFunc(s.handleCreateUrl)).Methods("POST")
	router.HandleFunc("/{url}", makeHttpHandleFunc(s.handleDeleteUrl)).Methods("DELETE")
	router.HandleFunc("/info/{url}", makeHttpHandleFunc(s.handleGetUrlInfo)).Methods("GET")

	log.Println("Running on port: " + s.address)

	http.ListenAndServe(s.address, router)
}

func (s *ApiServer) handleGetUrl(w http.ResponseWriter, r *http.Request) error {
	return writeJson(w, http.StatusOK, "Get works")
}

func (s *ApiServer) handleGetUrlInfo(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleCreateUrl(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleDeleteUrl(w http.ResponseWriter, r *http.Request) error {
	return nil
}
