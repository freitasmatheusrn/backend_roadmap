package reservations

import "fmt"

func IsValidTickets(tickets []ReservationTicket)error{
	if len(tickets) <= 0{
		return fmt.Errorf("nenhum assento encontrado")
	}
	return nil
}

