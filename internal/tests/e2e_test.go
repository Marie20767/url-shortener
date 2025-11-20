package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
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

	t.Run("redirects to long url by key", func(t *testing.T) {
		// longUrl := "https://myveryveryveryveryveryveryveryveryveryveryveryveryveryveryverylongurl.com"
		
		// TODO: insert 1 key into key db
		// TODO: make a POST /create request and check you get a 201 with the correct shortURL
		// TODO: make a GET /{$key} request with the key and check you get a 302 redirect to the longURL

		// resp := getCourse(t, testResources.AppUrl, id)
		// defer resp.Body.Close() //nolint:errcheck

		// if resp.StatusCode != http.StatusOK {
		// 	t.Fatalf("expected status 200, got %d", resp.StatusCode)
		// }

		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	t.Fatalf("failed to read response body: %v", err)
		// }

		// var result sqlc.Course
		// err = json.Unmarshal(body, &result)
		// if err != nil {
		// 	t.Fatalf("failed to parse JSON response: %v. Body: %s", err, string(body))
		// }

		// if result.Title.String != expectedTitle {
		// 	t.Errorf("expected title '%s', got '%s'", expectedTitle, result.Title.String)
		// }

		// if result.Description.String != expectedDescription {
		// 	t.Errorf("expected description '%s', got '%s'", expectedDescription, result.Description.String)
		// }

		// if result.ID.String() != id.String() {
		// 	t.Errorf("expected id '%s', got '%s'", id.String(), result.ID)
		// }
	})

	// t.Run("returns not found error", func(t *testing.T) {
	// 	nonExistentID := uuid.New()

	// 	resp := getLongUrl(t, testResources.AppUrl, nonExistentID)
	// 	defer resp.Body.Close() //nolint:errcheck

	// 	if resp.StatusCode != http.StatusNotFound {
	// 		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	// 	}
	// })
}