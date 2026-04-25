// internal/domain/book/book.go
package book

import "errors"

type Book struct {
	ID    uint32
	Title string
	ISBN  ISBN
	Price float64
}

// Factory para garantir invariantes de negócio
func NewBook(title string, isbn ISBN, price float64) (*Book, error) {
	if price <= 0 {
		return nil, errors.New("preço inválido")
	}
	return &Book{Title: title, ISBN: isbn, Price: price}, nil
}
