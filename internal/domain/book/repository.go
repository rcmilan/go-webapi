// internal/domain/book/repository.go
package book

import "context"

type Repository interface {
	Create(ctx context.Context, b *Book) error
	Atomic(ctx context.Context, fn func(repo Repository) error) error
}
