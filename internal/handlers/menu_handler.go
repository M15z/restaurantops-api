package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/m15z/restaurantops-api/internal/models"
	"github.com/m15z/restaurantops-api/internal/services"
)

type MenuHandler struct {
	menu *services.MenuService
}

func NewMenuHandler(menu *services.MenuService) *MenuHandler {
	return &MenuHandler{menu: menu}
}

// --- Categories ---

type createCategoryRequest struct {
	Name string `json:"name"`
}

func (h *MenuHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.menu.GetCategories(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch categories"})
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

func (h *MenuHandler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	id, err := h.menu.CreateCategory(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create category"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// --- Menu items ---
type menuItemRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	CategoryID  int64   `json:"category_id"`
	IsAvailable bool    `json:"is_available"`
}

func (h *MenuHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.menu.GetItems(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch menu items"})
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *MenuHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req menuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	item := models.MenuItem{
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
		Category:    models.Category{ID: req.CategoryID},
	}

	id, err := h.menu.CreateItem(r.Context(), item)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, positive price, and category_id required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create item"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *MenuHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var req menuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	item := models.MenuItem{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
		Category:    models.Category{ID: req.CategoryID},
	}

	if err := h.menu.UpdateItem(r.Context(), item); err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, positive price, and category_id required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not update item"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "Updated"})
}

type availabitlyRequest struct {
	IsAvailable bool `json:"is_available"`
}

func (h *MenuHandler) SetAvailability(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt("id", 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req availabitlyRequest
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if err := h.menu.SetItemAvailability(r.Context(), id, req.IsAvailable); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "could not update availability"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "Updated"})
}
