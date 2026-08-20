package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/models"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchExchangeRates(appID string) (*models.ExchangeRates, error) {
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

	var currencyRates *models.ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&currencyRates); err != nil {
		return nil, err
	}

	return currencyRates, nil
}
