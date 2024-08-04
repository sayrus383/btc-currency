package app

import (
	"context"
	"fmt"
	"github.com/namsral/flag"
	"github.com/sayrus383/btc-currency/internal/service/kraken"
	"github.com/sayrus383/btc-currency/internal/servicemanager"
	"github.com/sayrus383/btc-currency/internal/transport/http"
	"github.com/sayrus383/btc-currency/pkg/logger"
	"strings"
)

func Run() error {
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	log := logger.New()
	log.Info("application starting...")

	var pairFollows string
	flag.StringVar(&pairFollows, "KRAKEN_PAIR_FOLLOWS", "BTC/USD,BTC/CHF,BTC/EUR", "pair follows environment variable")
	kraken := kraken.New(log, strings.Split(pairFollows, ","))

	var port int
	flag.IntVar(&port, "HTTP_SERVER_PORT", 80, "HTTP server port")
	srv := http.NewServer(log, fmt.Sprintf(":%d", port), kraken)

	err := servicemanager.Register(ctx, log, srv, kraken)
	if err != nil {
		return err
	}

	return nil
}
