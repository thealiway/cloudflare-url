package controllers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type URLController struct {
	URLStore *sql.DB
}

type URL struct {
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

	stmt, err := u.URLStore.Prepare("INSERT INTO urls (long_url, shortened_url) VALUES ($1, $2)")
	if err != nil {
		fmt.Println("Unable to prepare")
		return nil, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(url, shortened); err != nil {
		fmt.Println("Unable to insert into DB")
		return nil, err
	}

	return &URL{
		LongURL:      url,
		ShortenedURL: string(shortened),
	}, nil
}

func (u *URLController) GetOriginalURL(shortenedURL string) (string, error) {
	row := u.URLStore.QueryRow("SELECT long_url FROM urls WHERE shortened_url = $1", shortenedURL)

	var longURL string
	err := row.Scan(&longURL)
	if err != nil {
		return "", err
	}

	fmt.Printf("longURL: %+v \n", longURL)

	return longURL, nil
}
