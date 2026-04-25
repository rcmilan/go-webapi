// cmd/api/main.go
package main

import (
	"bookstore-api/internal/application"
	"bookstore-api/internal/infrastructure/persistence/mysql"
	"bookstore-api/internal/presentation/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Conexão com Banco
	dsn := "root:root@tcp(127.0.0.1:3306)/livraria?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 2. Injeção de Dependências (Wired)
	repo := mysql.NewBookRepository(db)
	service := application.NewBookService(repo)
	handler := http.NewBookHandler(service)

	// 3. Roteamento
	r := gin.Default()
	r.POST("/books", handler.Register)

	// 4. Documentação Scalar
	r.GET("/docs", func(c *gin.Context) {
		html, _ := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL:  "./docs/swagger.json",
			DarkMode: true,
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", byte(html))
	})

	r.Run(":8080")
}
