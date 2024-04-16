package controllers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type URLController struct {
	URLStore *sql.DB
}

type URL struct {
	ID           string
	LongURL      string
	ShortenedURL string
}

func NewURLController(uStore *sql.DB) (*URLController, error) {
	return &URLController{
		URLStore: uStore,
	}, nil
}

func (u *URLController) CreateShortenedURL(url string) (*URL, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortened := make([]byte, keyLength)
	for i := range shortened {
		shortened[i] = charset[rand.Intn(len(charset))]
	}

	id := uuid.New().String()

	stmt, err := u.URLStore.Prepare("INSERT INTO urls (id, long_url, shortened_url) VALUES ($1, $2, $3)")
	if err != nil {
		fmt.Println("Unable to prepare")
		return nil, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id, url, shortened); err != nil {
		fmt.Println("Unable to insert into DB")
		return nil, err
	}

	return &URL{
		ID:           id,
		LongURL:      url,
		ShortenedURL: string(shortened),
	}, nil
}
