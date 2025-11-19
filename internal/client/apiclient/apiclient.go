// Package apiclient предоставляет обёртку над http.Client
// для выполнения HTTP-запросов из клиентского CLI приложения.
// Он стандартизирует таймауты, обработку ответа и формат возвращаемой структуры.
package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
)

// APIClient инкапсулирует http.Client и обеспечивает единый
// способ отправки запросов и получения ответов в формате models.Response.
//
// Клиент используется всеми сервисами CLI для общения с сервером.
type ApiClient struct {
	httpClient *http.Client
}

// NewAPIClient создаёт новый HTTP-клиент с заданным таймаутом.
// timeout — длительность тайм-аута, например: 10 * time.Second.
func NewApiClient(waitTime time.Duration) *ApiClient {
	return &ApiClient{
		httpClient: &http.Client{Timeout: waitTime * time.Second},
	}
}

// Do выполняет HTTP-запрос и возвращает структуру Response,
// содержащую статус, заголовки, тело ответа и cookies.
//
// Метод оборачивает http.Client.Do, полностью читает тело ответа
// и закрывает его.
//
// В случае сетевой ошибки или ошибки чтения тела возвращает error.
func (client *ApiClient) Do(req *http.Request) (*models.Response, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	return &models.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
		Cookies:    resp.Cookies(),
	}, nil
}
