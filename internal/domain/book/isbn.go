package book

import (
	"database/sql/driver"
	"errors"
)

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

func (i ISBN) Value() (driver.Value, error) {
	return i.value, nil
}

func (i *ISBN) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	v, ok := value.(string)
	if !ok {
		return errors.New("falha ao converter valor do banco para ISBN")
	}
	i.value = v
	return nil
}
