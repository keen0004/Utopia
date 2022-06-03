package cmc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// doReq HTTP client
func doReq(req *http.Request, proxy string) ([]byte, error) {
	requestTimeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   requestTimeout,
		Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
	}

	if proxy != "" {
		u, _ := url.Parse(proxy)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(u)}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// makeReq HTTP request helper
func (s *Client) makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", s.proAPIKey)

	resp, err := doReq(req, s.proxyUrl)
	if err != nil {
		return nil, err
	}

	return resp, err
}
