package opensea

import (
	"os"
)

// Client the CoinMarketCap client
type Client struct {
	proAPIKey string
	proxyUrl  string
}

var (
	baseURL = "https://api.opensea.io/api/v1"
)

// NewClient initializes a new client
func NewClient(key string, proxy string) *Client {
	if key == "" {
		key = os.Getenv("OPENSEA_API_KEY")
	}

	c := &Client{
		proAPIKey: key,
		proxyUrl:  proxy,
	}

	return c
}
