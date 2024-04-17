package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UsageController struct {
	UsageStore *sql.DB
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

func NewUsageController(store *sql.DB) *UsageController {
	return &UsageController{
		UsageStore: store,
	}
}

func (a *UsageController) GetUsage(shortenedURL string) (*Stats, error) {
	rows, err := a.UsageStore.Query("SELECT * FROM usage WHERE shortened_url = $1", shortenedURL)
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

func (a *UsageController) LogUsage(shortenedURL string) error {
	stmt, err := a.UsageStore.Prepare("INSERT INTO usage (id, shortened_url, usage_time) VALUES ($1, $2, $3)")
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
