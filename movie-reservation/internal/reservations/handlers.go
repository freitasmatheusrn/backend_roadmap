package reservations

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	cjson "github.com/freitasmatheusrn/movie-reservation/pkg/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
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
	_, claims, _ := jwtauth.FromContext(r.Context())
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "token inválido", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "id no formato inválido", http.StatusUnauthorized)
		return
	}

	showtimeID, err := strconv.ParseInt(chi.URLParam(r, "showtime_id"), 10, 64)
	if err != nil {
		http.Error(w, "id inválido na url", http.StatusBadRequest)
		return
	}
	body := r.Body
	defer body.Close()

	seats := r.Form["seat_ids"]
	var ids []int64

	for _, s := range seats {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			http.Error(w, "um assento inválido foi enviado", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	var tickets []ReservationTicket

	for _, id := range ids {
		tickets = append(tickets, ReservationTicket{SeatID: id})
	}

	input := InputDTO{
		UserID:     userID,
		ShowtimeID: showtimeID,
		Tickets:    tickets,
	}
	output, err := h.service.Create(r.Context(), input)

	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao criar reserva, %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "SSE não suportado", http.StatusInternalServerError)
		return
	}

	initial, _ := json.Marshal(output)
	fmt.Fprintf(w, "data: %s\n\n", initial)
	flusher.Flush()
	listener, err := h.service.Subscribe(output.ID.String())

	if err != nil {
		log.Println("erro ao inscrever SSE:", err)
		return
	}

	defer h.service.Unsubscribe(output.ID.String(), listener)
	log.Printf("SSE iniciado para reserva %s", output.ID)

	for {
		select {
		case update, ok := <-listener:
			if !ok {
				log.Printf("Canal SSE fechado para %s", output.ID)
				return
			}

			data, err := json.Marshal(update)
			if err != nil {
				log.Printf("Erro ao serializar update: %v", err)
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			log.Printf("📤 SSE enviado: %s - %s", output.ID, update.Status)

			if update.Status == "confirmed" || update.Status == "cancelled" {
				log.Printf("Finalizando SSE para %s", output.ID)
				return
			}

		case <-r.Context().Done():
			log.Printf("Cliente desconectou SSE %s", output.ID)
			return
		}
	}
}

func (h *handler) ConfirmReservation(w http.ResponseWriter, r *http.Request) {
	reservationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id inválido", http.StatusUnprocessableEntity)
		return
	}
	err = h.service.ConfirmReservation(r.Context(), reservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cjson.Send(w, http.StatusOK, "Confirmado com sucesso")
}

func (h *handler) PaymentPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE não suportado", http.StatusInternalServerError)
		return
	}
	listener, err := h.service.Subscribe(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer h.service.Unsubscribe(id, listener)
	for {
		select {
		case update, ok := <-listener:
			if !ok {
				log.Printf("📡 Canal fechado para %s", id)
				return
			}

			data, err := json.Marshal(update)
			if err != nil {
				log.Printf("Erro ao serializar update: %v", err)
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			log.Printf("📤 Enviado update SSE: %s - %s", id, update.Status)

			if update.Status == "confirmed" || update.Status == "cancelled" {
				log.Printf("📡 Status final alcançado, encerrando stream para %s", id)
				return
			}

		case <-r.Context().Done():
			log.Printf("📡 Cliente desconectou de %s", id)
			return
		}
	}

}


func (h *handler) List(w http.ResponseWriter, r *http.Request){
	_, claims, _ := jwtauth.FromContext(r.Context())
	userIDStr, ok := claims["sub"].(string)
	if !ok{
		http.Error(w, "token inválido", http.StatusUnauthorized)
		return

	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "id no formato inválido", http.StatusUnauthorized)
		return
	}
	list, err := h.service.List(r.Context(), userID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	cjson.Send(w, http.StatusOK, list)
}