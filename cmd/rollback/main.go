// Package main runs the last migration down (rollback).
package main

import (
	"log"
	"os"

	"back_testing/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: loading .env: %v", err)
	}
	if err := repository.RollbackMigrations(); err != nil {
		log.Fatalf("rollback: %v", err)
	}
	log.Println("Rollback completed.")
}
