package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sort"
)

func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return pgxpool.New(ctx, url)
}
func Migrate(ctx context.Context, p *pgxpool.Pool) error {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	for _, f := range files {
		b, e := os.ReadFile("migrations/" + f.Name())
		if e != nil {
			return e
		}
		if _, e = p.Exec(ctx, string(b)); e != nil {
			return fmt.Errorf("migration %s: %w", f.Name(), e)
		}
	}
	return nil
}
