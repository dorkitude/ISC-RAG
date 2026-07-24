package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func openDatabase(cmdProjectDir string) (*pgxpool.Pool, error) {
	projectDir, err := filepath.Abs(cmdProjectDir)
	if err != nil {
		return nil, err
	}
	if err := loadLocalEnv(projectDir); err != nil {
		return nil, err
	}
	url := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is required; start Postgres with docker compose up -d and set it in .env")
	}
	pool, err := storage.Open(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
