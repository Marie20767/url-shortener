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
	KeyDbPool    *pgxpool.Pool
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

	keyDbUrl, err := getKeyDbUrl(ctx, composeStack)
	if err != nil {
		return nil, err
	}

	keyDbPool, err := pgxpool.New(ctx, keyDbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new key db pool: %w", err)
	}

	if err := keyDbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to key db: %w", err)
	}

	return &TestResources{
		ComposeStack: composeStack,
		AppUrl:       appUrl,
		KeyDbPool:    keyDbPool,
	}, nil
}

func (tr *TestResources) Cleanup(ctx context.Context, t *testing.T) {
	if tr == nil {
		return
	}

	if tr.KeyDbPool != nil {
		tr.KeyDbPool.Close()
	}

	if tr.ComposeStack != nil {
		err := tr.ComposeStack.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
		if err != nil {
			t.Logf("failed to tear down compose stack: %v", err)
		}
	}
}

func getKeyDbUrl(ctx context.Context, composeStack *compose.DockerCompose) (string, error) {
	keyDbContainer, err := composeStack.ServiceContainer(ctx, "postgres")
	if err != nil {
		return "", fmt.Errorf("failed to get key db container: %w", err)
	}

	keyDbPort, err := keyDbContainer.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("failed to get key db mapped port: %w", err)
	}

	keyDbHost, err := keyDbContainer.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get key db host: %w", err)
	}

	return fmt.Sprintf(
		"postgres://testuser:password@%s:%s/keydb?sslmode=disable",
		keyDbHost,
		keyDbPort.Port(),
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
