package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
	"github.com/Damoz1606/ministock-backend/internal/category/response"
	"github.com/Damoz1606/ministock-backend/internal/category/usecase/create"
	"github.com/Damoz1606/ministock-backend/internal/category/usecase/delete"
	"github.com/Damoz1606/ministock-backend/internal/category/usecase/readAll"
	"github.com/Damoz1606/ministock-backend/internal/category/usecase/readOne"
	"github.com/Damoz1606/ministock-backend/internal/category/usecase/update"
)

type Handler struct {
	createHandler  *create.CreateHandler
	readAllHandler *readAll.ReadAllHandler
	readOneHandler *readOne.ReadOneHandler
	updateHandler  *update.UpdateHandler
	deleteHandler  *delete.DeleteHandler
}

func NewHandler(
	createHandler *create.CreateHandler,
	readAllHandler *readAll.ReadAllHandler,
	readOneHandler *readOne.ReadOneHandler,
	updateHandler *update.UpdateHandler,
	deleteHandler *delete.DeleteHandler,
) *Handler {
	return &Handler{
		createHandler:  createHandler,
		readAllHandler: readAllHandler,
		readOneHandler: readOneHandler,
		updateHandler:  updateHandler,
		deleteHandler:  deleteHandler,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	cmd := create.Command{Name: req.Name}
	category, err := h.createHandler.Handle(r.Context(), cmd)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, response.ToCategoryResponse(category))
}

func (h *Handler) ReadAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	name := r.URL.Query().Get("name")
	orderBy := r.URL.Query().Get("orderBy")
	orderDir := r.URL.Query().Get("order")

	q := readAll.Query{
		Page:       page,
		PageSize:   pageSize,
		NameFilter: name,
		OrderBy:    orderBy,
		OrderDir:   orderDir,
	}

	categories, total, err := h.readAllHandler.Handle(r.Context(), q)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response.ToListCategoriesResponse(categories, total, q.Page, q.PageSize))
}

func (h *Handler) ReadOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "categoryId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "invalid category id format")
		return
	}

	q := readOne.Query{CategoryID: id}
	category, err := h.readOneHandler.Handle(r.Context(), q)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response.ToCategoryResponse(category))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "categoryId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "invalid category id format")
		return
	}

	var req updateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	cmd := update.Command{CategoryID: id, Name: req.Name}
	category, err := h.updateHandler.Handle(r.Context(), cmd)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response.ToCategoryResponse(category))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "categoryId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "invalid category id format")
		return
	}

	cmd := delete.Command{CategoryID: id}
	err = h.deleteHandler.Handle(r.Context(), cmd)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
	case errors.Is(err, domain.ErrDuplicateName):
		writeError(w, http.StatusConflict, "DUPLICATE_NAME", "category name already exists")
	case errors.Is(err, domain.ErrInvalidName):
		writeError(w, http.StatusBadRequest, "INVALID_NAME", "invalid category name")
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(newErrorResponse(code, message))
}
