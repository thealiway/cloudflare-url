package apimodels

type Input struct {
	URL            string `json:"url"`
	ExpirationDate int64  `json:"expirationDate,omitempty"`
}
