package application

import (
	"bookstore-api/internal/domain/book"
	"context"
)

type BookService struct {
	repo book.Repository
}

func NewBookService(repo book.Repository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) RegisterBook(ctx context.Context, cmd RegisterBookCommand) (*BookResult, error) {
	var result *BookResult
	err := s.repo.Atomic(ctx, func(r book.Repository) error {
		isbn, err := book.NewISBN(cmd.ISBN)
		if err != nil {
			return ErrValidation{Msg: err.Error()}
		}
		b, err := book.NewBook(cmd.Title, isbn, cmd.Price)
		if err != nil {
			return ErrValidation{Msg: err.Error()}
		}
		if err := r.Create(ctx, b); err != nil {
			return err
		}
		result = toResult(b)
		return nil
	})
	return result, err
}

func (s *BookService) GetBook(ctx context.Context, id uint32) (*BookResult, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound{Msg: "livro não encontrado"}
	}
	return toResult(b), nil
}

func (s *BookService) ListBooks(ctx context.Context, filter BookFilter) ([]*BookResult, error) {
	books, err := s.repo.List(ctx, book.Filter{ID: filter.ID, ISBN: filter.ISBN})

	if err != nil {
		return nil, err
	}
	results := make([]*BookResult, len(books))
	for i, b := range books {
		results[i] = toResult(b)
	}
	return results, nil
}

func toResult(b *book.Book) *BookResult {
	return &BookResult{ID: b.ID, Title: b.Title, ISBN: b.ISBN.String(), Price: b.Price}
}
