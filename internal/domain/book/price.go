package book

import "errors"

type Price struct {
	value float64
}

func NewPrice(v float64) (Price, error) {
	if v <= 0 {
		return Price{}, errors.New("preço inválido: deve ser maior que zero")
	}
	return Price{value: v}, nil
}

func (p Price) Float64() float64 { return p.value }
