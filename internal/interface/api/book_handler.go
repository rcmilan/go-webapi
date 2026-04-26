package api

import (
	"bookstore-api/internal/application"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// bookService is the local interface the handler depends on.
// Defined here, consumed here — idiomatic Go.
type bookService interface {
	RegisterBook(ctx context.Context, cmd application.RegisterBookCommand) (*application.BookResult, error)
	GetBook(ctx context.Context, id uint32) (*application.BookResult, error)
	ListBooks(ctx context.Context, filter application.BookFilter) ([]*application.BookResult, error)
}

type BookHandler struct {
	service bookService
}

func NewBookHandler(s *application.BookService) *BookHandler {
	return &BookHandler{service: s}
}

// ── DTOs ────────────────────────────────────────────────────────────────────

type bookDTO struct {
	ID          uint32  `json:"id"`
	Title       string  `json:"title"`
	ISBN        string  `json:"isbn"`
	Price       float64 `json:"price"`
	ReleaseYear int     `json:"release_year"`
}

type createBookInput struct {
	Body struct {
		Title       string  `json:"title" doc:"Book title" minLength:"1"`
		ISBN        string  `json:"isbn" doc:"13-digit ISBN" minLength:"13" maxLength:"13"`
		Price       float64 `json:"price" doc:"Book price" minimum:"0.01"`
		ReleaseYear int     `json:"release_year" doc:"Year the book was released"`
	}
}

type createBookOutput struct{ Body bookDTO }

type getBookInput struct {
	ID uint32 `path:"id" doc:"Book ID"`
}

type getBookOutput struct{ Body bookDTO }

type listBooksInput struct {
	ID   uint32 `query:"id" doc:"Filter by book ID"`
	ISBN string `query:"isbn" doc:"Filter by ISBN"`
}

type listBooksOutput struct {
	Body struct {
		Books []bookDTO `json:"books"`
	}
}

// ── Mappers ──────────────────────────────────────────────────────────────────

func toBookDTO(r *application.BookResult) bookDTO {
	return bookDTO{ID: r.ID, Title: r.Title, ISBN: r.ISBN, Price: r.Price, ReleaseYear: r.ReleaseYear}
}

func toBookDTOs(results []*application.BookResult) []bookDTO {
	dtos := make([]bookDTO, len(results))
	for i, r := range results {
		dtos[i] = toBookDTO(r)
	}
	return dtos
}

// ── Handlers ─────────────────────────────────────────────────────────────────

func (h *BookHandler) createBook(ctx context.Context, input *createBookInput) (*createBookOutput, error) {
	result, err := h.service.RegisterBook(ctx, application.RegisterBookCommand{
		Title:       input.Body.Title,
		ISBN:        input.Body.ISBN,
		Price:       input.Body.Price,
		ReleaseYear: input.Body.ReleaseYear,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &createBookOutput{Body: toBookDTO(result)}, nil
}

func (h *BookHandler) getBook(ctx context.Context, input *getBookInput) (*getBookOutput, error) {
	result, err := h.service.GetBook(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &getBookOutput{Body: toBookDTO(result)}, nil
}

func (h *BookHandler) listBooks(ctx context.Context, input *listBooksInput) (*listBooksOutput, error) {
	results, err := h.service.ListBooks(ctx, application.BookFilter{ID: input.ID, ISBN: input.ISBN})
	if err != nil {
		return nil, mapError(err)
	}
	out := &listBooksOutput{}
	out.Body.Books = toBookDTOs(results)
	return out, nil
}

// ── Routing ──────────────────────────────────────────────────────────────────

func (h *BookHandler) Register(api huma.API) {
	huma.Register(api, op("create-book", http.MethodPost, "/api/v1/books", "Register a new book", "Books"), h.createBook)
	huma.Register(api, op("get-book", http.MethodGet, "/api/v1/books/{id}", "Get a book by ID", "Books"), h.getBook)
	huma.Register(api, op("list-books", http.MethodGet, "/api/v1/books", "List books", "Books"), h.listBooks)
}

func op(id, method, path, summary string, tags ...string) huma.Operation {
	return huma.Operation{OperationID: id, Method: method, Path: path, Summary: summary, Tags: tags}
}
