package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skantay/blockchange-sentinel/internal/controller/api"
	"github.com/skantay/blockchange-sentinel/internal/service"
	"github.com/skantay/blockchange-sentinel/internal/webapi/getblock"
	"github.com/skantay/blockchange-sentinel/pkg/config"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	getblock := getblock.New(&http.Client{}, cfg.APIKey)
	service := service.New(getblock)
	ctrl := api.New(service)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/block/most_changed", ctrl.GetMostChangedAddress)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.HTTPport),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		slog.Info(fmt.Sprintf("Starting server on port %v\n", cfg.HTTPport))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
		}
	}()

	<-signalCh
	slog.Info("shutting down")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("failed to shutdown the server: %w", err)
	}

	return nil
}
