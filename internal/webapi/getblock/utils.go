package getblock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrTooManyRequests = errors.New("too many requests")

const (
	includeTransactions = true
	defaultID           = "getblock.io"
	defaultJSONRPC      = "2.0"
	applicationJSON     = "application/json"
)

type request struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      string `json:"id"`
}

func (c *Client) post(data request) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	resp, err := c.h.Post(c.url, applicationJSON, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrTooManyRequests
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}

	return raw, nil
}
