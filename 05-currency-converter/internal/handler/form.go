package handler

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"charm.land/huh/v2"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/converter"
)

type ConversionForm struct {
	Amount float64
	From   string
	To     string
}

func RunCurrencyForm() (*ConversionForm, error) {
	currencyOptions := []huh.Option[string]{}
	for country, symbol := range converter.CurrencySymbols {
		currencyOptions = append(currencyOptions, huh.Option[string]{
			Key:   symbol + " " + country,
			Value: country,
		})
	}

	var currencyStr string

	formData := &ConversionForm{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your base currency").
				Options(currencyOptions...).
				Value(&formData.From),

			huh.NewSelect[string]().
				Title("What do you want to convert into").
				Options(currencyOptions...).
				Value(&formData.To),

			huh.NewInput().
				Title("How much to convert?").
				Value(&currencyStr).
				Validate(func(s string) error {
					num, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return errors.New("Invalid format")
					}
					formData.Amount = num
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(os.Stderr, "User exited")
			os.Exit(1)
		}
		return nil, err
	}

	return formData, nil
}

func RenderOutput(from, to string, originalAmount, convertedAmount float64) string {
	fromSym := converter.CurrencySymbols[from]
	toSym := converter.CurrencySymbols[to]

	return fmt.Sprintf("%s%s %s converts to %s%s %s\n",
		fromSym, converter.FormatAmount(originalAmount), from,
		toSym, converter.FormatAmount(convertedAmount), to,
	)
}
