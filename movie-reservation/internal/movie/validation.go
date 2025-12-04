package movie

import (
	"fmt"
	"net/url"
	"time"
)

func NameValid(name string) error {
	if len(name) < 5 {
		return fmt.Errorf("nome muito curto, precisa de no mínimo 5 letras")
	}
	return nil
}

func DescriptionValid(desc string) error {
	if len(desc) < 20 {
		return fmt.Errorf("descrição muito curta, precisa de pelo menos 20 caracteres")
	}
	return nil
}

func IsPosterValid(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("url inválida")
	}
	return nil
}

func IsValidDuration(duration int) error {
	if duration < 0 {
		return fmt.Errorf("duração inválida")
	}
	return nil
}

func IsValidReleaseDate(date time.Time) error {
	today := time.Now().Truncate(24 * time.Hour)

	date = date.Truncate(24 * time.Hour)
	if date.Before(today) {
		return fmt.Errorf("data de lançamento não pode ser no passado")
	}
	return nil
}

func IsValidRating(rating string) error {
	if len(rating) > 150 {
		return fmt.Errorf("nota não pode exceder 150 caracteres")
	}
	return nil
}
