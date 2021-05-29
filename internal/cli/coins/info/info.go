package info

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
	CmdName                 = "info"
	coinFlagValue           = "bitcoin"
	coinFlagDescription     = "cryptocurrency symbol. Refer to coins list command."
	currencyFlagValue       = "usd"
	currencyFlagDescription = "The target currency of the market data(ex. usd). One currency only."
)

type coinInfo struct {
	currency string
	config   *config.Config
	coin     *coingecko.Coin
}

func Cmd() error {
	infoFlags := flag.NewFlagSet(CmdName, flag.ExitOnError)

	var config config.Config
	config.Setup(infoFlags)
	config.YakSpinConfig.Suffix = " Getting Crypto Info from the Gecko"

	spinner, err := yacspin.New(config.YakSpinConfig)
	if err != nil {
		return errors.Wrap(err, "coin info spinner")
	}

	var coinStr string
	var currency string

	infoFlags.StringVar(&coinStr, "coin", coinFlagValue, coinFlagDescription)
	infoFlags.StringVar(&coinStr, "c", coinFlagValue, coinFlagDescription+" (shorthand)")
	infoFlags.StringVar(&currency, "currency", currencyFlagValue, currencyFlagDescription)
	infoFlags.StringVar(&currency, "cr", currencyFlagValue, currencyFlagDescription+" (shorthand)")

	if err := infoFlags.Parse(os.Args[2:]); err != nil {
		return errors.Wrap(err, "coin info flags")
	}
	spinner.Start()

	coinGeckoClient := coingecko.NewClient(nil)
	coinInfoRes, _, err := coinGeckoClient.Coins.GetCoin(coinStr)
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "coin info result")
	}

	coinInfo := coinInfo{
		currency: currency,
		config:   &config,
		coin:     coinInfoRes,
	}

	out, err := coinInfo.generateInfoScreen()
	if err != nil {
		spinner.StopFail()
		return errors.Wrap(err, "generate info screen")
	}

	spinner.Stop()
	fmt.Println(out)

	return nil
}

