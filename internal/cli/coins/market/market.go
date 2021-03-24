package market

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/markdl-dev/coin-dl/internal/config"
	"github.com/markdl-dev/go-coin-gecko/coingecko"
	"github.com/pkg/errors"
	"github.com/theckman/yacspin"
)

type Market struct {
	market *coingecko.CoinsMarket
	config *config.Config
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

	coinGeckoClient := coingecko.NewClient(nil)
	marketRes, _, err := coinGeckoClient.Coins.GetMarkets(*currency, coinsSlice)
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "market result")
	}
	spinner.Stop()
	fmt.Println(marketRes)
	return nil
}
