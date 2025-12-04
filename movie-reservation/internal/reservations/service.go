package reservations

import (
	"context"
	"errors"
	"strings"

	"github.com/freitasmatheusrn/movie-reservation/pkg/reference"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(ctx context.Context, input InputDTO) (OutputDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]ReservationSummary, error)
	ConfirmReservation(ctx context.Context, reservationID uuid.UUID) error
	Subscribe(id string) (chan StatusUpdate, error)
	Unsubscribe(reservationID string, listener chan StatusUpdate)
}

type Service struct {
	repo RepoInterface
}

func NewService(reservationRepo RepoInterface) Service {
	return Service{
		repo: reservationRepo,
	}
}

func (s *Service) Create(ctx context.Context, input InputDTO) (OutputDTO, error) {
	err := IsValidTickets(input.Tickets)
	if err != nil {
		return OutputDTO{}, err
	}
	bookinReference := reference.Generate()
	r := Reservation{
		UserID:           input.UserID,
		ShowtimeID:       input.ShowtimeID,
		Status:           StatusPending,
		BookingReference: bookinReference,
		Tickets:          input.Tickets,
	}
	output, err := s.repo.Create(ctx, r)
	if err != nil {
		return OutputDTO{}, err
	}
	return output, nil
}

func (s *Service) ConfirmReservation(ctx context.Context, reservationID uuid.UUID) error {
	err := s.repo.ConfirmReservation(ctx, reservationID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Subscribe(id string) (chan StatusUpdate, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id em branco")
	}
	return s.repo.Subscribe(id), nil
}

func (s *Service) Unsubscribe(reservationID string, listener chan StatusUpdate){
	s.repo.Unsubscribe(reservationID, listener)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error{
	err := s.repo.Delete(ctx, id)
	if err != nil{
		return err
	}
	return nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]ReservationSummary, error){
	outputList, err := s.repo.List(ctx, userID)
	if err != nil{
		return nil, err
	}
	return outputList, nil
}