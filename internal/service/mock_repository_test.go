package service

import (
	"github.com/nirrax/url_shortener/internal/entity"
	"github.com/nirrax/url_shortener/internal/repository"
	"github.com/stretchr/testify/mock"
)

type MockUrlRepository struct {
	mock.Mock
}

func (m *MockUrlRepository) GetByID(id int) (*entity.Url, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Url), args.Error(1)
}

func (m *MockUrlRepository) GetByLongUrl(longUrl string) (*entity.Url, error) {
	args := m.Called(longUrl)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Url), args.Error(1)
}

func (m *MockUrlRepository) Create(longUrl string) (*entity.Url, error) {
	args := m.Called(longUrl)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Url), args.Error(1)
}

func (m *MockUrlRepository) DeleteByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockRepository struct {
	mock.Mock
	urls *MockUrlRepository
}

func (m *MockRepository) Urls() repository.UrlRepository {
	return m.urls
}
