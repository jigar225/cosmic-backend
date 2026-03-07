package main
import (
"database/sql"
"fmt"
_ "github.com/lib/pq"
)
func main() {
  db, err := sql.Open("postgres", "postgresql://postgres:admin123@localhost:5432/cosmicq?sslmode=disable")
  if err != nil { panic(err) }
  rows, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_name = 'mediums'")
  if err != nil { panic(err) }
  for rows.Next() { var n string; rows.Scan(&n); fmt.Println(n) }
}