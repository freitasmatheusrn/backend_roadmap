package reservations

import (
	"context"
	"fmt"
	"log"
	"time"

	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RepoInterface interface {
	Create(ctx context.Context, reservation Reservation) (OutputDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]ReservationSummary, error)
	ConfirmReservation(ctx context.Context, reservationID uuid.UUID) error
	startReservationTimeout(reservationID string, confirmChan chan bool)
	CancelExpiredReservation(ctx context.Context, reservationID string) error
	Subscribe(reservationID string) chan StatusUpdate
	Unsubscribe(reservationID string, listener chan StatusUpdate)
	NotifyStatusChange(update StatusUpdate)
}

type Repository struct {
	DB      *pgx.Conn
	manager *ReservationManager
}

func NewReservationManager() *ReservationManager {
	return &ReservationManager{
		confirmations: make(map[string]chan bool),
		streams:       make(map[string][]chan StatusUpdate),
	}
}

func NewRepo(db *pgx.Conn) *Repository {
	manager := NewReservationManager()
	return &Repository{
		DB:      db,
		manager: manager,
	}
}

func (r *Repository) Create(ctx context.Context, rs Reservation) (OutputDTO, error) {
	var output OutputDTO

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return output, err
	}
	defer tx.Rollback(ctx)
	var (
		basePrice     float64
		theater       string
		reservationID uuid.UUID
	)
	err = tx.QueryRow(ctx,
		`SELECT s.base_price, t.name
		FROM showtimes AS s
		JOIN theaters AS t ON t.id = s.theater_id
		WHERE s.id = $1`,
		rs.ShowtimeID,
	).Scan(&basePrice, &theater)

	if err != nil {
		return output, fmt.Errorf("erro ao buscar showtime: %w", err)
	}

	totalAmount := float64(len(rs.Tickets)) * basePrice

	err = tx.QueryRow(ctx,
		`INSERT INTO reservations (user_id, showtime_id, total_amount, status, booking_reference)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id`,
		rs.UserID,
		rs.ShowtimeID,
		totalAmount,
		rs.Status,
		rs.BookingReference,
	).Scan(&reservationID)

	if err != nil {
		return output, fmt.Errorf("erro ao criar reserva: %w", err)
	}

	var tickets []ReservationTicket

	for _, t := range rs.Tickets {
		var ticket ReservationTicket

		err = tx.QueryRow(ctx,
			`INSERT INTO tickets (reservation_id, showtime_id, seat_id, price, reservation_status)
             VALUES ($1, $2, $3, $4, 'pending')
			 RETURNING seat_id, price`,
			reservationID,
			rs.ShowtimeID,
			t.SeatID,
			basePrice,
		).Scan(&ticket.SeatID, &ticket.Price)

		if err != nil {
			return OutputDTO{}, fmt.Errorf("assento %d indisponível: %w", t.SeatID, err)
		}

		tickets = append(tickets, ticket)
	}

	if err := tx.Commit(ctx); err != nil {
		return OutputDTO{}, fmt.Errorf("erro ao commitar transação: %w", err)
	}
	output = OutputDTO{
		ID:               reservationID,
		TotalAmount:      totalAmount,
		Theater:          theater,
		BookingReference: rs.BookingReference,
		ShowtimeID:       rs.ShowtimeID,
		UserID:           rs.UserID,
		Status:           rs.Status,
		Tickets:          tickets,
	}

	confirmChan := make(chan bool, 1)

	r.manager.mu.Lock()
	r.manager.confirmations[reservationID.String()] = confirmChan
	r.manager.mu.Unlock()

	go r.startReservationTimeout(reservationID.String(), confirmChan)

	return output, nil
}

