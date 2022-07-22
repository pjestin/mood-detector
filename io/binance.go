package io

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pjestin/mood-detector/model/binance"
	"github.com/pjestin/mood-detector/util"
)

type BinanceClient struct {
	apiKey    string
	secretKey string
	client    http.Client
}

func (c *BinanceClient) Init() {
	c.apiKey = os.Getenv("BINANCE_API_KEY")
	c.secretKey = os.Getenv("BINANCE_SECRET_KEY")
}

func (c *BinanceClient) buildParamString(params map[string]string) string {
	var paramString string
	for k, v := range params {
		if len(paramString) > 0 {
			paramString += "&"
		}
		paramString += fmt.Sprintf("%s=%s", k, v)
	}
	paramString += fmt.Sprintf("&timestamp=%s", fmt.Sprint(time.Now().UnixMilli()))

	signature := util.HmacSha256(paramString, c.secretKey)

	paramString += fmt.Sprintf("&signature=%s", signature)
	return paramString
}

func (c *BinanceClient) SendRequest(method string, path string, params map[string]string) ([]byte, error) {
	paramString := c.buildParamString(params)

	req, err := http.NewRequest(method, fmt.Sprintf("https://api.binance.com%s?%s", path, paramString), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "MoodDetector/0.1.0 (by pjestin)")
	req.Header.Set("X-MBX-APIKEY", c.apiKey)

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error status code: %d", resp.StatusCode)
		} else {
			return nil, fmt.Errorf("error status code: %d, response body: %s", resp.StatusCode, b)
		}
	}

	b, err := io.ReadAll(resp.Body)
	return b, err
}

func (c *BinanceClient) GetLastTrade(symbol string) (binance.Trade, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  "1",
	}

	b, err := c.SendRequest("GET", "/api/v3/myTrades", params)

	if err != nil {
		return binance.Trade{}, err
	}

	var body []binance.Trade
	err = json.Unmarshal(b, &body)

	return body[0], err
}

func (c *BinanceClient) PostOrder(symbol string, orderType string, side string, quantity string) (binance.Order, error) {
	params := map[string]string{
		"symbol":   symbol,
		"type":     orderType,
		"side":     side,
		"quantity": quantity,
	}

	b, err := c.SendRequest("POST", "/api/v3/order", params)

	if err != nil {
		return binance.Order{}, err
	}

	var body binance.Order
	err = json.Unmarshal(b, &body)

	return body, err
}
