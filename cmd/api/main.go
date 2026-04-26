package main

import (
	"bookstore-api/ent"
	"bookstore-api/internal/application"
	"bookstore-api/internal/infrastructure/persistence"
	"bookstore-api/internal/interface/api"
	"bookstore-api/internal/observability"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	observability.Setup()

	dsn := "root:root_password@tcp(127.0.0.1:3306)/bookstore_db?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		slog.Error("falha ao conectar no MySQL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer client.Close()

	repo := persistence.NewEntBookRepository(client)
	service := application.NewBookService(repo)
	bookHandler := api.NewBookHandler(service)
	recordHandler := api.NewRecordHandler()

	r := gin.New()
	r.Use(api.CorrelationID(), api.RequestLogger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	config := huma.DefaultConfig("Bookstore API", "1.0.0")
	config.Info.Description = "API de livraria seguindo DDD com Ent ORM."
	config.SchemasPath = ""
	config.CreateHooks = nil
	humaAPI := humagin.New(r, config)

	bookHandler.Register(humaAPI)
	recordHandler.Register(humaAPI)

	slog.Info("servidor iniciado", slog.String("addr", ":8080"))
	if err := r.Run(":8080"); err != nil {
		slog.Error("falha ao iniciar servidor", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
