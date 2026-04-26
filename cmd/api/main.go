package main

import (
	"bookstore-api/internal/application"
	"bookstore-api/internal/infrastructure/persistence"
	"bookstore-api/internal/interface/api"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "root:root_password@tcp(127.0.0.1:3306)/bookstore_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("falha ao conectar no MySQL: %v", err)
	}

	repo, err := persistence.NewGormBookRepository(db)
	if err != nil {
		log.Fatalf("falha ao inicializar repositório: %v", err)
	}

	service := application.NewBookService(repo)
	bookHandler := api.NewBookHandler(service)
	recordHandler := api.NewRecordHandler()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	config := huma.DefaultConfig("Bookstore API", "1.0.0")
	config.Info.Description = "API de livraria seguindo DDD com GORM."
	config.SchemasPath = ""
	config.CreateHooks = nil
	humaAPI := humagin.New(r, config)

	bookHandler.Register(humaAPI)
	recordHandler.Register(humaAPI)

	log.Printf("Servidor iniciado na porta :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
	}
}
