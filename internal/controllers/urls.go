package controllers

import (
	apimodels "cloudflareurl/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type URLController struct {
	URLStore *sql.DB
}

type URL struct {
	LongURL        string `json:"longURL"`
	ShortenedURL   string `json:"shortenedURL"`
	ExpirationDate int64  `json:"expirationDate"`
}

func NewURLController(uStore *sql.DB) *URLController {
	return &URLController{
		URLStore: uStore,
	}
}

func (u *URLController) CreateShortenedURL(input *apimodels.Input) (*URL, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortened := make([]byte, keyLength)
	for i := range shortened {
		shortened[i] = charset[rand.Intn(len(charset))]
	}

	stmt, err := u.URLStore.Prepare("INSERT INTO urls (long_url, shortened_url, expiration_date) VALUES ($1, $2, $3)")
	if err != nil {
		fmt.Println("Unable to prepare")
		return nil, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.URL, shortened, input.ExpirationDate); err != nil {
		fmt.Println("Unable to insert into DB")
		return nil, err
	}

	s := fmt.Sprintf("https://cloudflare-url-ownx73g3lq-uw.a.run.app/s/%s", string(shortened))

	return &URL{
		LongURL:        input.URL,
		ShortenedURL:   s,
		ExpirationDate: input.ExpirationDate,
	}, nil
}

func (u *URLController) GetOriginalURL(shortenedURL string) (string, error) {
	row := u.URLStore.QueryRow("SELECT * FROM urls WHERE shortened_url = $1", shortenedURL)

	var url URL
	err := row.Scan(&url.LongURL, &url.ShortenedURL, &url.ExpirationDate)
	if err != nil {
		return "", err
	}

	uTime := time.Unix(url.ExpirationDate, 0)
	if url.ExpirationDate != 0 && time.Now().After(uTime) {
		fmt.Println("link is expired")
		expiredErr := errors.New("link is expired")
		return "", expiredErr
	}

	return url.LongURL, nil
}

func (u *URLController) DeleteURL(shortenedURL string) error {
	_, err := u.URLStore.Exec("DELETE FROM urls WHERE shortened_url = $1", shortenedURL)
	if err != nil {
		fmt.Printf("Unable to delete: %+v \n", err)
		return err
	}

	return nil
}