func (r *Repository) startReservationTimeout(reservationID string, confirmChan chan bool) {
	ctx := context.Background()

	log.Printf("⏰ Iniciando timeout de 3 minutos para reserva %s", reservationID)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	duration := 3 * time.Minute

	for {
		select {
		case confirmed := <-confirmChan:
			if confirmed {
				log.Printf("✅ Reserva %s confirmada antes do timeout", reservationID)
				if r.manager != nil {
					r.NotifyStatusChange(StatusUpdate{
						ReservationID: reservationID,
						Status:        "confirmed",
						TimeRemaining: 0,
						Message:       "Reserva confirmada com sucesso!",
					})
				}
			}

			if r.manager != nil {
				r.manager.mu.Lock()
				delete(r.manager.confirmations, reservationID)
				r.manager.mu.Unlock()
			}
			return

		case <-ticker.C:
			elapsed := time.Since(startTime)
			remaining := duration - elapsed

			if remaining <= 0 {
				// Timeout expirou
				log.Printf("⏱️ Timeout expirado para reserva %s - cancelando", reservationID)

				if err := r.CancelExpiredReservation(ctx, reservationID); err != nil {
					log.Printf("❌ Erro ao cancelar reserva: %v", err)
				}

				// Notificar listeners
				if r.manager != nil {
					r.NotifyStatusChange(StatusUpdate{
						ReservationID: reservationID,
						Status:        "cancelled",
						TimeRemaining: 0,
						Message:       "Tempo expirado. Reserva cancelada.",
					})

					r.manager.mu.Lock()
					delete(r.manager.confirmations, reservationID)
					r.manager.mu.Unlock()
				}
				return
			}

			log.Println("notificando o tempo")
			r.NotifyStatusChange(StatusUpdate{
				ReservationID: reservationID,
				Status:        "pending",
				TimeRemaining: int(remaining.Seconds()),
				Message:       fmt.Sprintf("Tempo restante: %d segundos", int(remaining.Seconds())),
			})

		}
	}
}

