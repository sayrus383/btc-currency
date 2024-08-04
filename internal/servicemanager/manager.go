package servicemanager

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type CanInit interface {
	Init() error
}

type CanStart interface {
	Start(ctx context.Context) error
}

type CanStop interface {
	Stop()
}

func Register(ctx context.Context, log *slog.Logger, services ...any) error {
	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	for _, svc := range services {
		if s, ok := svc.(CanInit); ok {
			if err := s.Init(); err != nil {
				return err
			}
		}
	}

	// Graceful Start
	for _, svc := range services {
		go func(svc any) {
			if s, ok := svc.(CanStart); ok {
				if err := s.Start(ctx); err != nil {
					errChan <- fmt.Errorf("graceful start error: %v", err)
				}
			}
		}(svc)
	}

	// Graceful Stop
	for _, svc := range services {
		defer func(svc any) {
			if s, ok := svc.(CanStop); ok {
				s.Stop()
			}
		}(svc)
	}

	select {
	case err := <-errChan:
		return err
	case stop := <-stopChan:
		log.Debug("app is finished", "signal", stop.String())
	}

	return nil
}
