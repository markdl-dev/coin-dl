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
)

const defaultCurrency = "usd"

func Cmd() error {
	xrFlags := flag.NewFlagSet("xr", flag.ContinueOnError)
	// TODO add usage
	var config config.Config
	xrFlags.BoolVar(&config.ShowNotifications, "showNotifs", true, "show notification pop ups")
	xrFlags.BoolVar(&config.PlayNotificationsBeep, "playNotifs", true, "plays a beep when a notification shows")
	convertToCurrency := xrFlags.String("to", "usd", "Get BTC value in specified currency")
	if err := xrFlags.Parse(os.Args[2:]); err != nil {
		return errors.Wrap(err, "xr flags")
	}

	convertToCurrencies := strings.Split(*convertToCurrency, " ")
	if len(convertToCurrencies) == 0 {
		convertToCurrencies = append(convertToCurrencies, defaultCurrency)
	}

	coinGeckoClient := coingecko.NewClient(nil)
	xrList, _, err := coinGeckoClient.ExchangeRate.GetExchangeRates()
	if err != nil {
		return errors.Wrap(err, "xr list")
	}

	out, err := generateExchangeRateScreen(xrList, convertToCurrencies)
	if err != nil {
		return errors.Wrap(err, "xr screen")
	}

	fmt.Print(out)

	return nil
}

func generateExchangeRateScreen(exchangeRates *coingecko.ExchangeRates, currencies []string) (string, error) {
	var md string
	md += "# Exchange Rates  \n"
	md += time.Now().Format(time.ANSIC) + " \n"
	md += "## 1 BTC \n"
	md += "| Currency | Value | Type |"
	md += "\n"
	md += "| --- | --- | --- | \n"

	for _, currency := range currencies {
		if val, ok := exchangeRates.Rates[currency]; ok {
			ac := accounting.Accounting{Symbol: "", Precision: 2}
			bigFloatValue := big.NewFloat(val.Value)
			md += "| " + val.Name + " | " + ac.FormatMoneyBigFloat(bigFloatValue) + " | " + val.Type + " \n"
			// fmt.Println(val)
		}
	}

	out, err := glamour.Render(md, "dark")
	if err != nil {
		return "", errors.Wrap(err, "glamour render")
	}

	return out, nil
}
