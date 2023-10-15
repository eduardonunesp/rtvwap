package internals

import (
	"container/list"
	"context"
	"math/big"

	"github.com/google/uuid"
)

// Max number of samples to calculate
const queueBufferSize = 200

// TradeChannel is the channel is used to receive trades from the provider
type TradeChannel chan Trade

// TradePair is the representation of the trading pair left and right pairs
// For instance ETH-USD, BTC-USD, USD-ETH
type TradePair struct {
	From string
	To   string
}

// Trade represents a trade that matched/closed on the provider
type Trade struct {
	UUID      uuid.UUID
	TradePair TradePair
	Price     *big.Float
	Quantity  *big.Float
}

// VWAP represesnts the Volume-Weighted Average Price calculation
type VWAP struct {
	ctx          context.Context
	tradeChan    TradeChannel
	tradeSamples *list.List
}

// VWAPResult serves as a container for the result of a VWAP calculation
type VWAPResult struct {
	Pair      TradePair
	VWAPValue *big.Float
}

// TradeProvider is the core interface that all trade providers must implement
type TradeProvider interface {
	GetTradeChannel(pair TradePair) (TradeChannel, error)
}
