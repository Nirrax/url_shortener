package api

import (
	"github.com/nirrax/url_shortener/internal/entity"
	"github.com/nirrax/url_shortener/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockUrlService struct {
	mock.Mock
}

func (m *MockUrlService) GetUrl(url string) (*entity.UrlDTO, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UrlDTO), args.Error(1)
}

func (m *MockUrlService) CreateUrl(url string) (*entity.UrlDTO, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UrlDTO), args.Error(1)
}

func (m *MockUrlService) DeleteUrl(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

type MockService struct {
	mock.Mock
	urlService *MockUrlService
}

func newMockService() (*MockService, *MockUrlService) {
	urlSvc := new(MockUrlService)
	svc := &MockService{urlService: urlSvc}
	return svc, urlSvc
}

func (m *MockService) UrlService() service.UrlService {
	return m.urlService
}
