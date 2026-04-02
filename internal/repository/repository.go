package repository

import "github.com/nirrax/url_shortener/internal/database"

type Repository interface {
	Urls() UrlRepository
}

type RepositoryImpl struct {
	urls UrlRepository
}

func (r *RepositoryImpl) Urls() UrlRepository {
	return r.urls
}

func NewRepository(db database.DBclient) Repository {
	return &RepositoryImpl{
		urls: newUrlRepository(db),
	}
}
