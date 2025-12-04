package showtimes

import (
	"net/http"
	"strconv"

	"github.com/freitasmatheusrn/movie-reservation/pkg/conversor"
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Não foi possível processar os dados", http.StatusBadRequest)
		return
	}
	movieID, err := strconv.Atoi(r.FormValue("movie_id"))
	if err != nil {
		http.Error(w, "Filme inválido", http.StatusBadRequest)
		return
	}
	theaterID, err := strconv.Atoi(r.FormValue("theater_id"))
	if err != nil {
		http.Error(w, "Sala inválida", http.StatusBadRequest)
		return
	}
	startDate, err := conversor.StrToDateTime(r.FormValue("start_date"), "dateTimeLayout")
	if err != nil {
		http.Error(w, "Horário de início inválido", http.StatusBadRequest)
		return
	}
	endDate, err := conversor.StrToDateTime(r.FormValue("end_date"), "dateTimeLayout")
	if err != nil {
		http.Error(w, "Horário de fim de sessão inválido", http.StatusBadRequest)
		return
	}
	basePrice, err := strconv.ParseFloat(r.FormValue("base_price"), 64)
	if err != nil {
		http.Error(w, "Horário de fim de sessão inválido", http.StatusBadRequest)
		return
	}
	i := InputDTO{
		MovieID:   int64(movieID),
		TheaterID: int64(theaterID),
		StartTime: startDate,
		EndTime:   endDate,
		BasePrice: basePrice,
	}

	output, err := h.service.Create(r.Context(), i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.Send(w, http.StatusCreated, output)

}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	strDay := r.URL.Query().Get("date")
	day, err := conversor.StrToDateTime(strDay, "dateLayout")
	if err != nil {
		http.Error(w, "Data inválida", http.StatusBadRequest)
		return
	}
	outputList, err := h.service.List(r.Context(), day)
	if err != nil {
		http.Error(w, "Data inválida", http.StatusInternalServerError)
		return
	}
	json.Send(w, http.StatusOK, outputList)
}
