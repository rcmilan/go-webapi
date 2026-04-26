// cmd/migrate syncs the Ent schema into bookstore_dev so Atlas can use it
// as the desired state when generating versioned migration files.
//
// Run before every `atlas migrate diff`:
//
//	go run ./cmd/migrate
package main

import (
	"bookstore-api/ent"
	"bookstore-api/ent/migrate"
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const devDSN = "root:root_password@tcp(localhost:3306)/bookstore_dev?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	client, err := ent.Open("mysql", devDSN)
	if err != nil {
		log.Fatalf("failed connecting to bookstore_dev: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		log.Fatalf("failed syncing schema to bookstore_dev: %v", err)
	}

	log.Println("bookstore_dev synced with Ent schema — ready for: atlas migrate diff <name> --env local")
}
