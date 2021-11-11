package internals

import "math/big"

type Pair struct {
	Left  string
	Right string
}

type Trade struct {
	Pair
	Price    *big.Float
	Quantity *big.Float
}

type TradeProviderChan struct {
	Pair
	TradeChan chan Trade
}

type TradeProviderInterface interface {
	GetTradeProviderChan(pair Pair) (TradeProviderChan, error)
}

type TraderProvider struct {
	tradeProviderChan TradeProviderChan
}

func NewTradeProvider(lPair, rPair string, tradeFeeder TradeProviderInterface) (TraderProvider, error) {
	tFeeder := TraderProvider{}
	tPair := Pair{
		Left:  lPair,
		Right: rPair,
	}

	traderProvider, err := tradeFeeder.GetTradeProviderChan(tPair)
	if err != nil {
		return tFeeder, err
	}

	tFeeder.tradeProviderChan = traderProvider

	return tFeeder, nil
}

func (tf TraderProvider) Listen() <-chan Trade {
	return tf.tradeProviderChan.TradeChan
}
