package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/godra-y/go-project/pkg/api/models"
)

type CategoryHandler struct {
	CategoryModel *models.CategoryModel
}

func (ch *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := ch.CategoryModel.Insert(&category); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (ch *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	category, err := ch.CategoryModel.Get(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (ch *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := ch.CategoryModel.Update(&category); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (ch *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := ch.CategoryModel.Delete(id); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Category deleted successfully"}
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок Content-Type
	w.WriteHeader(http.StatusOK)                       // Устанавливаем статус код ответа
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal("Failed to encode JSON response:", err)
	}
}
