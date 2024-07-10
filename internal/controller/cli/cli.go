package cli

import (
	"context"
	"fmt"
)

type service interface {
	GetMostChangedAddress(context.Context) (string, error)
}

type CLI struct {
	service service
}

func New(service service) *CLI {
	return &CLI{
		service: service,
	}
}

func (c CLI) Run(ctx context.Context) error {
	address, err := c.service.GetMostChangedAddress(ctx)
	if err != nil {
		return fmt.Errorf("failed to get most changed address: %w", err)
	}

	fmt.Printf("Most changed address: %v\n", address)

	return nil
}
