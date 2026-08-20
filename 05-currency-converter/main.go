package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/huh/v2"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/models"
)

var CurrencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"INR": "₹",
}

func main() {
	appID, exists := os.LookupEnv("API_TOKEN")
	if !exists {
		log.Fatalln(
			"Must set the env var API_TOKEN value. Vist openexchangerates.org to get one.",
		)
	}

	var (
		from        string
		to          string
		currencyStr string
		amount      float64
	)

	currencyOptions := []huh.Option[string]{}
	for country, symbol := range CurrencySymbols {
		currencyOptions = append(currencyOptions, huh.Option[string]{
			Key:   symbol + " " + country,
			Value: country,
		})
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your base currency").
				Options(currencyOptions...).
				Value(&from),

			huh.NewSelect[string]().
				Title("What do you want to convert into").
				Options(currencyOptions...).
				Value(&to),

			huh.NewInput().
				Title("How much to convert?").
				Value(&currencyStr).
				Validate(func(s string) error {
					num, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return errors.New("Invalid format")
					}
					amount = num
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
		log.Fatal(err)
	}

	currencyRates, err := getCurrencyRates(appID)
	if err != nil {
		log.Fatal(err)
	}

	rateFrom, exists := currencyRates.Rates[from]
	if !exists {
		log.Fatalf("%s value doesn't exists", from)
	}
	rateTo, exists := currencyRates.Rates[to]
	if !exists {
		log.Fatalf("%s value doesn't exists", to)
	}

	var convertedAmount = (amount / rateFrom) * rateTo
	fmt.Printf(
		"%s converts to %s\n",
		CurrencySymbols[from]+formatAmount(amount)+" "+from,
		CurrencySymbols[to]+formatAmount(convertedAmount)+" "+to,
	)
}

func formatAmount(val float64) string {
	s := strconv.FormatFloat(val, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func getCurrencyRates(appID string) (*models.OpenExchangeRatesAPIResponse, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", appID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v\n", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK http status code: %d\n", resp.StatusCode)
	}

	var currencyRates *models.OpenExchangeRatesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&currencyRates); err != nil {
		return nil, err
	}

	return currencyRates, nil
}
