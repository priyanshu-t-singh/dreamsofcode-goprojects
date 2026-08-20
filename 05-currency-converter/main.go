package main

import (
	"fmt"
	"log"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/client"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/converter"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/05-currency-converter/internal/handler"
)

func main() {
	appID, exists := os.LookupEnv("API_TOKEN")
	if !exists {
		log.Fatalln(
			"Must set the env var API_TOKEN value. Vist openexchangerates.org to get one.",
		)
	}

	formData, err := handler.RunCurrencyForm()
	if err != nil {
		log.Fatal(err)
	}

	currencyRates, err := client.FetchExchangeRates(appID)
	if err != nil {
		log.Fatal(err)
	}

	var (
		from   string  = formData.From
		to     string  = formData.To
		amount float64 = formData.Amount
	)

	converterService := converter.NewService()
	convertedAmount, err := converterService.Convert(amount, from, to, currencyRates)
	if err != nil {
		log.Fatalln(err)
	}

	output := handler.RenderOutput(from, to, amount, convertedAmount)
	fmt.Print(output)
}
