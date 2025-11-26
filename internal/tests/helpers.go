package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func createShortUrl(t *testing.T, baseUrl, longUrl string) *http.Response {
	t.Helper()

	urlString := fmt.Sprintf("%s/create", baseUrl)
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		t.Fatalf("failed to parse POST /create url: %v", err)
	}

	b, err := json.Marshal(map[string]string{
		"url": longUrl,
	})
	if err != nil {
		t.Fatalf("failed to parse POST /create request body: %v", err)
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, parsedUrl.String(), bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("failed to create POST /create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make POST /create request: %v", err)
	}

	return resp
}

func getLongUrl(t *testing.T, baseUrl, key string) *http.Response {
	t.Helper()

	urlString := fmt.Sprintf("%s/%s", baseUrl, key)
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		t.Fatalf("failed to parse GET /%s url: %v", key, err)
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, parsedUrl.String(), http.NoBody)
	if err != nil {
		t.Fatalf("failed to create GET /%s request: %v", key, err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// stop following redirects
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make GET /%s request: %v", key, err)
	}

	return resp
}

func parseJSONResponse[T any](t *testing.T, resp *http.Response) *T {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var result T
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v. Body: %s", err, string(body))
	}

	return &result
}
