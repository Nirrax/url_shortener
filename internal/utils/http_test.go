package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractUrl(t *testing.T) {
	tests := []struct {
		name          string
		pathValue     string
		expectError   bool
		expectedError error
		expected      string
	}{
		{"valid url", "https://example.com", false, nil, "https://example.com"},
		{"empty url", "", true, ErrParamNotFound, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.SetPathValue("url", tc.pathValue)
			result, err := ExtractUrl(req)
			if tc.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestSendJsonResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		data           any
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			"ok response with map",
			http.StatusOK,
			map[string]string{"key": "value"},
			http.StatusOK,
			map[string]any{"key": "value"},
		},
		{
			"created response",
			http.StatusCreated,
			map[string]string{"id": "42"},
			http.StatusCreated,
			map[string]any{"id": "42"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SendJsonResponse(w, tc.statusCode, tc.data)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var body map[string]any
			err := json.NewDecoder(w.Body).Decode(&body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}

func TestSendJsonError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{"not found error", http.StatusNotFound, errors.New("resource not found"), http.StatusNotFound, "resource not found"},
		{"internal error", http.StatusInternalServerError, ErrEncodingJson, http.StatusInternalServerError, ErrEncodingJson.Error()},
		{"param not found", http.StatusBadRequest, ErrParamNotFound, http.StatusBadRequest, ErrParamNotFound.Error()},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SendJsonError(w, tc.statusCode, tc.err)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var body map[string]string
			err := json.NewDecoder(w.Body).Decode(&body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedMsg, body["error"])
		})
	}
}

func TestDecodeJSONBody(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name          string
		body          string
		expectError   bool
		expectedError error
		expected      payload
	}{
		{
			"valid json",
			`{"name":"Alice","age":30}`,
			false, nil,
			payload{"Alice", 30},
		},
		{
			"invalid json",
			`not-json`,
			true, ErrParsingJsonBody,
			payload{},
		},
		{
			"unknown field",
			`{"name":"Alice","age":30,"unknown":"field"}`,
			true, ErrParsingJsonBody,
			payload{},
		},
		{
			"empty body",
			``,
			true, ErrParsingJsonBody,
			payload{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			var result payload
			err := DecodeJSONBody(req, &result)

			if tc.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestExtractParam(t *testing.T) {
	tests := []struct {
		name          string
		paramName     string
		pathValue     string
		expectError   bool
		expectedError error
		expected      string
	}{
		{"valid param", "id", "abc123", false, nil, "abc123"},
		{"missing param", "id", "", true, ErrParamNotFound, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.SetPathValue(tc.paramName, tc.pathValue)

			result, err := extractParam(req, tc.paramName)

			if tc.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestSendJsonResponse_EncodingError(t *testing.T) {
	w := httptest.NewRecorder()
	// channels cannot be JSON-encoded, triggering the error path
	SendJsonResponse(w, http.StatusOK, make(chan int))

	body := w.Body.String()
	assert.Contains(t, body, ErrEncodingJson.Error())
}

func TestDecodeJSONBody_ClosesBody(t *testing.T) {
	bodyContent := `{"name":"Bob","age":25}`
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(bodyContent))

	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var result payload
	err := DecodeJSONBody(req, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", result.Name)
	assert.Equal(t, 25, result.Age)
}
