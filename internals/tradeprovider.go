package internals

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

func NewTradeProvider() TradeProvider {
	return TradeProvider{
		TradeChan: make(chan Trade),
	}
}
