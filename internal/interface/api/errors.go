package api

import (
	"bookstore-api/internal/application"
	"errors"

	"github.com/danielgtaylor/huma/v2"
)

func mapError(err error) error {
	var validation application.ErrValidation
	if errors.As(err, &validation) {
		return huma.Error400BadRequest(validation.Error())
	}
	var notFound application.ErrNotFound
	if errors.As(err, &notFound) {
		return huma.Error404NotFound(notFound.Error())
	}
	return huma.Error500InternalServerError("erro interno")
}
