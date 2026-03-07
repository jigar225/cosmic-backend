package main
import (
"fmt"
"github.com/golang-migrate/migrate/v4"
_ "github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
)
func main() {
  m, err := migrate.New("file://internal/migrations", "postgresql://postgres:admin123@localhost:5432/cosmicq?sslmode=disable")
  if err != nil { fmt.Println(err); return }
  err = m.Force(5)
  fmt.Println("Forced 5:", err)
}