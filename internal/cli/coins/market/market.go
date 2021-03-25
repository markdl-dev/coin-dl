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

type market struct {
	currency    string
	coinsMarket *coingecko.CoinsMarket
	config      *config.Config
}

func Cmd() error {
	marketFlags := flag.NewFlagSet("market", flag.ExitOnError)

	var config config.Config
	config.Setup(marketFlags)
	config.YakSpinConfig.Suffix = " Getting Market Data from the Gecko"

	spinner, err := yacspin.New(config.YakSpinConfig)
	if err != nil {
		return errors.Wrap(err, "market spinner")
	}

	coinsStr := marketFlags.String("coins", "bitcoin", "Space separated cryptocurrency symbols(ex. bitcoin, etheruem). Refer to coins list command.")
	currency := marketFlags.String("currency", "usd", "The target currency of the market data(ex. usd). One currency only.")

	if err := marketFlags.Parse(os.Args[2:]); err != nil {
		return errors.Wrap(err, "market flags")
	}

	spinner.Start()

	var coinsSlice []string
	if len(*coinsStr) > 0 {
		coinsSlice = strings.Split(*coinsStr, " ")
	}

	fmt.Println(coinsSlice)
	coinGeckoClient := coingecko.NewClient(nil)
	marketRes, _, err := coinGeckoClient.Coins.GetMarkets(*currency, coinsSlice)
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "market result")
	}

	market := market{
		currency:    *currency,
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
	var md string
	md += "# Market Data \n"
	md += "🕔 " + time.Now().Format(time.ANSIC) + " \n"
	md += "## Currency: " + strings.ToUpper(m.currency) + " \n"
	md += "| | Name | Price | 1h | 24h | 7d | 14d | 30d |"
	md += "\n"
	md += "| --- | --- | --- | --- | --- | --- | --- | --- | \n"

	for _, marketData := range *m.coinsMarket {
		currentPrice := coinvalue.New(marketData.CurrentPrice)
		pcp1h := coinvalue.New(marketData.PriceChangePercentage1hInCurrency)
		pcp24h := coinvalue.New(marketData.PriceChange24H)
		pcp7D := coinvalue.New(marketData.PriceChangePercentage7DInCurrency)
		pcp14D := coinvalue.New(marketData.PriceChangePercentage14DInCurrency)
		pcp30D := coinvalue.New(marketData.PriceChangePercentage30DInCurrency)

		md += "| " + strings.ToUpper(marketData.Symbol) + " | " +
			strings.ToUpper(marketData.Name) + " | " +
			currentPrice.Value + " | " +
			pcp1h.Emojicon + " " + pcp1h.Value + " | " +
			pcp24h.Emojicon + " " + pcp24h.Value + " | " +
			pcp7D.Emojicon + " " + pcp7D.Value + " | " +
			pcp14D.Emojicon + " " + pcp14D.Value + " | " +
			pcp30D.Emojicon + " " + pcp30D.Value + " | \n"

	}

	glam, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)

	out, err := glam.Render(md)
	if err != nil {
		return "", errors.Wrap(err, "glamour render")
	}

	return out, nil
}
