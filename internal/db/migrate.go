package db

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func Migrate() {
	migrationsDir := "/app/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations dir: %v", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	ctx := context.Background()
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}
		path := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file.Name(), err)
		}

		if _, err := DB.Exec(ctx, string(content)); err != nil {
			log.Fatalf("Failed to run migration %s: %v", file.Name(), err)
		}
		log.Printf("Connecting to DB: host=%s dbname=%s user=%s", os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"))

		log.Printf("Migrated migration %s", file.Name())
	}
}
