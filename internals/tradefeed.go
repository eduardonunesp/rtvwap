package internals

// TradeFeed produces the trading information based on the given provider
type TradeFeed struct {
	tradeProvider TradeProvider
}

// NewTradeFeed creates a new trade feed passing the trade pair and the trade provider to use
func NewTradeFeed(pair TradePair, tradeProviderCreator TradeProviderCreator) (TradeFeed, error) {
	tradeFeed := TradeFeed{}

	tradeProvider, err := tradeProviderCreator.CreateTradeProvider(pair)
	if err != nil {
		return tradeFeed, err
	}

	tradeFeed.tradeProvider = tradeProvider

	return tradeFeed, nil
}

// TradeChan will return the channel the trading information
func (tf TradeFeed) TradeChan() <-chan Trade {
	return tf.tradeProvider.TradeChan
}
