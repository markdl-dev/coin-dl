package market

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
	CmdName                 = "market"
	coinsFlagValue          = "bitcoin"
	coinsFlagDescription    = "Space separated cryptocurrency symbols(ex. bitcoin, etheruem). Refer to coins list command."
	currencyFlagValue       = "usd"
	currencyFlagDescription = "The target currency of the market data(ex. usd). One currency only."
)

type market struct {
	currency    string
	coinsMarket *coingecko.CoinsMarketData
	config      *config.Config
}

func Cmd() error {
	marketFlags := flag.NewFlagSet(CmdName, flag.ExitOnError)

	var config config.Config
	config.Setup(marketFlags)
	config.YakSpinConfig.Suffix = " Getting Market Data from the Gecko"

	spinner, err := yacspin.New(config.YakSpinConfig)
	if err != nil {
		return errors.Wrap(err, "market spinner")
	}

	var coinsStr string
	var currency string

	marketFlags.StringVar(&coinsStr, "coins", coinsFlagValue, coinsFlagDescription)
	marketFlags.StringVar(&coinsStr, "c", coinsFlagValue, coinsFlagDescription+" (shorthand)")
	marketFlags.StringVar(&currency, "currency", currencyFlagValue, currencyFlagDescription)
	marketFlags.StringVar(&currency, "cr", currencyFlagValue, currencyFlagDescription+" (shorthand)")

	if err := marketFlags.Parse(os.Args[2:]); err != nil {
		return errors.Wrap(err, "market flags")
	}
	spinner.Start()

	var coinsSlice []string
	if len(coinsStr) > 0 {
		coinsSlice = strings.Split(coinsStr, " ")
	}

	coinGeckoClient := coingecko.NewClient(nil)
	marketRes, _, err := coinGeckoClient.Coins.GetMarkets(currency, coinsSlice)
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "market result")
	}

	market := market{
		currency:    currency,
		coinsMarket: marketRes,
		config:      &config,
	}

	out, err := market.generateMarketScreen()
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "generate market screen")
	}

	spinner.Stop()
	fmt.Println(out)

	return nil
}

func (m *market) generateMarketScreen() (string, error) {
	var markdownBuilder strings.Builder
	var markdown string

	markdownBuilder.WriteString("# Market Data \n")
	markdown = fmt.Sprintf("🕔 %s \n\n", time.Now().Format(time.ANSIC))
	markdownBuilder.WriteString(markdown)
	markdown = fmt.Sprintf("## Currency: %s \n", strings.ToUpper(m.currency))
	markdownBuilder.WriteString(markdown)
	markdownBuilder.WriteString("| | Name | Price | 1h | 24h | 7d | 14d | 30d | \n")
	markdownBuilder.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- | \n")

	for _, marketData := range *m.coinsMarket {
		currentPrice := coinvalue.New(&marketData.CurrentPrice)
		PCP1H := coinvalue.New(marketData.PriceChangePercentage1HInCurrency)
		PCP24H := coinvalue.New(marketData.PriceChangePercentage24HInCurrency)
		PCP7D := coinvalue.New(marketData.PriceChangePercentage7DInCurrency)
		PCP14D := coinvalue.New(marketData.PriceChangePercentage14DInCurrency)
		PCP30D := coinvalue.New(marketData.PriceChangePercentage30DInCurrency)

		markdown = fmt.Sprintf("| %s | %s | %s | %s %s | %s %s | %s %s | %s %s | %s %s | \n",
			strings.ToUpper(marketData.Symbol),
			strings.ToUpper(marketData.Name),
			currentPrice.Value,
			PCP1H.Emojicon, PCP1H.Value,
			PCP24H.Emojicon, PCP24H.Value,
			PCP7D.Emojicon, PCP7D.Value,
			PCP14D.Emojicon, PCP14D.Value,
			PCP30D.Emojicon, PCP30D.Value)

		markdownBuilder.WriteString(markdown)
	}

	glam, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)

	out, err := glam.Render(markdownBuilder.String())
	if err != nil {
		return "", errors.Wrap(err, "glamour render")
	}

	return out, nil
}
