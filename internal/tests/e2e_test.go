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

		resp := createShortUrl(t, testResources.AppUrl, longUrl)
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read POST /create response body: %v", err)
		}

		var res CreateResponse
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("failed to parse POST /create JSON response: %v. Body: %s", err, string(body))
		}

		expectedUrl := fmt.Sprintf("%s/%s", apiDomain, newKey)
		if res.Url != expectedUrl {
			t.Errorf("expected url '%s', got '%s'", expectedUrl, res.Url)
		}

		// TODO: make a GET /{$key} request with the key and check you get a 302 redirect to the longURL
	})

	// t.Run("returns url not found error", func(t *testing.T) {

	// 	resp := getLongUrl(t, testResources.AppUrl, nonExistentID)
	// 	defer resp.Body.Close() //nolint:errcheck

	// 	if resp.StatusCode != http.StatusNotFound {
	// 		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	// 	}
	// })
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
