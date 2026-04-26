package application

import "bookstore-api/internal/domain/book"

type RegisterBookCommand struct {
	Title       string
	ISBN        string
	Price       float64
	ReleaseYear int
}

type BookFilter struct {
	ID   book.BookID
	ISBN string
}

type BookResult struct {
	ID          book.BookID
	Title       string
	ISBN        string
	Price       float64
	ReleaseYear int
}
