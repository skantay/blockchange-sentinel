package main

import (
	"log/slog"

	"github.com/skantay/blockchange-sentinel/internal/app/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		slog.Error(err.Error())
	}
}
