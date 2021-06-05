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
	var markdownBuilder strings.Builder
	var markdown string
	markdown = fmt.Sprintf("# %s (%s) \n", i.coin.Name, i.coin.Symbol)
	markdownBuilder.WriteString(markdown)

	// date time
	markdown = fmt.Sprintf("🕔 %s \n\n", time.Now().Format(time.ANSIC))
	markdownBuilder.WriteString(markdown)

	// current price / percentage
	currencyDisplay := strings.ToUpper(i.currency)
	markdown = fmt.Sprintf("## Current Price - %s \n", currencyDisplay)
	markdownBuilder.WriteString(markdown)

	currentPriceFloat := i.coin.MarketData.CurrentPrice[i.currency]
	currentPrice := coinvalue.New(&currentPriceFloat)

	PCP24HFloat := i.coin.MarketData.PriceChangePercentage24HInCurrency[i.currency]
	PCP24H := coinvalue.New(&PCP24HFloat)

	// market cap
	markdown = fmt.Sprintf("- %s **%s** %s *%s%% in the last 24 hours* \n", currencyDisplay,
		currentPrice.Value, PCP24H.Emojicon, PCP24H.Value)
	markdownBuilder.WriteString(markdown)

	marketCapFloat := i.coin.MarketData.MarketCap[i.currency]
	marketCap := coinvalue.New(&marketCapFloat)

	marketCapRank := strconv.Itoa(int(i.coin.MarketCapRank))
	markdown = fmt.Sprintf("## Market Cap - Rank %s \n", marketCapRank)
	markdownBuilder.WriteString(markdown)

	markdown = fmt.Sprintf("- %s %s \n\n", currencyDisplay, marketCap.Value)
	markdownBuilder.WriteString(markdown)

	circulatingSupply := coinvalue.New(&i.coin.MarketData.CirculatingSupply)
	markdown = fmt.Sprintf("- Circulating Supply: %s \n", circulatingSupply.Value)
	markdownBuilder.WriteString(markdown)

	// 24 hour info
	markdownBuilder.WriteString("## 24 Hour Update \n\n")

	tradingVolumeFloat := i.coin.MarketData.TotalVolume[i.currency]
	tradingVolume := coinvalue.New(&tradingVolumeFloat)

	markdown = fmt.Sprintf("- Trading Vol. %s \n\n", tradingVolume.Value)
	markdownBuilder.WriteString(markdown)

	high24HFloat := i.coin.MarketData.High24H[i.currency]
	high24H := coinvalue.New(&high24HFloat)

	low24HFloat := i.coin.MarketData.Low24H[i.currency]
	low24H := coinvalue.New(&low24HFloat)
	// md += "| 24h Low | 24h High |"
	// md += "\n"
	// md += "| --- | --- | \n"
	markdown = "| 24h Low | 24h High | \n | --- | --- | \n"
	markdownBuilder.WriteString(markdown)

	// md += "| " + low24H.Value + " | " + high24H.Value + "| \n"
	markdown = fmt.Sprintf("| %s | %s | \n", low24H.Value, high24H.Value)
	markdownBuilder.WriteString(markdown)

	linksMarkDown := i.BuildLinksMarkDown()
	markdownBuilder.WriteString(linksMarkDown)

	// TODO add description only on full details sub command
	// md += "## 📖 Description \n"
	// md += i.coin.Description["en"]

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

func (i *coinInfo) BuildLinksMarkDown() string {
	var linksMarkDown strings.Builder
	links := i.coin.Links

	linksMarkDown.WriteString("## 🌏 Website \n")

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

	linksMarkDown.WriteString("## 🔎 Explorers \n")

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

	linksMarkDown.WriteString("## 🗣  Community \n")

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

	linksMarkDown.WriteString("## 🧑‍💻 Source Code \n")

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
