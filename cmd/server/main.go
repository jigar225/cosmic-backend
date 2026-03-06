// Package main starts the Fiber HTTP server: connects to Postgres, runs migrations, then serves.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"back_testing/internal/repository"
	"back_testing/internal/storage"
	"back_testing/internal/transport"
	"back_testing/internal/transport/handlers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: loading .env: %v", err)
	}

	ctx := context.Background()
	pool, err := repository.NewDBPool(ctx)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	if err := repository.RunMigrations(); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("Migrations applied.")

	s3Uploader, err := storage.NewS3Uploader(storage.S3FromEnv())
	if err != nil {
		log.Printf("warning: S3 uploader init failed (chapter PDF uploads disabled): %v", err)
		s3Uploader, _ = storage.NewS3Uploader(storage.S3Config{}) // no-op when not configured
	}

	h := &handlers.Handlers{
		BoardRepo:        repository.NewBoardRepo(pool),
		CountryRepo:      repository.NewCountryRepo(pool),
		GradeMethodRepo:  repository.NewGradeMethodRepo(pool),
		GradeRepo:        repository.NewGradeRepo(pool),
		MediumRepo:       repository.NewMediumRepo(pool),
		LanguageRepo:     repository.NewLanguageRepo(pool),
		SubjectRepo:      repository.NewSubjectRepo(pool),
		BookRepo:         repository.NewBookRepo(pool),
		ChapterRepo:      repository.NewChapterRepo(pool),
		S3Uploader:       s3Uploader,
		UserRepo:         repository.NewUserRepo(pool),
		RefreshTokenRepo: repository.NewRefreshTokenRepo(pool),
	}
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		h.AuthConfig = &handlers.AuthConfig{
			JWTSecret:     []byte(secret),
			AccessExpiry:  3 * time.Minute,
			RefreshExpiry: 30 * 24 * time.Hour, // 30 days
		}
	}
	app := transport.NewApp(h)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("listen: %v", err)
		}
	}()
	log.Println("Server started. DB connected. Fiber HTTP server on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
