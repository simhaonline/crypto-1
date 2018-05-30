package crypto

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/meeDamian/crypto/currencies"
	"github.com/meeDamian/crypto/orderbook"
)

type (
	Exchange struct {
		// each exchange SHOULD specify its human-readable name
		Name string

		// each exchange MUST specify its generic top level domain (`www.`, etc should be omitted)
		Domain string

		/**
		 * public
		**/
		// returns OrderBook for requested Market or error
		OrderBook func(Market) (orderbook.OrderBook, error)

		// returns a list of all Markets on available on a given exchange. Includes disabled Markets.
		//      Limited to supported currencies only, see currencies/currencies.go and currencies/symbols/symbols.go for more
		Markets func() ([]Market, error)

		/**
		 * private
		**/
		// returns all exchange Balances (Total, Available, Locked) for account credentials provided
		Balances func(Credentials) (Balances, error)

		/**
		 * optional
		**/
		// returns OrderBooks of ALL available markets. Should only be implemented if a "shortcut" endpoint exists
		//      If only some markets couldn't be downloaded, error should be logged, but not returned.
		//      Error only if no usable data can be returned
		AllOrderBooks func() ([]orderbook.OrderBook, error)
	}

	Credentials struct {
		// human-readable Name of the account
		Name string

		// API Key & Secret
		Key, Secret string

		// Id needed by bitstamp and tdax
		Id *string
	}
)

var pairRegExp regexp.Regexp

func init() {
	var symbols []string
	for symbol := range currencies.All() {
		symbols = append(symbols, symbol)
	}

	orSymbols := strings.Join(symbols, "|")
	pairRegExp = *regexp.MustCompile(fmt.Sprintf(`(?i)^[ZX]?(%[1]s)\/?[_ZX]?(%[1]s)$`, orSymbols))
}
