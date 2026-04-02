package entity

import "time"

type Url struct {
	ID        int       `json:"id" db:"id"`
	LongUrl   string    `json:"longUrl" db:"long_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type UrlRequest struct {
	Url string `json:"url"`
}

type UrlDTO struct {
	ID        int       `json:"id" db:"id"`
	ShortUrl  string    `json:"shortUrl" db:"short_url"`
	LongUrl   string    `json:"longUrl" db:"long_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
