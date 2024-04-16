package controllers

import (
	"database/sql"
	"math/rand"
	"time"
)

type URLController struct {
	URLStore *sql.DB
}

type URL struct {
	ID           string
	LongURL      string
	ShortenedURL string
}

func NewURLController(sStore *sql.DB) *URLController {
	return &URLController{}
}

func (u *URLController) CreateShortenedURL(url string) *URL {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortened := make([]byte, keyLength)
	for i := range shortened {
		shortened[i] = charset[rand.Intn(len(charset))]
	}

	u.URLStore.Exec("INSERT INTO urls ")

	return &URL{
		LongURL:      url,
		ShortenedURL: string(shortened),
	}
}
