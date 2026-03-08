package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/parseMachineReborn/url_shortener/internal/apperror"
	"github.com/parseMachineReborn/url_shortener/internal/service"
)

type handler struct {
	s *service.URLService
}

func NewHandler(s *service.URLService) *handler {
	return &handler{
		s: s,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /geturl/{shortURL}", h.GetURL)
	mux.HandleFunc("GET /getAll", h.GetAll)
	mux.HandleFunc("DELETE /delete/{shortURL}", h.Delete)
}

func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Не валидный json в body", http.StatusBadRequest)
		return
	}

	res := h.s.Shorten(r.Context(), input.Address)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Проблема перевода в JSON укороченной ссылки", http.StatusInternalServerError)
	}
}

func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) {
	url := r.PathValue("shortURL")

	if _, err := h.s.GetURL(r.Context(), url); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			http.Error(w, "URL не найден", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	storage, err := h.s.GetAll(r.Context())

	if err != nil {
		http.Error(w, "Ошибка при получении списка сохраненных укороченных URL", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(storage); err != nil {
		http.Error(w, "Ошибка при маршаллинге списка сохраненных укороченных URL", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("shortURL")

	if err := h.s.Delete(r.Context(), shortURL); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			http.Error(w, "Не найдено элемента с таким ключом(shortURL)", http.StatusNotFound)
			return
		}

		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
	}
}
