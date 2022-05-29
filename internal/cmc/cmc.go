// Package coinmarketcap Coin Market Cap API client for Go
package cmc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Client the CoinMarketCap client
type Client struct {
	proAPIKey string
	proxyUrl  string
}

// Status is the status structure
type Status struct {
	Timestamp    string  `json:"timestamp"`
	ErrorCode    int     `json:"error_code"`
	ErrorMessage *string `json:"error_message"`
	Elapsed      int     `json:"elapsed"`
	CreditCount  int     `json:"credit_count"`
}

// Response is the response structure
type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

type MapListing struct {
	ID                  float64  `json:"id"`
	Name                string   `json:"name"`
	Symbol              string   `json:"symbol"`
	Slug                string   `json:"slug"`
	Rank                int      `json:"rank"`
	IsActive            int      `json:"is_active"`
	FirstHistoricalData string   `json:"first_historical_data"`
	LastHistoricalData  string   `json:"last_historical_data"`
	Platform            Platform `json:"platform"`
}

// CryptocurrencyInfo options
type CryptocurrencyInfo struct {
	ID       float64                `json:"id"`
	Name     string                 `json:"name"`
	Symbol   string                 `json:"symbol"`
	Category string                 `json:"category"`
	Slug     string                 `json:"slug"`
	Logo     string                 `json:"logo"`
	Tags     []string               `json:"tags"`
	Urls     map[string]interface{} `json:"urls"`
	Platform Platform               `json:"platform"`
}

// Listing is the listing structure
type Listing struct {
	ID                float64           `json:"id"`
	Name              string            `json:"name"`
	Symbol            string            `json:"symbol"`
	Slug              string            `json:"slug"`
	CirculatingSupply float64           `json:"circulating_supply"`
	TotalSupply       float64           `json:"total_supply"`
	MaxSupply         float64           `json:"max_supply"`
	DateAdded         string            `json:"date_added"`
	NumMarketPairs    float64           `json:"num_market_pairs"`
	CMCRank           float64           `json:"cmc_rank"`
	LastUpdated       string            `json:"last_updated"`
	Platform          Platform          `json:"platform"`
	Quote             map[string]*Quote `json:"quote"`
}

