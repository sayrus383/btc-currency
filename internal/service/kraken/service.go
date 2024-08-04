package kraken

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const krakenCurrencyPairUrl = "https://api.kraken.com/0/public/Ticker"

type Service interface {
	GetCurrencyPairs([]string) []CurrencyPair
}

type service struct {
	log         *slog.Logger
	pairFollows []string

	pairsCache map[string]CurrencyPair
	httpClient *http.Client
}

func New(log *slog.Logger, pairFollows []string) Service {
	return &service{
		log:         log,
		pairFollows: pairFollows,
		pairsCache:  make(map[string]CurrencyPair),
		httpClient:  &http.Client{},
	}

}

func (s *service) Start(ctx context.Context) error {
	s.log.Info("starting kraken service")
	if err := s.updateCache(ctx); err != nil {
		return err
	}

	// Update currency pair cache every 5 second
	t := time.NewTicker(time.Second * 5)
	for range t.C {
		if err := s.updateCache(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) GetCurrencyPairs(pairs []string) []CurrencyPair {
	curPairs := make([]CurrencyPair, 0, len(pairs))
	for _, pair := range pairs {
		if cache, ok := s.pairsCache[pair]; ok {
			curPairs = append(curPairs, cache)
		}
	}

	return curPairs
}

func (s *service) updateCache(ctx context.Context) error {
	for _, pair := range s.pairFollows {
		kraken, err := s.getPairFromKraken(ctx, pair)
		if err != nil {
			return err
		}

		t, ok := kraken.Result[pair]
		if !ok {
			continue
		}

		amount, err := decimal.NewFromString(t.LastTrade[0])
		if err != nil {
			return err
		}

		s.pairsCache[pair] = CurrencyPair{
			Name:   pair,
			Amount: amount,
		}
	}

	return nil
}

func (s *service) getPairFromKraken(ctx context.Context, pair string) (krakenResponse, error) {
	url := fmt.Sprintf("%s?pair=%s", krakenCurrencyPairUrl, pair)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return krakenResponse{}, fmt.Errorf("getPairFromKraken create req err: %v", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return krakenResponse{}, fmt.Errorf("getPairFromKraken do req err: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return krakenResponse{}, fmt.Errorf("getPairFromKraken read body err: %v", err)
	}

	var response krakenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return krakenResponse{}, fmt.Errorf("getPairFromKraken unmarshal err: %v", err)
	}

	return response, nil
}
