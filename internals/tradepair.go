package internals

import "fmt"

// NewTradePair create new trade pair
func NewTradePair(from, to string) TradePair {
	return TradePair{from, to}
}

func (tp TradePair) String() string {
	return fmt.Sprintf("%s-%s", tp.From, tp.To)
}
