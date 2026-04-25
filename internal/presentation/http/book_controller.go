// internal/presentation/http/book_controller.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBookRequest struct {
	Title string  `json:"title" binding:"required"`
	ISBN  string  `json:"isbn" binding:"required,len=13"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

// @Summary Criar Livro
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) Create(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validação falhou: campos obrigatórios ausentes"})
		return
	}
	// Chamar application service...
	c.JSON(http.StatusCreated, gin.H{"message": "Livro cadastrado"})
}
