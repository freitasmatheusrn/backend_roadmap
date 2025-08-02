package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/freitasmatheusrn/url-shortening-app/internal/dtos"
	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
	"github.com/freitasmatheusrn/url-shortening-app/internal/usecase"
	"github.com/freitasmatheusrn/url-shortening-app/internal/webserver/views"
	"github.com/go-chi/chi/v5"
)

type UrlHandler struct {
	UrlRepository repo.UrlRepositoryInterface
}

func NewUrlHandler(urlRepository *repo.UrlRepository) *UrlHandler {
	return &UrlHandler{
		UrlRepository: urlRepository,
	}
}

func (h *UrlHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	input := dtos.UrlInputDTO{
		OriginalURL: r.FormValue("url"),
	}
	uc := usecase.NewCreateUrlUseCase(h.UrlRepository)
	err := uc.Execute(ctx, input)
	if err != nil {
		http.Error(w, "Error creating short url", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UrlHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	uc := usecase.NewListAllUrlUseCase(h.UrlRepository)
	urlList, err := uc.Execute(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	baseUrl := "http://localhost:8000"
	err = views.Home(urlList, baseUrl).Render(ctx, w)
	if err != nil {
		http.Error(w, "Could not render page", http.StatusInternalServerError)
		return
	}
}
func (h *UrlHandler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	id := chi.URLParam(r, "id")
	err := views.EditForm(id).Render(ctx, w)
	if err != nil {
		http.Error(w, "Could not render page", http.StatusInternalServerError)
		return
	}
}

func (h *UrlHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	input := dtos.UpdateUrlInputDTO{
		ID:          chi.URLParam(r, "id"),
		OriginalURL: r.FormValue("url"),
	}
	uc := usecase.NewUpdateUrlUseCase(h.UrlRepository)
	err := uc.Execute(ctx, input)
	if err != nil {
		http.Error(w, "Could not update", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UrlHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	uc := usecase.NewDeleteUrlUseCase(h.UrlRepository)
	id := chi.URLParam(r, "id")
	err := uc.Execute(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UrlHandler) FowardRequest(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	id := chi.URLParam(r, "shortUrl")
	uc := usecase.NewFowardRequestUseCase(h.UrlRepository)
	originalUrl, err := uc.Execute(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)

}
