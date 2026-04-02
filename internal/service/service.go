package service

import (
	"github.com/nirrax/url_shortener/internal/repository"
	"github.com/nirrax/url_shortener/internal/utils"
)

type Service interface {
	UrlService() UrlService
}

type ServiceImpl struct {
	urls UrlService
}

func (s *ServiceImpl) UrlService() UrlService {
	return s.urls
}

func NewService(repository repository.Repository, encoder utils.Encoder) Service {
	return &ServiceImpl{
		urls: newUrlService(repository, encoder),
	}
}
