package converter

import (
	"errors"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Convert(amount float64, from string, to string, exchangeRates *models.ExchangeRates) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}
	if from == to {
		return amount, nil
	}

	rates := exchangeRates.Rates
	rateFrom, okFrom := rates[from]
	rateTo, okTo := rates[to]

	if !okFrom || !okTo {
		return 0, errors.New("unsupported currency code")
	}

	return (amount / rateFrom) * rateTo, nil
}
