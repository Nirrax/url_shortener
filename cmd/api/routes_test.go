package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nirrax/url_shortener/internal/entity"
	"github.com/nirrax/url_shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestApp(svc service.Service) *App {
	return &App{Service: svc}
}

func TestShowUrlHandler_Success(t *testing.T) {
	svc, urlSvc := newMockService()
	expected := &entity.UrlDTO{ShortUrl: "abc123", LongUrl: "https://example.com"}
	urlSvc.On("GetUrl", "abc123").Return(expected, nil)

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/url/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var got entity.UrlDTO
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, expected.ShortUrl, got.ShortUrl)
	assert.Equal(t, expected.LongUrl, got.LongUrl)
	urlSvc.AssertExpectations(t)
}

func TestShowUrlHandler_NotFound(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("GetUrl", "missing").Return(nil, service.ErrUrlNotFound)

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/url/missing", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	urlSvc.AssertExpectations(t)
}

func TestShowUrlHandler_InternalError(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("GetUrl", "abc123").Return(nil, errors.New("db connection lost"))

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/url/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	urlSvc.AssertExpectations(t)
}

func TestRedirectUrlHandler_Success(t *testing.T) {
	svc, urlSvc := newMockService()
	dto := &entity.UrlDTO{ShortUrl: "abc123", LongUrl: "https://example.com"}
	urlSvc.On("GetUrl", "abc123").Return(dto, nil)

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	assert.Equal(t, dto.LongUrl, rec.Header().Get("Location"))
	urlSvc.AssertExpectations(t)
}

func TestRedirectUrlHandler_NotFound(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("GetUrl", "missing").Return(nil, service.ErrUrlNotFound)

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	urlSvc.AssertExpectations(t)
}

func TestRedirectUrlHandler_InternalError(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("GetUrl", "abc123").Return(nil, errors.New("timeout"))

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	urlSvc.AssertExpectations(t)
}

func TestCreateUrlHandler_Success(t *testing.T) {
	svc, urlSvc := newMockService()
	created := &entity.UrlDTO{ShortUrl: "newshort", LongUrl: "https://new.example.com"}
	urlSvc.On("CreateUrl", "https://new.example.com").Return(created, nil)

	router := newTestApp(svc).NewRouter()

	body, _ := json.Marshal(entity.UrlRequest{Url: "https://new.example.com"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var got entity.UrlDTO
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, created.ShortUrl, got.ShortUrl)
	urlSvc.AssertExpectations(t)
}

func TestCreateUrlHandler_InvalidBody(t *testing.T) {
	svc, urlSvc := newMockService()

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	urlSvc.AssertNotCalled(t, "CreateUrl", mock.Anything)
}

func TestCreateUrlHandler_ServiceError(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("CreateUrl", "https://bad.example.com").Return(nil, errors.New("create failed"))

	router := newTestApp(svc).NewRouter()

	body, _ := json.Marshal(entity.UrlRequest{Url: "https://bad.example.com"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	urlSvc.AssertExpectations(t)
}

func TestDeleteUrlHandler_Success(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("DeleteUrl", "abc123").Return(nil)

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodDelete, "/abc123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	urlSvc.AssertExpectations(t)
}

func TestDeleteUrlHandler_NotFound(t *testing.T) {
	svc, urlSvc := newMockService()
	urlSvc.On("DeleteUrl", "missing").Return(service.ErrUrlNotFound)

	router := newTestApp(svc).NewRouter()

	req := httptest.NewRequest(http.MethodDelete, "/missing", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	urlSvc.AssertExpectations(t)
}
