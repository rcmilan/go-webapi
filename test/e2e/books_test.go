package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBook_success(t *testing.T) {
	cleanBooks(t)

	resp, body := doJSON(t, http.MethodPost, "/api/v1/books", map[string]any{
		"title":        "The Go Programming Language",
		"isbn":         "9780134190440",
		"price":        49.99,
		"release_year": 2015,
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "The Go Programming Language", body["title"])
	assert.Equal(t, "9780134190440", body["isbn"])
	assert.Equal(t, 49.99, body["price"])
	assert.Equal(t, float64(2015), body["release_year"])
	assert.NotZero(t, body["id"])
}

func TestCreateBook_duplicateISBN(t *testing.T) {
	cleanBooks(t)

	payload := map[string]any{
		"title":        "Clean Code",
		"isbn":         "9780132350884",
		"price":        39.99,
		"release_year": 2008,
	}
	_, _ = doJSON(t, http.MethodPost, "/api/v1/books", payload)
	resp, _ := doJSON(t, http.MethodPost, "/api/v1/books", payload)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestCreateBook_invalidISBN(t *testing.T) {
	cleanBooks(t)

	// ISBN shorter than 13 chars — Huma rejects at validation layer.
	resp, _ := doJSON(t, http.MethodPost, "/api/v1/books", map[string]any{
		"title":        "Bad Book",
		"isbn":         "123",
		"price":        10.00,
		"release_year": 2020,
	})

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestGetBook_success(t *testing.T) {
	cleanBooks(t)

	_, created := doJSON(t, http.MethodPost, "/api/v1/books", map[string]any{
		"title":        "Designing Data-Intensive Applications",
		"isbn":         "9781491903100",
		"price":        59.99,
		"release_year": 2017,
	})
	require.NotZero(t, created["id"])
	id := int(created["id"].(float64))

	resp, body := doJSON(t, http.MethodGet, fmt.Sprintf("/api/v1/books/%d", id), nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Designing Data-Intensive Applications", body["title"])
	assert.Equal(t, "9781491903100", body["isbn"])
	assert.Equal(t, 59.99, body["price"])
	assert.Equal(t, float64(2017), body["release_year"])
}

func TestGetBook_notFound(t *testing.T) {
	cleanBooks(t)

	resp, _ := doJSON(t, http.MethodGet, "/api/v1/books/999999", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestListBooks(t *testing.T) {
	cleanBooks(t)

	for _, b := range []map[string]any{
		{"title": "Book A", "isbn": "9780000000001", "price": 10.00, "release_year": 2020},
		{"title": "Book B", "isbn": "9780000000002", "price": 20.00, "release_year": 2021},
	} {
		_, _ = doJSON(t, http.MethodPost, "/api/v1/books", b)
	}

	resp, body := doJSON(t, http.MethodGet, "/api/v1/books", nil)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	list, ok := body["books"].([]any)
	require.True(t, ok)
	assert.Len(t, list, 2)
}

func TestListBooks_filterByISBN(t *testing.T) {
	cleanBooks(t)

	_, _ = doJSON(t, http.MethodPost, "/api/v1/books", map[string]any{
		"title": "Book A", "isbn": "9780000000003", "price": 10.00, "release_year": 2020,
	})
	_, _ = doJSON(t, http.MethodPost, "/api/v1/books", map[string]any{
		"title": "Book B", "isbn": "9780000000004", "price": 20.00, "release_year": 2021,
	})

	resp, body := doJSON(t, http.MethodGet, "/api/v1/books?isbn=9780000000003", nil)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	list, ok := body["books"].([]any)
	require.True(t, ok)
	require.Len(t, list, 1)
	assert.Equal(t, "Book A", list[0].(map[string]any)["title"])
}

// doJSON fires an HTTP request with an optional JSON body and decodes the response.
func doJSON(t *testing.T, method, path string, payload any) (*http.Response, map[string]any) {
	t.Helper()

	var req *http.Request
	if payload != nil {
		b, err := json.Marshal(payload)
		require.NoError(t, err)
		req, _ = http.NewRequest(method, testServer.URL+path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, testServer.URL+path, nil)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	return resp, result
}
