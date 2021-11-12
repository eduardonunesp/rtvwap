package internals

import "math/big"

// TradePair is the representation of the trading pair left and right pairs
// For instance ETH-USD, BTC-USD, USD-ETH
type TradePair struct {
	Left  string
	Right string
}

// NewTradePair create new trade pair
func NewTradePair(lPair, rPair string) TradePair {
	return TradePair{lPair, rPair}
}

// Trade represents a trade that matched/closed on the provider
type Trade struct {
	// TODO: Add some ID to indicate which provider was used
	TradePair
	Price    *big.Float
	Quantity *big.Float
}

// TradeProvider is the core structure for the provider which serves the go channel
// The go channel produces the trading stream information
type TradeProvider struct {
	TradeChan chan Trade
}

// TradeProviderCreator represents the interfae that all trade providers must have
// in order to create the TradeProvider
type TradeProviderCreator interface {
	// Should pass the trade pair expected from this trade provider
	CreateTradeProvider(pair TradePair) (TradeProvider, error)
}

// TradeFeed produces the trading information based on the given provider
type TradeFeed struct {
	tradeProviderChan TradeProvider
}

// NewTradeFeed creates a new trade feed passing the trade pair and the trade provider to use
func NewTradeFeed(lPair, rPair string, tradeProvider TradeProviderCreator) (TradeFeed, error) {
	return NewTradeFeedWithPair(TradePair{lPair, rPair}, tradeProvider)
}

// NewTradeFeedWithPair creates a new trade feed passing the trade pair and the trade provider to use
func NewTradeFeedWithPair(pair TradePair, tradeProvider TradeProviderCreator) (TradeFeed, error) {
	tProvider := TradeFeed{}

	tProviderChan, err := tradeProvider.CreateTradeProvider(pair)
	if err != nil {
		return tProvider, err
	}

	tProvider.tradeProviderChan = tProviderChan

	return tProvider, nil
}

// GetFeedChan will return the channel the trading information
func (tf TradeFeed) TradeChan() <-chan Trade {
	return tf.tradeProviderChan.TradeChan
}
