package kraken

import "github.com/shopspring/decimal"

type CurrencyPair struct {
	Name   string
	Amount decimal.Decimal
}

type ticker struct {
	Ask          []string `json:"a"`
	Bid          []string `json:"b"`
	LastTrade    []string `json:"c"`
	Volume       []string `json:"v"`
	VWAP         []string `json:"p"`
	Trades       []int    `json:"t"`
	Low          []string `json:"l"`
	High         []string `json:"h"`
	OpeningPrice string   `json:"o"`
}

type result map[string]ticker

type krakenResponse struct {
	Error  []interface{} `json:"error"`
	Result result        `json:"result"`
}