func (r *Repository) CancelExpiredReservation(ctx context.Context, reservationID string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentStatus string
	err = tx.QueryRow(ctx,
		`UPDATE reservations 
		 SET status = 'cancelled', updated_at = NOW()
		 WHERE id = $1 AND status = 'pending'
		 RETURNING status`,
		reservationID,
	).Scan(&currentStatus)

	if err != nil {
		return fmt.Errorf("erro ao cancelar reserva: %w", err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM tickets WHERE reservation_id = $1`,
		reservationID,
	)
	if err != nil {
		return fmt.Errorf("erro ao deletar tickets: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao commitar: %w", err)
	}

	return nil
}

func (r *Repository) ConfirmReservation(ctx context.Context, reservationID uuid.UUID) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	var status string
	var createdAt time.Time
	err = tx.QueryRow(ctx,
		`SELECT status, created_at 
		 FROM reservations 
		 WHERE id = $1`,
		reservationID,
	).Scan(&status, &createdAt)

	if err != nil {
		return fmt.Errorf("reserva não encontrada: %w", err)
	}

	if time.Since(createdAt) > 3*time.Minute {
		return fmt.Errorf("reserva expirou, não pode ser confirmada")
	}

	if status == "cancelled" {
		return fmt.Errorf("reserva já foi cancelada")
	}

	if status == "confirmed" {
		return fmt.Errorf("reserva já está confirmada")
	}

	_, err = tx.Exec(ctx,
		`UPDATE reservations 
		 SET status = 'confirmed', updated_at = NOW()
		 WHERE id = $1`,
		reservationID,
	)
	if err != nil {
		return fmt.Errorf("erro ao confirmar reserva: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE tickets 
		 SET reservation_status = 'confirmed'
		 WHERE reservation_id = $1`,
		reservationID,
	)
	if err != nil {
		return fmt.Errorf("erro ao confirmar tickets: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao commitar: %w", err)
	}
	if r.manager != nil {
		r.manager.mu.RLock()
		if confirmChan, exists := r.manager.confirmations[reservationID.String()]; exists {
			select {
			case confirmChan <- true:
				log.Printf("📢 Notificação de confirmação enviada para reserva %s", reservationID)
			default:
				log.Printf("⚠️ Não foi possível notificar canal para reserva %s", reservationID.String())
			}
		}
		r.manager.mu.RUnlock()
	}

	log.Printf("✅ Reserva %s confirmada com sucesso", reservationID)
	return nil
}

func (r *Repository) Subscribe(reservationID string) chan StatusUpdate {
	r.manager.mu.Lock()
	defer r.manager.mu.Unlock()

	listener := make(chan StatusUpdate, 5)
	r.manager.streams[reservationID] = append(r.manager.streams[reservationID], listener)

	log.Printf("📡 Novo listener inscrito para reserva %s (total: %d)",
		reservationID, len(r.manager.streams[reservationID]))

	return listener
}

func (r *Repository) Unsubscribe(reservationID string, listener chan StatusUpdate) {
	r.manager.mu.Lock()
	defer r.manager.mu.Unlock()

	listeners := r.manager.streams[reservationID]
	for i, l := range listeners {
		if l == listener {
			r.manager.streams[reservationID] = slices.Delete(listeners, i, i+1)
			close(listener)
			log.Printf("📡 Listener removido para reserva %s (restantes: %d)",
				reservationID, len(r.manager.streams[reservationID]))
			break
		}
	}

	if len(r.manager.streams[reservationID]) == 0 {
		delete(r.manager.streams, reservationID)
	}
}

func (r *Repository) NotifyStatusChange(update StatusUpdate) {
	r.manager.mu.RLock()
	listeners := r.manager.streams[update.ReservationID]
	r.manager.mu.RUnlock()

	log.Printf("📢 Notificando %d listeners sobre mudança de status para %s: %s",
		len(listeners), update.ReservationID, update.Status)

	for _, listener := range listeners {
		select {
		case listener <- update:
		default:
			log.Printf("⚠️ Não foi possível notificar listener (canal cheio/fechado)")
		}
	}
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	cmdTag, err := r.DB.Exec(ctx,
		`DELETE FROM reservations WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reservation %s not found", id)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID) ([]ReservationSummary, error) {
    query := `
        SELECT 
            r.id,
            th.name as theater,
            st.start_time as showtime,
            t.seat_id,
            s.seat_number,
            s.row_label,
            CONCAT(s.row_label, s.seat_number) as seat_fullname,
            t.price
        FROM reservations r
        INNER JOIN showtimes st ON r.showtime_id = st.id
        INNER JOIN theaters th ON st.theater_id = th.id
        INNER JOIN tickets t ON t.reservation_id = r.id
        INNER JOIN seats s ON t.seat_id = s.id
        WHERE r.user_id = $1
        ORDER BY r.created_at DESC, t.id
    `
    
    rows, err := r.DB.Query(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query reservations: %w", err)
    }
    defer rows.Close()
    
    // Mapa para agrupar tickets por reserva
    reservationsMap := make(map[uuid.UUID]*ReservationSummary)
    
    for rows.Next() {
        var (
            reservationID uuid.UUID
            theater       string
            showtime      time.Time
            ticket        ReservationTicket
        )
        
        err := rows.Scan(
            &reservationID,
            &theater,
            &showtime,
            &ticket.SeatID,
            &ticket.SeatNumber,
            &ticket.RowLabel,
            &ticket.SeatFullName,
            &ticket.Price,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        
        // Se a reserva já existe no mapa, adiciona o ticket
        if res, exists := reservationsMap[reservationID]; exists {
            res.Tickets = append(res.Tickets, ticket)
        } else {
            // Cria nova reserva
            reservationsMap[reservationID] = &ReservationSummary{
                ID:       reservationID,
                Theater:  theater,
                Showtime: showtime,
                Tickets:  []ReservationTicket{ticket},
            }
        }
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    
    // Converte o mapa para slice
    result := make([]ReservationSummary, 0, len(reservationsMap))
    for _, reservation := range reservationsMap {
        result = append(result, *reservation)
    }
    
    return result, nil
}