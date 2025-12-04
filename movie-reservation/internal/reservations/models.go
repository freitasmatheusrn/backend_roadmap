package reservations

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusConfirmed ReservationStatus = "confirmed"
	StatusCanceled  ReservationStatus = "canceled"
)

type ReservationTicket struct {
	SeatID       int64   `json:"seat_id"`
	SeatNumber   int     `json:"seat_number"`
	RowLabel     string  `json:"seat_label"`
	SeatFullName string  `json:"seat_fullname"`
	Price        float64 `json:"price"`
}

type Reservation struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	ShowtimeID       int64
	TotalAmount      float64
	Status           ReservationStatus
	BookingReference string
	Tickets          []ReservationTicket
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type InputDTO struct {
	UserID      uuid.UUID           `json:"user_id"`
	ShowtimeID  int64               `json:"showtime_id"`
	TotalAmount float64             `json:"-"`
	Tickets     []ReservationTicket `json:"-"`
}

type OutputDTO struct {
	ID               uuid.UUID           `json:"id"`
	UserID           uuid.UUID           `json:"user_id"`
	ShowtimeID       int64               `json:"showtime_id"`
	TotalAmount      float64             `json:"total_amount"`
	Status           ReservationStatus   `json:"status"`
	BookingReference string              `json:"booking_reference"`
	Theater          string              `json:"theater"`
	Tickets          []ReservationTicket `json:"tickets"`
}
type ReservationSummary struct {
	ID               uuid.UUID           `json:"id"`
	Theater          string              `json:"theater,omitempty"`
	BookingReference string              `json:"booking_reference"`
	Showtime         time.Time           `json:"showtime,omitempty"`
	Tickets          []ReservationTicket `json:"tickets,omitempty"`
}

type ReservationManager struct {
	confirmations map[string]chan bool
	streams       map[string][]chan StatusUpdate
	mu            sync.RWMutex
}

type StatusUpdate struct {
	ReservationID string
	Status        string
	TimeRemaining int
	Message       string
}
