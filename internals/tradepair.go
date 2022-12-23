package internals

// TradePair is the representation of the trading pair left and right pairs
// For instance ETH-USD, BTC-USD, USD-ETH
type TradePair struct {
	From string
	To   string
}

// NewTradePair create new trade pair
func NewTradePair(from, to string) TradePair {
	return TradePair{from, to}
}
