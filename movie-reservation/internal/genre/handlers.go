package genre

import (
	"fmt"
	"net/http"

	"github.com/freitasmatheusrn/movie-reservation/pkg/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Send(w, http.StatusCreated, output)
}
func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	outputList, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.Send(w, http.StatusOK, outputList)
}
