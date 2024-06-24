package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	address  string
	database databaseI
}

// Constructor
func NewApiServer(address string, db databaseI) *ApiServer {
	return &ApiServer{
		address:  address,
		database: db,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHttpHandleFunc(s.handleCreateUrl)).Methods("POST")
	router.HandleFunc("/{url}", makeHttpHandleFunc(s.handleGetUrl)).Methods("GET")
	router.HandleFunc("/{url}", makeHttpHandleFunc(s.handleDeleteUrl)).Methods("DELETE")
	router.HandleFunc("/info/{url}", makeHttpHandleFunc(s.handleGetUrlInfo)).Methods("GET")

	log.Println("Running on port: " + s.address)

	http.ListenAndServe(s.address, router)
}

func (s *ApiServer) handleGetUrl(w http.ResponseWriter, r *http.Request) error {
	param := mux.Vars(r)["url"]
	object, err := s.database.GetUrlByShortUrl(param)

	if err != nil {
		return fmt.Errorf("invalid url: %s", param)
	}

	http.Redirect(w, r, object.LongUrl, http.StatusMovedPermanently)
	return nil
}

func (s *ApiServer) handleGetUrlInfo(w http.ResponseWriter, r *http.Request) error {
	param := mux.Vars(r)["url"]
	object, err := s.database.GetUrlByShortUrl(param)

	if err != nil {
		return fmt.Errorf("invalid url: %v", param)
	}

	return writeJson(w, http.StatusOK, object)
}

func (s *ApiServer) handleCreateUrl(w http.ResponseWriter, r *http.Request) error {
	//url.ParseRequestURI
	newDto := new(UrlDto)

	if err := json.NewDecoder(r.Body).Decode(newDto); err != nil {
		return err
	}

	//check if the url is valid
	if _, err := url.ParseRequestURI(newDto.LongUrl); err != nil {
		return err
	}

	newUrl, err := NewUrl(newDto.LongUrl)

	if err != nil {
		return err
	}

	if err := s.database.CreateUrl(*newUrl); err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, newUrl)
}

func (s *ApiServer) handleDeleteUrl(w http.ResponseWriter, r *http.Request) error {
	param := mux.Vars(r)["url"]
	err := s.database.DeleteUrlByShortUrl(param)

	if err != nil {
		//return err
		return fmt.Errorf("invalid url: %v", param)
	}

	return writeJson(w, http.StatusOK, "Deleted url: "+param)
}

func writeJson(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
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
