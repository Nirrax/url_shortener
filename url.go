package main

import (
	"time"

	"github.com/theckman/go-securerandom"
)

type Url struct {
	ID        int32     `json:"id"`
	LongUrl   string    `json:"longUrl"`
	ShortUrl  string    `json:"shortUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type UrlDto struct {
	LongUrl string `json:"longUrl"`
}

func NewUrl(longUrl string) (*Url, error) {
	shortUrl, err := securerandom.URLBase64InBytes(8)

	if err != nil {
		return nil, err
	}

	return &Url{
		LongUrl:   longUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now().UTC(),
	}, nil
}
