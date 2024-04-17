package controllers

import (
	apimodels "cloudflareurl/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type URLController struct {
	URLStore *sql.DB
}

type URL struct {
	LongURL        string `json:"longURL"`
	ShortenedURL   string `json:"shortenedURL"`
	ExpirationDate int64  `json:"expirationDate"`
}

type Usage struct {
	ID           string `json:"id"`
	ShortenedURL string `json:"shortenedURL"`
	UsageTime    int64  `json:"usageTime"`
}

type Stats struct {
	Day     int `json:"day"`
	Week    int `json:"week"`
	AllTime int `json:"allTime"`
}

func NewURLController(uStore *sql.DB) (*URLController, error) {
	return &URLController{
		URLStore: uStore,
	}, nil
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

	err = u.LogUsage(shortenedURL)
	if err != nil {
		return "", nil
	}

	return url.LongURL, nil
}

func (u *URLController) LogUsage(shortenedURL string) error {
	stmt, err := u.URLStore.Prepare("INSERT INTO usage (id, shortened_url, usage_time) VALUES ($1, $2, $3)")
	if err != nil {
		fmt.Println("Unable to prepare usage")
		return err
	}
	defer stmt.Close()

	id := uuid.New().String()
	_, err = stmt.Exec(id, shortenedURL, time.Now().Unix())
	if err != nil {
		fmt.Printf("Unable to insert usage into DB: %+v \n", err)
		return err
	}

	return nil
}

func (u *URLController) GetUsage(shortenedURL string) (*Stats, error) {
	rows, err := u.URLStore.Query("SELECT * FROM usage WHERE shortened_url = $1", shortenedURL)
	if err != nil {
		fmt.Printf("unable to get all usages: %+v \n", err)
		return nil, err
	}

	var usages []Usage
	for rows.Next() {
		var usage Usage
		err = rows.Scan(&usage.ID, &usage.ShortenedURL, &usage.UsageTime)
		if err != nil {
			fmt.Printf("unable to unmarshal usage rows: %+v \n", err)
			return nil, err
		}
		usages = append(usages, usage)
	}

	if len(usages) == 0 {
		emptyErr := errors.New("URL does not exist")
		return nil, emptyErr
	}

	usageStats := &Stats{
		Day:     0,
		Week:    0,
		AllTime: len(usages),
	}

	for _, use := range usages {
		dayAgo := time.Now().Add(-time.Hour * 24)
		weekAgo := time.Now().Add(-time.Hour * 168)
		uTime := time.Unix(use.UsageTime, 0)
		if dayAgo.Before(uTime) {
			usageStats.Day = usageStats.Day + 1
			usageStats.Week = usageStats.Week + 1
			continue
		}
		if weekAgo.Before(uTime) {
			usageStats.Week = usageStats.Week + 1
			continue
		}
	}

	return usageStats, nil
}

func (u *URLController) DeleteURL(shortenedURL string) error {
	_, err := u.URLStore.Exec("DELETE FROM urls WHERE shortened_url = $1", shortenedURL)
	if err != nil {
		fmt.Printf("Unable to delete: %+v \n", err)
		return err
	}

	return nil
}
