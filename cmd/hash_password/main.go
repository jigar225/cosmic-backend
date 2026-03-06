package main

import (
	"fmt"
	"log"

	"back_testing/internal/auth"
)

func main() {
	const pwd = "cosmicq@123"
	hash, err := auth.HashPassword(pwd)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	fmt.Println("hash:", hash)
}
