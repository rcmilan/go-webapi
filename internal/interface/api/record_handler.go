package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type RecordHandler struct{}

func NewRecordHandler() *RecordHandler {
	return &RecordHandler{}
}

// ── DTOs ─────────────────────────────────────────────────────────────────────

type recordDTO struct {
	ID     uint32  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Year   uint16  `json:"year"`
	Genre  string  `json:"genre"`
	Price  float64 `json:"price"`
}

type createRecordInput struct {
	Body struct {
		Title  string  `json:"title" doc:"Record title" minLength:"1"`
		Artist string  `json:"artist" doc:"Artist name" minLength:"1"`
		Year   uint16  `json:"year" doc:"Release year"`
		Genre  string  `json:"genre" doc:"Music genre"`
		Price  float64 `json:"price" doc:"Record price" minimum:"0.01"`
	}
}

type createRecordOutput struct{ Body recordDTO }

type getRecordInput struct {
	ID uint32 `path:"id" doc:"Record ID"`
}

type getRecordOutput struct{ Body recordDTO }

type listRecordsInput struct {
	Artist string `query:"artist" doc:"Filter by artist"`
	Genre  string `query:"genre" doc:"Filter by genre"`
}

type listRecordsOutput struct {
	Body struct {
		Records []recordDTO `json:"records"`
	}
}

// ── Handlers ─────────────────────────────────────────────────────────────────

func (h *RecordHandler) createRecord(_ context.Context, _ *createRecordInput) (*createRecordOutput, error) {
	return &createRecordOutput{}, nil
}

func (h *RecordHandler) getRecord(_ context.Context, _ *getRecordInput) (*getRecordOutput, error) {
	return &getRecordOutput{}, nil
}

func (h *RecordHandler) listRecords(_ context.Context, _ *listRecordsInput) (*listRecordsOutput, error) {
	out := &listRecordsOutput{}
	out.Body.Records = []recordDTO{}
	return out, nil
}

// ── Routing ──────────────────────────────────────────────────────────────────

func (h *RecordHandler) Register(api huma.API) {
	huma.Register(api, op("create-record", http.MethodPost, "/api/v1/records", "Register a new record", "Records"), h.createRecord)
	huma.Register(api, op("get-record", http.MethodGet, "/api/v1/records/{id}", "Get a record by ID", "Records"), h.getRecord)
	huma.Register(api, op("list-records", http.MethodGet, "/api/v1/records", "List records", "Records"), h.listRecords)
}
