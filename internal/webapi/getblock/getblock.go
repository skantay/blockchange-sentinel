package getblock

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/skantay/blockchange-sentinel/internal/entities"
)

type Client struct {
	h   *http.Client
	url string
}

func New(h *http.Client, APIkey string) *Client {
	return &Client{
		h:   h,
		url: "https://go.getblock.io/" + APIkey,
	}
}

type responseGetLastBlock struct {
	Result string `json:"result"`
}

func (c *Client) GetLastBlock() (string, error) {
	raw := request{
		JSONRPC: defaultJSONRPC,
		Method:  "eth_blockNumber",
		Params:  []any{},
		ID:      defaultID,
	}

	data, err := c.post(raw)
	if err != nil {
		return "", err
	}

	var result responseGetLastBlock
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal: %w", err)
	}

	return result.Result, nil
}

type responseGetBlockInfoByIndex struct {
	Result struct {
		Transactions []entities.Transaction `json:"transactions"`
	} `json:"result"`
}

func (c *Client) GetBlockTransactionsByIndex(block string, index int) ([]entities.Transaction, error) {
	block = getBlockByIndex(block, index)

	raw := request{
		JSONRPC: defaultJSONRPC,
		Method:  "eth_getBlockByNumber",
		Params: []any{
			block,
			includeTransactions,
		},
		ID: defaultID,
	}

	data, err := c.post(raw)
	if err != nil {
		return nil, err
	}

	var result responseGetBlockInfoByIndex
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return result.Result.Transactions, nil
}

func getBlockByIndex(block string, index int) string {
	x := &big.Int{}
	x.SetString(block[2:], 16)

	x.Sub(x, big.NewInt(int64(index)))

	return fmt.Sprintf("0x%x", x)
}
