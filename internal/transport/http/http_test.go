package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sayrus383/btc-currency/internal/service/kraken"
	"github.com/sayrus383/btc-currency/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type KrakenServiceMock struct{}

func (m *KrakenServiceMock) GetCurrencyPairs(_ []string) []kraken.CurrencyPair {
	return []kraken.CurrencyPair{
		{Name: "BTC/CHF", Amount: decimal.RequireFromString("49000.12")},
		{Name: "BTC/EUR", Amount: decimal.RequireFromString("50000.12")},
		{Name: "BTC/USD", Amount: decimal.RequireFromString("52000.12")},
	}
}

func TestLTP(t *testing.T) {
	e := echo.New()
	mockService := &KrakenServiceMock{}
	server := NewServer(logger.New(), ":8080", mockService)
	server.Init()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ltp?pairs=BTC/CHF,BTC/EUR,BTC/USD", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, server.ltp(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response LTPResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := LTPResponse{
			LTP: []Price{
				{Pair: "BTC/CHF", Amount: "49000.12"},
				{Pair: "BTC/EUR", Amount: "50000.12"},
				{Pair: "BTC/USD", Amount: "52000.12"},
			},
		}

		assert.Equal(t, expected, response)
	}
}