func (i *coinInfo) generateInfoScreen() (string, error) {
	var md string
	md += "# " + i.coin.Name + " (" + i.coin.Symbol + ") \n"
	// date time
	md += "🕔  " + time.Now().Format(time.ANSIC) + " \n\n"
	// current price / percentage
	md += "## Current Price - " + strings.ToUpper(i.currency) + " \n"
	currentPriceFloat := i.coin.MarketData.CurrentPrice[i.currency]
	currentPrice := coinvalue.New(&currentPriceFloat)

	PCP24HFloat := i.coin.MarketData.PriceChangePercentage24HInCurrency[i.currency]
	PCP24H := coinvalue.New(&PCP24HFloat)

	// market cap
	currencyDisplay := strings.ToUpper(i.currency)
	md += "- " + currencyDisplay + " **" + currentPrice.Value + "** " + PCP24H.Emojicon + " *" + PCP24H.Value +
		"% in the last 24 hours* \n"

	marketCapFloat := i.coin.MarketData.MarketCap[i.currency]
	marketCap := coinvalue.New(&marketCapFloat)

	marketCapRank := strconv.Itoa(int(i.coin.MarketCapRank))
	md += "## Market Cap - Rank " + marketCapRank + "\n"
	md += "- " + currencyDisplay + " " + marketCap.Value + " \n\n"

	circulatingSupply := coinvalue.New(&i.coin.MarketData.CirculatingSupply)
	md += "- Circulating Supply: " + circulatingSupply.Value + " \n"

	// 24 hour info
	md += "## 24 Hour Update \n\n"

	tradingVolumeFloat := i.coin.MarketData.TotalVolume[i.currency]
	tradingVolume := coinvalue.New(&tradingVolumeFloat)

	md += "- Trading Vol. " + tradingVolume.Value + " \n\n"
	high24HFloat := i.coin.MarketData.High24H[i.currency]
	high24H := coinvalue.New(&high24HFloat)

	low24HFloat := i.coin.MarketData.Low24H[i.currency]
	low24H := coinvalue.New(&low24HFloat)
	md += "| 24h Low | 24h High |"
	md += "\n"
	md += "| --- | --- | \n"
	md += "| " + low24H.Value + " | " + high24H.Value + "| \n"

	// links - website
	md += "## 🔗 Links \n"

	linksMarkDown := i.BuildLinksMarkDown()
	md += linksMarkDown

	// TODO add description only on full details sub command
	// md += "## 📖 Description \n"
	// md += i.coin.Description["en"]

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

func (i *coinInfo) BuildLinksMarkDown() string {
	var linksMarkDown strings.Builder
	links := i.coin.Links

	linksMarkDown.WriteString("\n Website \n")

	// home page
	var homePageLinksStringBuilder strings.Builder
	for _, homePageLink := range links.HomePage {
		if len(homePageLink) == 0 {
			continue
		}
		homePageLinkMD := "- [Homepage](" + homePageLink + ") \n"
		homePageLinksStringBuilder.WriteString(homePageLinkMD)
		// break on first value that has a home page
		break
	}
	linksMarkDown.WriteString(homePageLinksStringBuilder.String())

	linksMarkDown.WriteString("\n Explorers \n")

	// block chain site
	var blockChainLinksStringBuilder strings.Builder
	for _, blockChainSiteLink := range links.BlockChainSite {
		if len(blockChainSiteLink) == 0 {
			continue
		}
		blockChainSiteLinkMD := fmt.Sprintf("- [Blockchain Site](%s) \n", blockChainSiteLink)
		blockChainLinksStringBuilder.WriteString(blockChainSiteLinkMD)
	}
	linksMarkDown.WriteString(blockChainLinksStringBuilder.String())

	linksMarkDown.WriteString("\n Community \n")

	// official forum
	var officialForumLinkStringBuilder strings.Builder
	for _, officialForumLink := range links.OfficialForumURL {
		if len(officialForumLink) == 0 {
			continue
		}
		officialForumLinkMD := fmt.Sprintf("- [Official Forum](%s) \n", officialForumLink)
		officialForumLinkStringBuilder.WriteString(officialForumLinkMD)
	}
	linksMarkDown.WriteString(officialForumLinkStringBuilder.String())

	// twitter
	if len(links.TwitterScreenName) > 0 {
		twitterLinkMD := fmt.Sprintf("- [Twitter](https://twitter.com/%s) \n", links.TwitterScreenName)
		linksMarkDown.WriteString(twitterLinkMD)
	}

	// facebook
	if len(links.FacebookUsername) > 0 {
		facebookLinkMD := fmt.Sprintf("- [Facebook](https://facebook.com/%s) \n", links.FacebookUsername)
		linksMarkDown.WriteString(facebookLinkMD)
	}

	// subreddit
	if len(links.SubredditURL) > 0 {
		subredditLinkMD := fmt.Sprintf("- [Subreddit](%s) \n", links.SubredditURL)
		linksMarkDown.WriteString(subredditLinkMD)
	}

	linksMarkDown.WriteString("\n Source Code \n")

	var reposLinkStringBuilder strings.Builder
	// github repos url
	for _, githubRepoURL := range links.ReposURL.Github {
		if len(githubRepoURL) == 0 {
			continue
		}
		githubRepoURLMD := fmt.Sprintf("- [Github](%s) \n", githubRepoURL)
		reposLinkStringBuilder.WriteString(githubRepoURLMD)
	}
	// bitbucket url
	for _, bitbucketRepoURL := range links.ReposURL.Bitbucket {
		if len(bitbucketRepoURL) == 0 {
			continue
		}
		bitbucketRepoURLMD := fmt.Sprintf("- [Bitbucket](%s) \n", bitbucketRepoURL)
		reposLinkStringBuilder.WriteString(bitbucketRepoURLMD)
	}
	linksMarkDown.WriteString(reposLinkStringBuilder.String())

	return linksMarkDown.String()
}
