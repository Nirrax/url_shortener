package main

import (
	"time"

	"github.com/theckman/go-securerandom"
)

type Url struct {
	ID        int32
	LongUrl   string
	ShortUrl  string
	CreatedAt time.Time
}

func NewUrl(longUrl string) (*Url, error) {

	id, err := securerandom.Int32()

	if err != nil {
		return nil, err
	}

	shortUrl, err := securerandom.URLBase64InBytes(8)

	if err != nil {
		return nil, err
	}

	return &Url{
		ID:        id,
		LongUrl:   longUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now(),
	}, nil
}
