package user

import (
	"net/http"
	"time"

	"github.com/freitasmatheusrn/movie-reservation/internal/auth"
	"github.com/freitasmatheusrn/movie-reservation/pkg/json"
)

type handler struct {
	service  *Service
	tokenGen *auth.TokenGenerator
}

func NewHandler(service *Service, tokenGen *auth.TokenGenerator) *handler {
	return &handler{
		service:  service,
		tokenGen: tokenGen,
	}
}

func (h *handler) Signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := InputDTO{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	output, err := h.service.Signup(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.Send(w, http.StatusCreated, output)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := LoginDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	output, err := h.service.Login(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, ttl, err := h.tokenGen.GenerateAccessToken(output.ID, output.Role)
	if err != nil {
		http.Error(w, "Erro ao emitir token, logue novamente", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(ttl),
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})
	json.Send(w, http.StatusOK, output)
}