// ConvertListing is the converted listing structure
type ConvertListing struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`
	Symbol      string                   `json:"symbol"`
	Amount      float64                  `json:"amount"`
	LastUpdated string                   `json:"last_updated"`
	Quote       map[string]*ConvertQuote `json:"quote"`
}

// MapListing is the structure of a map listing
type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

// ConvertQuote is the converted listing structure
type ConvertQuote struct {
	Price       float64 `json:"price"`
	LastUpdated string  `json:"last_updated"`
}

// Quote is the quote structure
type Quote struct {
	Price            float64 `json:"price"`
	Volume24H        float64 `json:"volume_24h"`
	PercentChange1H  float64 `json:"percent_change_1h"`
	PercentChange24H float64 `json:"percent_change_24h"`
	PercentChange7D  float64 `json:"percent_change_7d"`
	MarketCap        float64 `json:"market_cap"`
	LastUpdated      string  `json:"last_updated"`
}

// MapOptions options
type MapOptions struct {
	Start  int
	Limit  int
	Symbol string
}

// InfoOptions options
type InfoOptions struct {
	ID     string
	Symbol string
	Slug   string
}

// ListingOptions options
type ListingOptions struct {
	Start   int
	Limit   int
	Convert string
	Sort    string
}

// QuoteOptions options
type QuoteOptions struct {
	// Covert suppots multiple currencies command separated. eg. "BRL,USD"
	Convert string
	// Symbols suppots multiple tickers command separated. eg. "BTC,ETH,XRP"
	Symbol string
}

// ConvertOptions options
type ConvertOptions struct {
	Amount  float64
	Symbol  string
	Convert string
}

var (
	// ErrTypeAssertion is type assertion error
	ErrTypeAssertion = errors.New("type assertion error")
)

var (
	baseURL   = "https://pro-api.coinmarketcap.com/v1"
	basev2URL = "https://pro-api.coinmarketcap.com/v2"
)

// NewClient initializes a new client
func NewClient(key string, proxy string) *Client {
	if key == "" {
		key = os.Getenv("CMC_PRO_API_KEY")
	}

	c := &Client{
		proAPIKey: key,
		proxyUrl:  proxy,
	}

	return c
}

// Info returns all static metadata for one or more cryptocurrencies including name, symbol, logo, and its various registered URLs.
func (s *Client) Info(options *InfoOptions) (map[string]*CryptocurrencyInfo, error) {
	var params []string
	if options == nil {
		options = new(InfoOptions)
	}
	if options.ID != "" {
		params = append(params, fmt.Sprintf("id=%s", options.ID))
	}
	if options.Symbol != "" {
		params = append(params, fmt.Sprintf("symbol=%s", options.Symbol))
	}
	if options.Slug != "" {
		params = append(params, fmt.Sprintf("slug=%s", strings.ToLower(options.Slug)))
	}

	url := fmt.Sprintf("%s/cryptocurrency/info?%s", basev2URL, strings.Join(params, "&"))
	body, err := s.makeReq(url)
	if err != nil {
		return nil, err
	}

	resp := new(Response)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	var result = make(map[string]*CryptocurrencyInfo)
	ifcs, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, ErrTypeAssertion
	}

	for k, v := range ifcs {
		info := new(CryptocurrencyInfo)
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, info)
		if err != nil {
			return nil, err
		}
		result[k] = info
	}

	return result, nil
}

// LatestListings gets a paginated list of all cryptocurrencies with latest market data. You can configure this call to sort by market cap or another market ranking field. Use the "convert" option to return market values in multiple fiat and cryptocurrency conversions in the same call.
func (s *Client) LatestListings(options *ListingOptions) ([]*Listing, error) {
	var params []string
	if options == nil {
		options = new(ListingOptions)
	}
	if options.Start != 0 {
		params = append(params, fmt.Sprintf("start=%v", options.Start))
	}
	if options.Limit != 0 {
		params = append(params, fmt.Sprintf("limit=%v", options.Limit))
	}
	if options.Convert != "" {
		params = append(params, fmt.Sprintf("convert=%s", options.Convert))
	}
	if options.Sort != "" {
		params = append(params, fmt.Sprintf("sort=%s", options.Sort))
	}

	url := fmt.Sprintf("%s/cryptocurrency/listings/latest?%s", baseURL, strings.Join(params, "&"))
	body, err := s.makeReq(url)
	if err != nil {
		return nil, err
	}

	resp := new(Response)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("JSON Error: [%s]. Response body: [%s]", err.Error(), string(body))
	}

	var listings []*Listing
	ifcs, ok := resp.Data.([]interface{})
	if !ok {
		return nil, ErrTypeAssertion
	}

	for i := range ifcs {
		ifc := ifcs[i]
		listing := new(Listing)
		b, err := json.Marshal(ifc)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, listing)
		if err != nil {
			return nil, err
		}

		listings = append(listings, listing)
	}

	return listings, nil
}

// Map returns a paginated list of all cryptocurrencies by CoinMarketCap ID.
func (s *Client) Map(options *MapOptions) ([]*MapListing, error) {
	var params []string
	if options == nil {
		options = new(MapOptions)
	}

	if options.Start != 0 {
		params = append(params, fmt.Sprintf("start=%d", options.Start))
	}
	if options.Limit != 0 {
		params = append(params, fmt.Sprintf("limit=%d", options.Limit))
	}
	if options.Symbol != "" {
		params = append(params, fmt.Sprintf("symbol=%s", options.Symbol))
	}

	url := fmt.Sprintf("%s/cryptocurrency/map?%s", baseURL, strings.Join(params, "&"))

	body, err := s.makeReq(url)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("JSON Error: [%s]. Response body: [%s]", err.Error(), string(body))
	}

	var result []*MapListing
	ifcs, ok := resp.Data.(interface{})
	if !ok {
		return nil, ErrTypeAssertion
	}

	for _, item := range ifcs.([]interface{}) {
		value := new(MapListing)
		b, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, value)
		if err != nil {
			return nil, err
		}

		result = append(result, value)
	}

	return result, nil
}

// LatestQuotes gets latest quote for each specified symbol. Use the "convert" option to return market values in multiple fiat and cryptocurrency conversions in the same call.
func (s *Client) LatestQuotes(options *QuoteOptions) ([]*Listing, error) {
	var params []string
	if options == nil {
		options = new(QuoteOptions)
	}

	if options.Symbol != "" {
		params = append(params, fmt.Sprintf("symbol=%s", options.Symbol))
	}

	if options.Convert != "" {
		params = append(params, fmt.Sprintf("convert=%s", options.Convert))
	}

	url := fmt.Sprintf("%s/cryptocurrency/quotes/latest?%s", basev2URL, strings.Join(params, "&"))

	body, err := s.makeReq(url)
	resp := new(Response)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("JSON Error: [%s]. Response body: [%s]", err.Error(), string(body))
	}

	var quotesLatest []*Listing
	ifcs, ok := resp.Data.(interface{})
	if !ok {
		return nil, ErrTypeAssertion
	}

	for _, coinObj := range ifcs.(map[string]interface{}) {
		for _, obj := range coinObj.([]interface{}) {
			quoteLatest := new(Listing)
			b, err := json.Marshal(obj)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(b, quoteLatest)
			if err != nil {
				return nil, err
			}

			quotesLatest = append(quotesLatest, quoteLatest)
		}
	}

	return quotesLatest, nil
}

// PriceConversion Convert an amount of one currency into multiple cryptocurrencies or fiat currencies at the same time using the latest market averages. Optionally pass a historical timestamp to convert values based on historic averages.
func (s *Client) PriceConversion(options *ConvertOptions) (*ConvertListing, error) {
	var params []string
	if options == nil {
		options = new(ConvertOptions)
	}

	if options.Amount != 0 {
		params = append(params, fmt.Sprintf("amount=%f", options.Amount))
	}

	if options.Symbol != "" {
		params = append(params, fmt.Sprintf("symbol=%s", options.Symbol))
	}

	if options.Convert != "" {
		params = append(params, fmt.Sprintf("convert=%s", options.Convert))
	}

	url := fmt.Sprintf("%s/tools/price-conversion?%s", basev2URL, strings.Join(params, "&"))

	body, err := s.makeReq(url)

	resp := new(Response)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("JSON Error: [%s]. Response body: [%s]", err.Error(), string(body))
	}

	ifc, ok := resp.Data.(interface{})
	if !ok {
		return nil, ErrTypeAssertion
	}

	listing := new(ConvertListing)
	b, err := json.Marshal(ifc)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, listing)
	if err != nil {
		return nil, err
	}

	return listing, nil
}
