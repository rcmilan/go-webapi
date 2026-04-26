package book

import "errors"

type Book struct {
	ID          uint32
	Title       string
	ISBN        ISBN
	Price       float64
	ReleaseYear int
}

func NewBook(title string, isbn ISBN, price float64, releaseYear int) (*Book, error) {
	if title == "" {
		return nil, errors.New("título é obrigatório")
	}
	if price <= 0 {
		return nil, errors.New("preço inválido: deve ser maior que zero")
	}
	return &Book{Title: title, ISBN: isbn, Price: price, ReleaseYear: releaseYear}, nil
}
