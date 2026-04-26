package book

import "errors"

type ISBN struct {
	value string
}

func NewISBN(v string) (ISBN, error) {
	if len(v) != 13 {
		return ISBN{}, errors.New("ISBN inválido: deve possuir exatamente 13 dígitos")
	}
	return ISBN{value: v}, nil
}

func (i ISBN) String() string { return i.value }
