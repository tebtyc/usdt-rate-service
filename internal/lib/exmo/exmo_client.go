package exmo

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
	"usdt-rate-service/internal/domain"
)

const exmoURL = "https://api.exmo.com/v1.1/order_book?pair=USDT_USD&limit=1"

type ExmoClient struct {
	httpClient *http.Client
}

func New() *ExmoClient {
	return &ExmoClient{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

type ResponseExmo struct {
	USDT_USD struct {
		Ask [][]string `json:"ask"`
		Bid [][]string `json:"bid"`
	} `json:"USDT_USD"`
}

func (c *ExmoClient) FetchRates(ctx context.Context) (domain.Rate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, exmoURL, nil)
	if err != nil {
		return domain.Rate{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Rate{}, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Rate{}, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var data ResponseExmo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return domain.Rate{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.USDT_USD.Ask) == 0 || len(data.USDT_USD.Bid) == 0 {
		return domain.Rate{}, fmt.Errorf("empty ask or bid")
	}

	askPrice, err := strconv.ParseFloat(data.USDT_USD.Ask[0][0], 64)
	if err != nil {
		return domain.Rate{}, fmt.Errorf("failed to parse ask price: %w", err)
	}
	bidPrice, err := strconv.ParseFloat(data.USDT_USD.Bid[0][0], 64)
	if err != nil {
		return domain.Rate{}, fmt.Errorf("failed to parse bid price: %w", err)
	}

	return domain.Rate{
		Ask:       askPrice,
		Bid:       bidPrice,
		Timestamp: time.Now().Unix(),
	}, nil
}
