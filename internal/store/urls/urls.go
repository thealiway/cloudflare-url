package urls

import (
	cloudsql "cloudflareurl/internal/store"
	"database/sql"
)

func NewURLStore() (*sql.DB, error) {
	db, err := cloudsql.ConnectWithConnector("urls")
	if err != nil {
		return nil, err
	}

	return db, nil
}
