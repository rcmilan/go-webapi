package api

import (
	"bookstore-api/internal/application"
	"bookstore-api/internal/domain/book"
	"bookstore-api/internal/observability"
	"context"
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

func mapError(ctx context.Context, err error) error {
	var validation application.ErrValidation
	if errors.As(err, &validation) {
		return huma.Error400BadRequest(validation.Error())
	}
	var notFound application.ErrNotFound
	if errors.As(err, &notFound) {
		return huma.Error404NotFound(notFound.Error())
	}
	var conflict book.ErrConflict
	if errors.As(err, &conflict) {
		return huma.Error409Conflict(conflict.Error())
	}
	observability.FromContext(ctx).ErrorContext(ctx, "unhandled error", slog.String("error", err.Error()))
	return huma.Error500InternalServerError("erro interno")
}
