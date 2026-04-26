package persistence

import (
	"bookstore-api/ent"
	entbook "bookstore-api/ent/book"
	"bookstore-api/internal/domain/book"
	"context"
	"fmt"
)

type EntBookRepository struct {
	client *ent.Client
}

func NewEntBookRepository(client *ent.Client) *EntBookRepository {
	return &EntBookRepository{client: client}
}

func (r *EntBookRepository) Create(ctx context.Context, b *book.Book) error {
	created, err := r.client.Book.
		Create().
		SetTitle(b.Title).
		SetIsbn(b.ISBN.String()).
		SetPrice(b.Price.Float64()).
		SetReleaseYear(b.ReleaseYear).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return book.ErrConflict{Msg: "ISBN já cadastrado"}
		}
		return fmt.Errorf("ent: create book: %w", err)
	}
	b.ID = book.BookID(created.ID)
	return nil
}

func (r *EntBookRepository) FindByID(ctx context.Context, id book.BookID) (*book.Book, error) {
	e, err := r.client.Book.Get(ctx, uint32(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("ent: find book by id: %w", err)
	}
	return toDomainBook(e)
}

func (r *EntBookRepository) List(ctx context.Context, filter book.Filter) ([]*book.Book, error) {
	q := r.client.Book.Query()
	if filter.ID != 0 {
		q = q.Where(entbook.ID(uint32(filter.ID)))
	}
	if filter.ISBN != "" {
		q = q.Where(entbook.Isbn(filter.ISBN))
	}
	rows, err := q.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ent: list books: %w", err)
	}
	result := make([]*book.Book, len(rows))
	for i, e := range rows {
		b, err := toDomainBook(e)
		if err != nil {
			return nil, err
		}
		result[i] = b
	}
	return result, nil
}

func (r *EntBookRepository) Atomic(ctx context.Context, fn func(book.Repository) error) error {
	tx, err := r.client.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ent: begin transaction: %w", err)
	}
	txRepo := &EntBookRepository{client: tx.Client()}
	if err := fn(txRepo); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func toDomainBook(e *ent.Book) (*book.Book, error) {
	isbn, err := book.NewISBN(e.Isbn)
	if err != nil {
		return nil, fmt.Errorf("ent: invalid isbn in database %q: %w", e.Isbn, err)
	}
	price, err := book.NewPrice(e.Price)
	if err != nil {
		return nil, fmt.Errorf("ent: invalid price in database %v: %w", e.Price, err)
	}
	return &book.Book{
		ID:          book.BookID(e.ID),
		Title:       e.Title,
		ISBN:        isbn,
		Price:       price,
		ReleaseYear: e.ReleaseYear,
	}, nil
}
