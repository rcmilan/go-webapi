package book

import "context"

type Filter struct {
	ID   BookID
	ISBN string
}

type Repository interface {
	Create(ctx context.Context, b *Book) error
	FindByID(ctx context.Context, id BookID) (*Book, error)
	List(ctx context.Context, filter Filter) ([]*Book, error)
	Atomic(ctx context.Context, fn func(Repository) error) error
}
