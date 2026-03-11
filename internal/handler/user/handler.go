package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parseMachineReborn/url_shortener/internal/apperror"
	"github.com/parseMachineReborn/url_shortener/internal/service/user"
)

type handler struct {
	userS     *user.Service
	secretKey string
}

func NewHandler(userS *user.Service, secretKey string) *handler {
	return &handler{
		userS:     userS,
		secretKey: secretKey,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", h.Registration)
	mux.HandleFunc("POST /login", h.LogIn)
}

func (h *handler) Registration(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Не валидный json", http.StatusBadRequest)
		return
	}

	err := h.userS.Register(r.Context(), input.Email, input.Pass)
	if err != nil {
		if errors.Is(err, apperror.ErrAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Проблема при регистрации", http.StatusInternalServerError)
		return
	}
}

func (h *handler) LogIn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Не валидный json", http.StatusBadRequest)
		return
	}

	userId, err := h.userS.LogIn(r.Context(), input.Email, input.Pass)
	if err != nil {
		http.Error(w, "Не правильный email или пароль", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		http.Error(w, "Проблема создания токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"token": tokenStr}); err != nil {
		http.Error(w, "Ошибка маршаллинга токена", http.StatusUnauthorized)
		return
	}
}
