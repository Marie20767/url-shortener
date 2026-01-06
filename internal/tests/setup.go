package tests

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

type TestResources struct {
	ComposeStack *compose.DockerCompose
	DbPool       *pgxpool.Pool
	AppUrl       string
}

// setupTestResources creates and starts all required containers for testing
func setupTestResources(ctx context.Context, t *testing.T) (*TestResources, error) {
	t.Helper()
	composeStack, err := compose.NewDockerCompose("./docker-compose.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create compose stack: %w", err)
	}

	err = composeStack.Up(ctx, compose.Wait(true))
	defer func() {
		// handle cleanup here if setup fails halfway through
		if err != nil {
			cleanupErr := composeStack.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
			slog.Error("cleanup error", slog.Any("error", cleanupErr))
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("failed to start compose stack: %w", err)
	}

	appUrl, err := getAppUrl(ctx, composeStack)
	if err != nil {
		return nil, err
	}

	urldbUrl, err := getDbUrl(ctx, composeStack)
	if err != nil {
		return nil, err
	}

	urldbPool, err := pgxpool.New(ctx, urldbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new db pool: %w", err)
	}

	if err := urldbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &TestResources{
		ComposeStack: composeStack,
		AppUrl:       appUrl,
		DbPool:       urldbPool,
	}, nil
}

func (tr *TestResources) Cleanup(ctx context.Context, t *testing.T) {
	if tr == nil {
		return
	}

	if tr.DbPool != nil {
		tr.DbPool.Close()
	}

	if tr.ComposeStack != nil {
		err := tr.ComposeStack.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
		if err != nil {
			t.Logf("failed to tear down compose stack: %v", err)
		}
	}
}

func getDbUrl(ctx context.Context, composeStack *compose.DockerCompose) (string, error) {
	urldbContainer, err := composeStack.ServiceContainer(ctx, "postgres")
	if err != nil {
		return "", fmt.Errorf("failed to get db container: %w", err)
	}

	urldbPort, err := urldbContainer.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("failed to get db mapped port: %w", err)
	}

	urldbHost, err := urldbContainer.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get db host: %w", err)
	}

	return fmt.Sprintf(
		"postgres://testuser:password@%s:%s/urldb?sslmode=disable",
		urldbHost,
		urldbPort.Port(),
	), nil
}

func getAppUrl(ctx context.Context, composeStack *compose.DockerCompose) (string, error) {
	appContainer, err := composeStack.ServiceContainer(ctx, "url-shortener-server")
	if err != nil {
		return "", fmt.Errorf("failed to get app container: %w", err)
	}

	appPort, err := appContainer.MappedPort(ctx, "3001")
	if err != nil {
		return "", fmt.Errorf("failed to get app mapped port : %w", err)
	}

	appHost, err := appContainer.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get app host: %w", err)
	}

	return fmt.Sprintf("http://%s:%s", appHost, appPort.Port()), nil
}
