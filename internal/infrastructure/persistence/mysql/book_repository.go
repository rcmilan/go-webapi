// internal/infrastructure/persistence/mysql/book_repository.go
package mysql

import (
	"bookstore-api/internal/domain/book"
	"context"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func (r *BookRepository) Atomic(ctx context.Context, fn func(repo book.Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Criamos uma nova instância do repo usando a transação (tx)
		return fn(&BookRepository{db: tx})
	})
}

func (r *BookRepository) Create(ctx context.Context, b *book.Book) error {
	// Mapeamento: Domínio -> Persistência
	model := BookModel{Title: b.Title, ISBN: b.ISBN.String(), Price: b.Price}
	return r.db.WithContext(ctx).Create(&model).Error
}
