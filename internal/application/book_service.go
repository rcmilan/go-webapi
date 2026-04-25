// internal/application/book_service.go
package application

import (
	"bookstore-api/internal/domain/book"
	"context"
)

type BookService struct {
	repo book.Repository
}

func (s *BookService) RegisterBook(ctx context.Context, title, isbnStr string, price float64) error {
	return s.repo.Atomic(ctx, func(r book.Repository) error {
		isbn, _ := book.NewISBN(isbnStr)
		newBook, _ := book.NewBook(title, isbn, price)

		return r.Create(ctx, newBook)
	})
}
