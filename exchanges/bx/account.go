package bx

import (
	"encoding/json"

	"github.com/meeDamian/crypto"
)

type balResp struct {
	Balances map[string]struct {
		Total     interface{} `json:"total"`
		Available interface{} `json:"available"`
	} `json:"balance"`
}

const balancesUrl = "https://bx.in.th/api/balance/"

func Balances(c crypto.Credentials) (balances crypto.Balances, err error) {
	res, err := privateRequest(c, "POST", balancesUrl, nil)
	if err != nil {
		return
	}

	defer res.Body.Close()

	var r balResp
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return
	}

	balances = make(crypto.Balances)
	for currency, b := range r.Balances {
		err := balances.Add(currency, b.Available, b.Total, nil)
		if err != nil {
			log.Debugf("skipping balance of %s: %v", currency, err)
		}
	}

	return
}
