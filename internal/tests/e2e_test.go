package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type CreateResponse struct {
	Url string `json:"url"`
}

func TestIntegration(t *testing.T) {
	_ = godotenv.Load(".env.test")
	ctx := context.Background()

	testResources, err := setupTestResources(ctx, t)
	if err != nil {
		fmt.Printf("setup tests failed: %s", err)
		testResources.Cleanup(ctx, t)
		os.Exit(1)
	}

	t.Cleanup(func() {
		testResources.Cleanup(ctx, t)
	})

	t.Run("redirects to original url by short url", func(t *testing.T) {
		apiDomain := os.Getenv("API_DOMAIN")
		newKey := "abcde123"
		longUrl := "https://myveryveryveryveryveryveryveryveryveryveryveryveryveryveryverylongurl.com"

		rows, err := testResources.KeyDbPool.Query(
			ctx,
			`INSERT INTO keys (key_value) VALUES ($1)`,
			newKey,
		)
		if err != nil {
			t.Fatalf("failed to insert new key: %v", err)
		}
		defer rows.Close()

		createResp := createShortUrl(t, testResources.AppUrl, longUrl)
		defer createResp.Body.Close() //nolint:errcheck

		if createResp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", createResp.StatusCode)
		}

		body, err := io.ReadAll(createResp.Body)
		if err != nil {
			t.Fatalf("failed to read POST /create response body: %v", err)
		}

		var createRes CreateResponse
		err = json.Unmarshal(body, &createRes)
		if err != nil {
			t.Fatalf("failed to parse POST /create JSON response: %v. Body: %s", err, string(body))
		}

		expectedUrl := fmt.Sprintf("%s/%s", apiDomain, newKey)
		if createRes.Url != expectedUrl {
			t.Errorf("expected url '%s', got '%s'", expectedUrl, createRes.Url)
		}

		getResp := getLongUrl(t, testResources.AppUrl, newKey)
		if getResp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("expected 301 redirect, got %d", getResp.StatusCode)
		}

		location := getResp.Header.Get("Location")
		if location != longUrl {
			t.Fatalf("expected redirect to %s, got %s", longUrl, location)
		}
	})

	t.Run("returns url not found error", func(t *testing.T) {
		nonExistentKey := "1234bcde"
		resp := getLongUrl(t, testResources.AppUrl, nonExistentKey)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns validation error with invalid key", func(t *testing.T) {
		invalidKey := "1234bcd"
		resp := getLongUrl(t, testResources.AppUrl, invalidKey)

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", resp.StatusCode)
		}
	})
}

func createShortUrl(t *testing.T, baseUrl string, longUrl string) *http.Response {
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

func getLongUrl(t *testing.T, baseUrl string, key string) *http.Response {
	t.Helper()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// stop following redirects
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(fmt.Sprintf("%s/%s", baseUrl, key))
	if err != nil {
		t.Fatalf("failed to make GET /%s request: %v", key, err)
	}
	defer resp.Body.Close()

	return resp
}
