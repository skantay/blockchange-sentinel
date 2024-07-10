package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/skantay/blockchange-sentinel/internal/entities"
	"github.com/skantay/blockchange-sentinel/internal/webapi/getblock"
)

const (
	nBlocks = 100
	retries = 3
)

type blockService interface {
	GetLastBlock() (string, error)
	GetBlockTransactionsByIndex(string, int) ([]entities.Transaction, error)
}

type Service struct {
	blockService blockService
}

func New(blockService blockService) *Service {
	return &Service{
		blockService: blockService,
	}
}

func (s *Service) GetMostChangedAddress(ctx context.Context) (string, error) {
	var recentBlock string

	var err error

	for i := 0; i < retries; i++ {
		recentBlock, err = s.blockService.GetLastBlock()
		if err != nil {
			if errors.Is(err, getblock.ErrTooManyRequests) {
				continue
			}
			return "", fmt.Errorf("failed to get recent block: %w", err)
		}
		break
	}

	if recentBlock == "" {
		return "", err
	}

	transactions := make(map[string]*big.Int)

	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}
	errChan := make(chan error)

	for i := 0; i < nBlocks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var block []entities.Transaction

			for i := 0; i < retries; i++ {
				block, err = s.blockService.GetBlockTransactionsByIndex(recentBlock, i)
				if err != nil {
					if errors.Is(err, getblock.ErrTooManyRequests) {
						continue
					}

					select {
					case errChan <- fmt.Errorf("failed to get transactions by index: %w", err):
					default:
					}
					return
				}
				break
			}

			if len(block) == 0 {
				select {
				case errChan <- getblock.ErrTooManyRequests:
				default:
				}
				return
			}

			processBlock(block, transactions, m)

		}()
	}

	doneChan := make(chan struct{})

	go func() {
		wg.Wait()

		doneChan <- struct{}{}
	}()

	for {
		select {
		case err := <-errChan:
			return "", err
		case <-ctx.Done():
			return "", ctx.Err()
		case <-doneChan:
			return mostChangedAddress(transactions), nil
		}
	}
}

func processBlock(blockTransactions []entities.Transaction, allTransactions map[string]*big.Int, m *sync.Mutex) {
	wg := &sync.WaitGroup{}

	for _, tx := range blockTransactions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			amount := &big.Int{}
			amount.SetString(tx.Value[2:], 16)

			m.Lock()

			if _, ok := allTransactions[tx.From]; !ok {
				allTransactions[tx.From] = &big.Int{}
			}
			allTransactions[tx.From].Sub(allTransactions[tx.From], amount)

			if _, ok := allTransactions[tx.To]; !ok {
				allTransactions[tx.To] = &big.Int{}
			}
			allTransactions[tx.To].Add(allTransactions[tx.To], amount)

			m.Unlock()
		}()
	}

	wg.Wait()
}

func mostChangedAddress(addresses map[string]*big.Int) string {
	var result string

	var tmp *big.Int

	for address, amount := range addresses {
		if tmp == nil || tmp.Abs(tmp).Cmp(amount.Abs(amount)) < 0 {
			tmp = amount

			result = address
		}
	}

	return result
}

func retry(fn func() bool) {
	retry := fn()
	for !retry {
		retry = fn()
	}
}
