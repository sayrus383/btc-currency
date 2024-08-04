package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sayrus383/btc-currency/internal/service/kraken"
	"log/slog"
	"net/http"
	"strings"
)

type Server struct {
	log       *slog.Logger
	addr      string
	krakenSvc kraken.Service

	echo *echo.Echo
}

func NewServer(log *slog.Logger, addr string, krakenSvc kraken.Service) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		echo:      e,
		log:       log,
		addr:      addr,
		krakenSvc: krakenSvc,
	}
}

func (s *Server) Init() error {
	s.echo.GET("/api/v1/ltp", s.ltp)

	return nil
}

func (s *Server) Start(_ context.Context) error {
	s.log.Info(fmt.Sprintf("http server start on %s", s.addr))

	return s.echo.Start(s.addr)
}

func (s *Server) Stop() {
	_ = s.echo.Close()
}

func (s *Server) ltp(c echo.Context) error {
	pairs := strings.Split(c.QueryParam("pairs"), ",")
	if c.QueryParam("pairs") == "" || len(pairs) == 0 {
		return c.JSON(http.StatusBadRequest, ToError("pairs is empty"))
	}

	currencyPairs := s.krakenSvc.GetCurrencyPairs(pairs)

	return c.JSON(http.StatusOK, KrakenCurrencyPairToHttpResponse(currencyPairs))
}
