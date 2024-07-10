package main

import (
	"log/slog"

	"github.com/skantay/blockchange-sentinel/internal/app/api"
)

func main() {
	if err := api.Run(); err != nil {
		slog.Error(err.Error())
	}
}
