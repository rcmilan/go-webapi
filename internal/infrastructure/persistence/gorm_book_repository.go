package persistence

import (
	"bookstore-api/internal/domain/book"
	"context"
	"errors"

	"gorm.io/gorm"
)

type bookDPO struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"size:255;not null"`
	ISBN      book.ISBN `gorm:"type:varchar(13);uniqueIndex;not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt int64     `gorm:"autoCreateTime"`
}

func (bookDPO) TableName() string { return "books" }

func toDomain(d bookDPO) *book.Book {
	return &book.Book{ID: d.ID, Title: d.Title, ISBN: d.ISBN, Price: d.Price}
}

type GormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) (*GormBookRepository, error) {
	if err := db.AutoMigrate(&bookDPO{}); err != nil {
		return nil, err
	}
	return &GormBookRepository{db: db}, nil
}

func (r *GormBookRepository) Create(ctx context.Context, b *book.Book) error {
	dpo := bookDPO{Title: b.Title, ISBN: b.ISBN, Price: b.Price}
	if err := r.db.WithContext(ctx).Create(&dpo).Error; err != nil {
		return err
	}
	b.ID = dpo.ID
	return nil
}

func (r *GormBookRepository) FindByID(ctx context.Context, id uint32) (*book.Book, error) {
	var dpo bookDPO
	err := r.db.WithContext(ctx).First(&dpo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(dpo), nil
}

func (r *GormBookRepository) List(ctx context.Context, filter book.Filter) ([]*book.Book, error) {
	var dpos []bookDPO
	q := r.db.WithContext(ctx)
	if filter.ID != 0 {
		q = q.Where("id = ?", filter.ID)
	}
	if filter.ISBN != "" {
		q = q.Where("isbn = ?", filter.ISBN)
	}
	if err := q.Find(&dpos).Error; err != nil {
		return nil, err
	}
	books := make([]*book.Book, len(dpos))
	for i, d := range dpos {
		books[i] = toDomain(d)
	}
	return books, nil
}

func (r *GormBookRepository) Atomic(ctx context.Context, fn func(book.Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&GormBookRepository{db: tx})
	})
}
