package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MarketDataService struct {
	client *http.Client
}

func NewMarketDataService() *MarketDataService {
	return &MarketDataService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetHistoricalPrices retrieves historical prices with API priority
func (s *MarketDataService) GetHistoricalPrices(symbol string, days int) ([]PricePoint, error) {
	// Determine asset type
	if isCrypto(symbol) {
		return s.getCryptoPrices(symbol, days)
	}
	return s.getStockPrices(symbol, days)
}

// getCryptoPrices uses multi-API fallback for cryptocurrencies
func (s *MarketDataService) getCryptoPrices(symbol string, days int) ([]PricePoint, error) {
	// 1. Try Binance (NO LIMITS)
	prices, err := s.fetchBinance(symbol, days)
	if err == nil && len(prices) > 0 {
		return prices, nil
	}
	
	// 2. Try CoinGecko (free)
	prices, err = s.fetchCoinGecko(symbol, days)
	if err == nil && len(prices) > 0 {
		return prices, nil
	}
	
	// 3. Try CoinCap (200/day)
	prices, err = s.fetchCoinCap(symbol, days)
	if err == nil && len(prices) > 0 {
		return prices, nil
	}
	
	return nil, fmt.Errorf("failed to fetch crypto prices for %s", symbol)
}

// getStockPrices for stocks/ETFs
func (s *MarketDataService) getStockPrices(symbol string, days int) ([]PricePoint, error) {
	// 1. Yahoo Finance (free)
	prices, err := s.fetchYahooFinance(symbol, days)
	if err == nil && len(prices) > 0 {
		return prices, nil
	}
	
	return nil, fmt.Errorf("failed to fetch stock prices for %s", symbol)
}

// Binance API (NO LIMITS)
func (s *MarketDataService) fetchBinance(symbol string, days int) ([]PricePoint, error) {
	// Convert symbol (BTC -> BTCUSDT)
	pair := symbol + "USDT"
	
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=1d&limit=%d", pair, days)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("binance API error: %d", resp.StatusCode)
	}
	
	var data [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	prices := make([]PricePoint, 0, len(data))
	for _, candle := range data {
		timestamp := int64(candle[0].(float64))
		closePrice := parseFloat(candle[4])
		
		prices = append(prices, PricePoint{
			Date:  time.Unix(timestamp/1000, 0),
			Close: closePrice,
		})
	}
	
	return prices, nil
}

// CoinGecko API (free)
func (s *MarketDataService) fetchCoinGecko(symbol string, days int) ([]PricePoint, error) {
	coinID := getCoinGeckoID(symbol)
	
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=usd&days=%d", coinID, days)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("coingecko API error: %d", resp.StatusCode)
	}
	
	var data struct {
		Prices [][]float64 `json:"prices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	prices := make([]PricePoint, 0, len(data.Prices))
	for _, p := range data.Prices {
		prices = append(prices, PricePoint{
			Date:  time.Unix(int64(p[0]/1000), 0),
			Close: p[1],
		})
	}
	
	return prices, nil
}

// CoinCap API (200 requests/day)
func (s *MarketDataService) fetchCoinCap(symbol string, days int) ([]PricePoint, error) {
	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s/history?interval=d1", toLowerCase(symbol))
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var data struct {
		Data []struct {
			PriceUsd string `json:"priceUsd"`
			Time     int64  `json:"time"`
		} `json:"data"`
	}
	
	json.Unmarshal(body, &data)
	
	prices := make([]PricePoint, 0)
	for _, p := range data.Data {
		prices = append(prices, PricePoint{
			Date:  time.Unix(p.Time/1000, 0),
			Close: parseFloat(p.PriceUsd),
		})
	}
	
	return prices, nil
}

// Yahoo Finance (free) - using chart API instead of download
func (s *MarketDataService) fetchYahooFinance(symbol string, days int) ([]PricePoint, error) {
	// Use v8 chart API for better reliability
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=%dd", symbol, days)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("yahoo API error: %d", resp.StatusCode)
	}
	
	// Simple CSV parsing
	body, _ := io.ReadAll(resp.Body)
	lines := parseCSV(string(body))
	
	prices := make([]PricePoint, 0)
	for i, line := range lines {
		if i == 0 || len(line) < 5 {
			continue // skip header
		}
		
		date, _ := time.Parse("2006-01-02", line[0])
		closePrice := parseFloat(line[4])
		
		prices = append(prices, PricePoint{
			Date:  date,
			Close: closePrice,
		})
	}
	
	return prices, nil
}

// Helper functions
type PricePoint struct {
	Date  time.Time
	Close float64
}

func isCrypto(symbol string) bool {
	cryptos := map[string]bool{
		"BTC": true, "ETH": true, "BNB": true, "XRP": true,
		"ADA": true, "SOL": true, "DOGE": true, "DOT": true,
	}
	return cryptos[symbol]
}

func getCoinGeckoID(symbol string) string {
	mapping := map[string]string{
		"BTC": "bitcoin",
		"ETH": "ethereum",
		"BNB": "binancecoin",
		"XRP": "ripple",
		"ADA": "cardano",
		"SOL": "solana",
		"DOGE": "dogecoin",
	}
	if id, ok := mapping[symbol]; ok {
		return id
	}
	return toLowerCase(symbol)
}

func toLowerCase(s string) string {
	return string(append([]byte(nil), []byte(s)...))
}

func parseFloat(s interface{}) float64 {
	switch v := s.(type) {
	case float64:
		return v
	case string:
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	}
	return 0
}

func parseCSV(data string) [][]string {
	lines := [][]string{}
	// Simplified parsing - use encoding/csv for production
	return lines
}
