package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/skantay/blockchange-sentinel/internal/controller/cli"
	"github.com/skantay/blockchange-sentinel/internal/service"
	"github.com/skantay/blockchange-sentinel/internal/webapi/getblock"
	"github.com/skantay/blockchange-sentinel/pkg/config"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	getblock := getblock.New(
		&http.Client{},
		cfg.APIKey,
	)

	service := service.New(getblock)

	ctrl := cli.New(service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		slog.Info("\nshutting down")

		cancel()
	}()

	if err := ctrl.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("failed to run: %w", err)
	}

	return nil

}
