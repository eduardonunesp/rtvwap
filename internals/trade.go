package internals

import (
	"math/big"

	"github.com/google/uuid"
)

// Trade represents a trade that matched/closed on the provider
type Trade struct {
	UUID uuid.UUID
	TradePair
	Price    *big.Float
	Quantity *big.Float
}

func NewTrade(pair TradePair, price, quantity *big.Float) Trade {
	return Trade{
		UUID:      uuid.New(),
		TradePair: pair,
		Price:     price,
		Quantity:  quantity,
	}
}
