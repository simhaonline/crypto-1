package luno

import (
	"encoding/json"
	"fmt"
	"github.com/meeDamian/crypto/markets"

	"github.com/meeDamian/crypto/currencies"
	"github.com/meeDamian/crypto/orderbook"
	"github.com/meeDamian/crypto/utils"
)

type tickers struct {
	Tickers []struct {
		Pair string `json:"pair"`
	} `json:"tickers"`
}

const (
	orderBookUrl = "https://api.mybitx.com/api/1/orderbook?pair=%s%s"
	marketsUrl   = "https://api.mybitx.com/api/1/tickers"
)

var marketList []markets.Market

func morph(name string) string {
	return currencies.Morph(name, aliases)
}

func OrderBook(m markets.Market) (ob orderbook.OrderBook, err error) {
	url := fmt.Sprintf(orderBookUrl, morph(m.Asset), morph(m.PricedIn))
	return orderbook.Download(url)
}

func Markets() (_ []markets.Market, err error) {
	if len(marketList) > 0 {
		return marketList, nil
	}

	res, err := utils.NetClient().Get(marketsUrl)
	if err != nil {
		return []markets.Market{}, err
	}

	defer res.Body.Close()

	var ts tickers
	err = json.NewDecoder(res.Body).Decode(&ts)
	if err != nil {
		return
	}

	for _, m := range ts.Tickers {
		market, err := markets.NewFromSymbol(m.Pair)
		if err != nil {
			log.Debugf("skipping symbol %s: %v", m.Pair, err)
			continue
		}

		marketList = append(marketList, market)
	}

	return marketList, nil
}
