package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/skantay/blockchange-sentinel/internal/webapi/getblock"
)

type service interface {
	GetMostChangedAddress(context.Context) (string, error)
}

type API struct {
	service service
}

func New(service service) *API {
	return &API{
		service: service,
	}
}

func (a *API) GetMostChangedAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	address, err := a.service.GetMostChangedAddress(r.Context())
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, getblock.ErrTooManyRequests) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"most_changed_address": address,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
