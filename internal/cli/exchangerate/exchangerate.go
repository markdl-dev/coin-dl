package exchangerate

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/markdl-dev/coin-dl/internal/coinvalue"
	"github.com/markdl-dev/coin-dl/internal/config"
	"github.com/markdl-dev/go-coin-gecko/coingecko"
	"github.com/pkg/errors"
	"github.com/theckman/yacspin"
)

const (
	CmdName                 = "exchange"
	currencyFlagDescription = "Get BTC value in specified currency/currencies"
)

type exchangeRate struct {
	exchangeRates    *coingecko.ExchangeRates
	wantedCurrencies []string
	config           *config.Config
}

func Cmd() error {
	exchangeRateFlag := flag.NewFlagSet(CmdName, flag.ExitOnError)

	var config config.Config
	config.Setup(exchangeRateFlag)
	config.YakSpinConfig.Suffix = " Getting Exchange Rates from the Gecko"

	spinner, err := yacspin.New(config.YakSpinConfig)
	if err != nil {
		return errors.Wrap(err, "exchangerate spinner")
	}

	var toCurrenciesStr string
	exchangeRateFlag.StringVar(&toCurrenciesStr, "currency", "", currencyFlagDescription)
	exchangeRateFlag.StringVar(&toCurrenciesStr, "cr", "", currencyFlagDescription+" (shorthand)")

	if err := exchangeRateFlag.Parse(os.Args[2:]); err != nil {
		return errors.Wrap(err, "exchangerate flags")
	}

	spinner.Start()

	var toCurrenciesSlice []string
	if len(toCurrenciesStr) > 0 {
		toCurrenciesSlice = strings.Split(toCurrenciesStr, " ")
	}

	coinGeckoClient := coingecko.NewClient(nil)
	xrList, _, err := coinGeckoClient.ExchangeRate.GetExchangeRates()
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "exchangerate list")
	}

	xr := exchangeRate{exchangeRates: xrList,
		wantedCurrencies: toCurrenciesSlice,
		config:           &config}

	out, err := xr.generateExchangeRateScreen()
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "exchangerate screen")
	}

	spinner.Stop()
	fmt.Print(out)

	return nil
}

func (e *exchangeRate) generateExchangeRateScreen() (string, error) {
	var markdownBuilder strings.Builder
	var markdown string
	markdownBuilder.WriteString("# Exchange Rates  \n")
	markdown = fmt.Sprintf("🕔 %s \n", time.Now().Format(time.ANSIC))
	markdownBuilder.WriteString(markdown)
	markdownBuilder.WriteString("## 1 BTC \n")
	markdownBuilder.WriteString("| Currency | Value | Type | \n")
	markdownBuilder.WriteString("| --- | --- | --- | \n")

	// setup accounting with no symbol(rest api result does not have symbols) and 2 decimal precision
	if len(e.wantedCurrencies) == 0 {
		for _, val := range e.exchangeRates.Rates {
			exchangeRate := coinvalue.New(&val.Value)
			markdown = fmt.Sprintf("| %s | %s | %s | \n",
				val.Name,
				exchangeRate.Value,
				val.Type)
			markdownBuilder.WriteString(markdown)
		}
	}

	for _, currency := range e.wantedCurrencies {
		if val, ok := e.exchangeRates.Rates[currency]; ok {
			exchangeRate := coinvalue.New(&val.Value)
			markdown = fmt.Sprintf("| %s | %s | %s | \n",
				val.Name,
				exchangeRate.Value,
				val.Type)
			markdownBuilder.WriteString(markdown)
		}
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(markdownBuilder.String())
	if err != nil {
		return "", errors.Wrap(err, "exchangerate glamour render")
	}

	return out, nil
}
