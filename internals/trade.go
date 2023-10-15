package internals

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

func NewTrade(pair TradePair, price, quantity *big.Float) Trade {
	return Trade{
		UUID:      uuid.New(),
		TradePair: pair,
		Price:     price,
		Quantity:  quantity,
	}
}

func (t Trade) String() string {
	return fmt.Sprintf("%s %s %s %s", t.UUID, t.TradePair, t.Price, t.Quantity)
}
