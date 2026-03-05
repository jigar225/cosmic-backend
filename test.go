package main
import (
"fmt"
"os"
"github.com/golang-migrate/migrate/v4"
_ "github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
)
func main() {
  url := os.Getenv("DATABASE_URL")
  m, err := migrate.New("file://internal/migrations", url)
  fmt.Println(m, err)
}