package http

import "github.com/sayrus383/btc-currency/internal/service/kraken"

type Error struct {
	Error string `json:"error"`
}

type Price struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
}

type LTPResponse struct {
	LTP []Price `json:"ltp"`
}

func ToError(msg string) Error {
	return Error{Error: msg}
}

func KrakenCurrencyPairToHttpResponse(pairs []kraken.CurrencyPair) LTPResponse {
	httpResponse := LTPResponse{
		LTP: make([]Price, 0, len(pairs)),
	}

	for _, pair := range pairs {
		httpResponse.LTP = append(httpResponse.LTP, Price{
			Pair:   pair.Name,
			Amount: pair.Amount.String(),
		})
	}

	return httpResponse
}
