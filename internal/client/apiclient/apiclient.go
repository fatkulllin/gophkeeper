package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
)

type ApiClient struct {
	httpClient *http.Client
}

func NewApiClient(waitTime int64) *ApiClient {
	timeoutClient := time.Duration(waitTime)
	return &ApiClient{
		httpClient: &http.Client{Timeout: timeoutClient * time.Second},
	}
}

func (client *ApiClient) Do(req *http.Request) (*models.Response, error) {
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
