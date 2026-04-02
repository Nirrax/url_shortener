package service

import (
	"errors"
	"fmt"

	"github.com/nirrax/url_shortener/internal/entity"
	"github.com/nirrax/url_shortener/internal/repository"
	"github.com/nirrax/url_shortener/internal/utils"
)

var (
	ErrInvalidUrl  = fmt.Errorf("provided argument is not an url")
	ErrUrlNotFound = fmt.Errorf("url does not exist")
)

type UrlService interface {
	GetUrl(shortUrl string) (*entity.UrlDTO, error)
	CreateUrl(longUrl string) (*entity.UrlDTO, error)
	DeleteUrl(shortUrl string) error
}

type UrlServiceImpl struct {
	repository.Repository
	encoder utils.Encoder
}

func newUrlService(repository repository.Repository, encoder utils.Encoder) *UrlServiceImpl {
	return &UrlServiceImpl{
		repository,
		encoder,
	}
}

func (s *UrlServiceImpl) GetUrl(shortUrl string) (*entity.UrlDTO, error) {
	id, err := s.encoder.Decode(shortUrl)
	if err != nil {
		return nil, ErrInvalidUrl
	}

	url, err := s.Urls().GetByID(int(id))
	if err != nil {
		return nil, ErrUrlNotFound
	}

	return &entity.UrlDTO{
		ID:        url.ID,
		ShortUrl:  shortUrl,
		LongUrl:   url.LongUrl,
		CreatedAt: url.CreatedAt,
	}, nil
}

func (s *UrlServiceImpl) CreateUrl(longUrl string) (*entity.UrlDTO, error) {
	url, err := s.Urls().GetByLongUrl(longUrl)
	if err != nil && !errors.Is(err, ErrUrlNotFound) {
		return nil, err
	}

	if errors.Is(err, ErrUrlNotFound) {
		url, err = s.Urls().Create(longUrl)
		if err != nil {
			if utils.IsUniqueConstraintViolation(err) {
				url, err = s.Urls().GetByLongUrl(longUrl)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	return &entity.UrlDTO{
		ID:        url.ID,
		ShortUrl:  s.encoder.Encode(uint64(url.ID)),
		LongUrl:   url.LongUrl,
		CreatedAt: url.CreatedAt,
	}, nil
}

func (s *UrlServiceImpl) DeleteUrl(shortUrl string) error {
	id, err := s.encoder.Decode(shortUrl)
	if err != nil {
		return ErrInvalidUrl
	}

	err = s.Urls().DeleteByID(int(id))
	if err != nil {
		return ErrUrlNotFound
	}

	return nil
}
