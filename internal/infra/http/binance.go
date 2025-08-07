package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Homyakadze14/AFFARM_tz/internal/common"
)

type BinanceClient struct {
	log     *slog.Logger
	client  *http.Client
	baseURL string
}

func NewBinanceClient(log *slog.Logger, timeout time.Duration) *BinanceClient {
	client := &http.Client{
		Timeout: timeout,
	}

	baseURL := "https://api.binance.com/api/v3"

	return &BinanceClient{
		log:     log,
		client:  client,
		baseURL: baseURL,
	}
}

type Cryptocurrency struct {
	Price string `json:"price"`
}

func (c *BinanceClient) GetPrice(symbol string, currency string) (float64, error) {
	const op = "BinanceClient.GetPrice"
	log := c.log.With(slog.String("op", op),
		slog.String("symbol", symbol),
		slog.String("currency", currency))

	priceURL := c.baseURL + fmt.Sprintf("/ticker/price?symbol=%s%s", symbol, currency)
	resp, err := c.client.Get(priceURL)
	if err != nil {
		log.Error(fmt.Sprintf("fail to get price! error: %s", err))
		return 0, common.ErrUnexpected
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("bad status code! code: %s", resp.Status))
		if resp.StatusCode == http.StatusBadRequest {
			return 0, common.ErrBadData
		}
		return 0, common.ErrUnexpected
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("fail to read response body! error: %s", err))
		return 0, common.ErrUnexpected
	}

	cryptocur := &Cryptocurrency{}
	err = json.Unmarshal(data, cryptocur)
	if err != nil {
		log.Error(fmt.Sprintf("fail to unmarshal response body! error: %s", err))
		return 0, common.ErrUnexpected
	}

	price, err := strconv.ParseFloat(cryptocur.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func (c *BinanceClient) SymbolExists(symbol string) (bool, error) {
	const op = "BinanceClient.SymbolExists"
	log := c.log.With(slog.String("op", op),
		slog.String("symbol", symbol))

	_, err := c.GetPrice(symbol, "USDT")
	if err != nil {
		if errors.Is(err, common.ErrBadData) {
			return false, nil
		}
		log.Error(fmt.Sprintf("fail to check existence! error: %s", err))
		return false, err
	}

	return true, nil
}
