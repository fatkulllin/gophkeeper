package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(waitTime int64) *Client {
	timeoutClient := time.Duration(waitTime)
	return &Client{
		httpClient: &http.Client{Timeout: timeoutClient * time.Second},
	}
}

func (client *Client) Do(req *http.Request) (*models.Response, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return &models.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
		Cookies:    resp.Cookies(),
	}, nil
}
