package tests

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var testResources *TestResources

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env.test")
	ctx := context.Background()

	var err error
	testResources, err = setupTestResources(ctx, &testing.T{})
	if err != nil {
		fmt.Printf("setup tests failed: %s\n", err)
		if testResources != nil {
			testResources.Cleanup(ctx, &testing.T{})
		}
		os.Exit(1)
	}

	exitCode := m.Run()

	testResources.Cleanup(ctx, &testing.T{})
	if exitCode != 0 {
		slog.Error("tests failed", slog.Int("exit_code", exitCode))
	}

	os.Exit(exitCode)
}

type CreateResponse struct {
	Url string `json:"url"`
}

func TestUrl(t *testing.T) {
	t.Run("redirects to original url by short url", func(t *testing.T) {
		apiDomain := os.Getenv("API_DOMAIN")
		newKey := "aBcdE123"
		longUrl := "https://myveryveryveryveryveryveryveryveryveryveryveryveryveryveryverylongurl.com"

		_, err := testResources.DbPool.Exec(
			t.Context(),
			`INSERT INTO keys (id) VALUES ($1)`,
			newKey,
		)
		if err != nil {
			t.Fatalf("failed to insert new key: %v", err)
		}

		createResp := createShortUrl(t, testResources.AppUrl, longUrl)
		defer createResp.Body.Close() //nolint:errcheck

		if createResp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", createResp.StatusCode)
		}

		actualCreateRes := parseJSONResponse[CreateResponse](t, createResp)

		expectedUrl := fmt.Sprintf("%s/%s", apiDomain, newKey)
		if actualCreateRes.Url != expectedUrl {
			t.Errorf("expected url '%s', got '%s'", expectedUrl, actualCreateRes.Url)
		}

		getResp := getLongUrl(t, testResources.AppUrl, newKey)
		defer getResp.Body.Close() //nolint:errcheck
		if getResp.StatusCode != http.StatusFound {
			t.Fatalf("expected 302 redirect, got %d", getResp.StatusCode)
		}

		location := getResp.Header.Get("Location")
		if location != longUrl {
			t.Fatalf("expected redirect to %s, got %s", longUrl, location)
		}
	})

	t.Run("returns url not found error", func(t *testing.T) {
		nonExistentKey := "1234bcde"
		resp := getLongUrl(t, testResources.AppUrl, nonExistentKey)
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns validation error with invalid key", func(t *testing.T) {
		invalidKey := "1234Bcd"
		resp := getLongUrl(t, testResources.AppUrl, invalidKey)
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", resp.StatusCode)
		}
	})
}
