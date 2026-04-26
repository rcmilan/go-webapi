package book

import "errors"

type Book struct {
	ID          BookID
	Title       string
	ISBN        ISBN
	Price       Price
	ReleaseYear int
}

func NewBook(title string, isbn ISBN, price Price, releaseYear int) (*Book, error) {
	if title == "" {
		return nil, errors.New("título é obrigatório")
	}
	return &Book{Title: title, ISBN: isbn, Price: price, ReleaseYear: releaseYear}, nil
}
