package internals

import (
	"log"
	"math/big"
)

// Max number of samples to calculate
var QueueBufferSize = 200

// VWAP represesnts the Volume-Weighted Average Price calculation
type VWAP struct {
	tradeChan             <-chan Trade
	tradeSamples          []Trade
	totalPriceAndQuantity *big.Float
	totalQuantity         *big.Float
}

// NewVWAP creates new Volume-Weighted Average Price from a TradeFeed
func NewVWAP(tradeFeed TradeFeed) VWAP {
	vwap := VWAP{
		tradeChan:             tradeFeed.TradeChan(),
		totalPriceAndQuantity: new(big.Float),
		totalQuantity:         new(big.Float),
	}

	return vwap
}

func (vwap VWAP) Calculate() {
	go func() {
		for {
			select {
			case trade := <-vwap.tradeChan:
				var prevTrade Trade
				newTrade := trade

				for len(vwap.tradeSamples) > 0 {
					prevTrade = vwap.tradeSamples[0]
					vwap.tradeSamples = vwap.tradeSamples[1:]
				}
				vwap.tradeSamples = append(vwap.tradeSamples, newTrade)

				mulNewPriceAndQuantity := new(big.Float)
				mulNewPriceAndQuantity.Mul(newTrade.Price, newTrade.Quantity)

				if prevTrade.Price != nil && prevTrade.Quantity != nil {
					mulOldPriceAndQuantity := new(big.Float)
					mulOldPriceAndQuantity.Mul(prevTrade.Price, prevTrade.Quantity)

					vwap.totalPriceAndQuantity.Sub(vwap.totalPriceAndQuantity, mulOldPriceAndQuantity)
					vwap.totalPriceAndQuantity.Add(vwap.totalPriceAndQuantity, mulNewPriceAndQuantity)

					vwap.totalQuantity.Sub(vwap.totalQuantity, prevTrade.Quantity).Add(vwap.totalQuantity, newTrade.Quantity)
				} else {
					vwap.totalPriceAndQuantity.Add(vwap.totalPriceAndQuantity, mulNewPriceAndQuantity)
					vwap.totalQuantity.Add(vwap.totalQuantity, newTrade.Quantity)
				}

				vwapResult := new(big.Float)
				vwapResult.Quo(vwap.totalPriceAndQuantity, vwap.totalQuantity)

				log.Printf("VWAP for pair %s is %v", newTrade.Left+"-"+newTrade.Right, vwapResult)
			}
		}
	}()
}
