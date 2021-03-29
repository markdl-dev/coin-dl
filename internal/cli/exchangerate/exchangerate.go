package exchangerate

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/leekchan/accounting"
	"github.com/markdl-dev/coin-dl/internal/config"
	"github.com/markdl-dev/go-coin-gecko/coingecko"
	"github.com/pkg/errors"
	"github.com/theckman/yacspin"
)

// const CmdName = "exchange"

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
	var md string
	md += "# Exchange Rates  \n"
	md += "🕔 " + time.Now().Format(time.ANSIC) + " \n"
	md += "## 1 BTC \n"
	md += "| Currency | Value | Type |"
	md += "\n"
	md += "| --- | --- | --- | \n"

	// setup accounting with no symbol(rest api result does not have symbols) and 2 decimal precision
	ac := accounting.Accounting{Symbol: "", Precision: 2}
	if len(e.wantedCurrencies) == 0 {
		for _, val := range e.exchangeRates.Rates {
			bigFloatValue := big.NewFloat(val.Value)
			md += "| " + val.Name + " | " + ac.FormatMoneyBigFloat(bigFloatValue) + " | " + val.Type + " \n"
		}
	}

	for _, currency := range e.wantedCurrencies {
		if val, ok := e.exchangeRates.Rates[currency]; ok {
			bigFloatValue := big.NewFloat(val.Value)
			md += "| " + val.Name + " | " + ac.FormatMoneyBigFloat(bigFloatValue) + " | " + val.Type + " \n"
		}
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(md)
	if err != nil {
		return "", errors.Wrap(err, "exchangerate glamour render")
	}

	return out, nil
}
