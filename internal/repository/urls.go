package repository

import (
	"time"

	"github.com/nirrax/url_shortener/internal/database"
	"github.com/nirrax/url_shortener/internal/entity"
)

type UrlRepository interface {
	GetByID(id int) (*entity.Url, error)
	GetByLongUrl(longUrl string) (*entity.Url, error)
	Create(longUrl string) (*entity.Url, error)
	DeleteByID(id int) error
}

type UrlRepositoryImpl struct {
	db database.DBclient
}

func newUrlRepository(db database.DBclient) *UrlRepositoryImpl {
	return &UrlRepositoryImpl{
		db: db,
	}
}

func (r *UrlRepositoryImpl) GetByID(id int) (*entity.Url, error) {
	url := entity.Url{}

	query := `SELECT * FROM urls WHERE id = $1`
	err := r.db.FetchOne(&url, query, id)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *UrlRepositoryImpl) GetByLongUrl(longUrl string) (*entity.Url, error) {
	url := entity.Url{}

	query := `SELECT * FROM urls WHERE long_url = $1`
	err := r.db.FetchOne(&url, query, longUrl)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *UrlRepositoryImpl) Create(longUrl string) (*entity.Url, error) {
	url := entity.Url{}

	query := `INSERT INTO urls (long_url, created_at) VALUES ($1, $2) RETURNING *`
	err := r.db.FetchOne(&url, query, longUrl, time.Now())
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *UrlRepositoryImpl) DeleteByID(id int) error {
	query := `DELETE FROM urls WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
