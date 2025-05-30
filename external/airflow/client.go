package airflow

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	pageLimit = 99999
)

type Request struct {
	Path   string
	Method string
	Body   []byte
	Query  string
}

type Auth struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

type Client struct {
	client *http.Client
}

func NewAirflowClient() *Client {
	return &Client{client: &http.Client{}}
}

func (ac Client) Invoke(ctx context.Context, r Request, auth Auth) ([]byte, error) {
	var resp []byte

	endpoint := buildEndPoint(auth.Host, r.Path, r.Query)
	request, err := http.NewRequestWithContext(ctx, r.Method, endpoint, bytes.NewBuffer(r.Body))
	if err != nil {
		return resp, fmt.Errorf("failed to build http request for %s due to %w", endpoint, err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth.Token))))

	httpResp, respErr := ac.client.Do(request)
	if respErr != nil {
		return resp, fmt.Errorf("failed to call airflow %s due to %w", endpoint, respErr)
	}
	if httpResp.StatusCode != http.StatusOK {
		httpResp.Body.Close()
		return resp, fmt.Errorf("status code received %d on calling %s", httpResp.StatusCode, endpoint)
	}
	return parseResponse(httpResp)
}

func parseResponse(resp *http.Response) ([]byte, error) {
	var body []byte
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return body, fmt.Errorf("failed to read airflow response: %w", err)
	}
	return body, nil
}

func buildEndPoint(host, path, query string) string {
	host = strings.Trim(host, "/")
	u := &url.URL{
		Scheme:   "http",
		Host:     host,
		Path:     path,
		RawQuery: query,
	}
	return u.String()
}
