package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrParamNotFound   = errors.New("param not found")
	ErrParsingJsonBody = errors.New("failed to parse json body")
	ErrEncodingJson    = errors.New("failed to encode json body")
)

func ExtractUrl(r *http.Request) (string, error) {
	return extractParam(r, "url")
}

func SendJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		SendJsonError(w, http.StatusInternalServerError, ErrEncodingJson)
	}
}

func SendJsonError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := map[string]string{"error": err.Error()}
	SendJsonResponse(w, statusCode, errorResponse)
}

func DecodeJSONBody(r *http.Request, v any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(v)
	if err != nil {
		return ErrParsingJsonBody
	}
	return nil
}

func extractParam(r *http.Request, name string) (string, error) {
	param := r.PathValue(name)
	if param == "" {
		return "", ErrParamNotFound
	}

	return param, nil
}
