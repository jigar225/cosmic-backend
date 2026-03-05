package main
import (
"fmt"
"github.com/golang-migrate/migrate/v4"
_ "github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
)
func main() {
  _, err := migrate.New("file://internal/migrations", "postgres://a:b@c:5432/d")
  fmt.Println("file://internal/migrations", err)
}