package showtimes

import (
	"fmt"
	"time"
)

func IsValidTimeRange(startDate, endDate time.Time) error {
	if startDate.After(endDate){
		return fmt.Errorf("início da sessão não pode ser após o fim da sessão")
	}
	return nil
}

