package book

import "context"

type Filter struct {
	ID   uint32
	ISBN string
}

type Repository interface {
	Create(ctx context.Context, b *Book) error
	FindByID(ctx context.Context, id uint32) (*Book, error)
	List(ctx context.Context, filter Filter) ([]*Book, error)
	Atomic(ctx context.Context, fn func(Repository) error) error
}
