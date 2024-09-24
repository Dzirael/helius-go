package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const defaultUrl = "https://api.helius.xyz/"
const apiVersion = "v0"

type HttpClient interface {
	GetParsedTransaction(ctx context.Context, transactions ...string) (result []Transaction, err error)
	GetTransactionHistory(ctx context.Context, params *TransactionQuerry, address string) (result []Transaction, err error)
	GetAllTransactionHistory(ctx context.Context, params *TransactionQuerry, address string) TransactionHistory
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseUrl    *url.URL
}

func New(apiKey string) HttpClient {
	baseUrl, _ := url.Parse(defaultUrl + apiVersion)
	query := baseUrl.Query()
	query.Set("api-key", apiKey)
	baseUrl.RawQuery = query.Encode()

	return &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		baseUrl:    baseUrl,
	}
}

func (c *Client) GetParsedTransaction(ctx context.Context, transactions ...string) (result []Transaction, err error) {
	urlPath, err := c.getUrl("transactions/", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get url")
	}
	body := map[string]interface{}{
		"transactions": transactions,
	}

	req, err := c.newRequest(ctx, "POST", urlPath, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	err = c.call(req, &result)
	return
}

func (c *Client) GetTransactionHistory(ctx context.Context, params *TransactionQuerry, address string) (result []Transaction, err error) {
	urlPath, err := c.getUrl("addresses/"+address+"/transactions/", params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get url")
	}

	req, err := c.newRequest(ctx, "GET", urlPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	err = c.call(req, &result)
	return
}

func (c *Client) GetAllTransactionHistory(ctx context.Context, params *TransactionQuerry, address string) TransactionHistory {
	var result = TransactionHistory{
		channel: make(chan Transaction),
	}

	go func() {
		defer close(result.channel)
		if params == nil {
			params = &TransactionQuerry{}
		}
		for {
			transactions, err := c.GetTransactionHistory(ctx, params, address)
			if err != nil {
				result.err = err
			}

			if len(transactions) == 0 {
				break
			}

			for _, transaction := range transactions {
				result.channel <- transaction
			}
			params.Before = transactions[len(transactions)-1].Signature
		}
	}()
	return result
}

func (c *Client) getUrl(endpoint string, params *TransactionQuerry) (string, error) {
	u := *c.baseUrl
	newPath, err := url.JoinPath(u.Path, endpoint)
	if err != nil {
		return "", err
	}
	u.Path = newPath

	fmt.Println(u.String())

	if params != nil {
		query := u.Query()
		for key, value := range params.ToMap() {
			query.Set(key, value)
		}
		u.RawQuery = query.Encode()
	}
	return u.String(), nil
}

func (c *Client) newRequest(ctx context.Context, method, url string, req interface{}) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return request, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil

}

func (c *Client) call(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("bad status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	return nil
}
