package movie

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/freitasmatheusrn/movie-reservation/pkg/json"
	"github.com/go-chi/chi"
)

type handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input InputDTO
	if err := json.Read(r, &input); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao ler JSON: %s", err), http.StatusBadRequest)
		return
	}

	output, err := h.service.Create(r.Context(), input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar filme: %s", err), http.StatusBadRequest)
		return
	}

	json.Send(w, http.StatusCreated, output)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "estamos com problema de servidor, por favor tente novamente mais tarde", http.StatusInternalServerError)
		return
	}
	json.Send(w, http.StatusCreated, output)

}

func (h *handler) ListByGenre(w http.ResponseWriter, r *http.Request) {
	genres := r.URL.Query()["genres"]
	log.Printf("%v", genres)
	output, err := h.service.ListByGenre(r.Context(), genres)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Send(w, http.StatusOK, output)
}

func (h *handler) ListByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	outputList, err := h.service.ListByName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Send(w, http.StatusOK, outputList)
}

func (h *handler) Show(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
	}
	output, err := h.service.Show(r.Context(), id)
	if err != nil {
		http.Error(w, "id não encontrado", http.StatusNotFound)
	}
	json.Send(w, http.StatusCreated, output)
}
