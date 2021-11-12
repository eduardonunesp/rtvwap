package internals

import (
	"math/big"
)

var PreAllocatedSize = 200
var QueueBufferSize = 200

// VWAP represesnts the Volume-Weighted Average Price calculation
type VWAP struct {
	tradeChan             <-chan Trade
	samplesChan           chan Trade
	totalPriceAndQuantity *big.Float
	totalQuantity         *big.Float
}

// NewVWAP creates new Volume-Weighted Average Price from a TradeFeed
func NewVWAP(tradeFeed TradeFeed) VWAP {
	vwap := VWAP{
		tradeFeed.TradeChan(),
		make(chan Trade, QueueBufferSize),
		new(big.Float),
		new(big.Float),
	}

	// Pre allocate some blank values to accelarate calculations results
	for i := 0; i < PreAllocatedSize; i++ {
		vwap.samplesChan <- Trade{
			TradePair: TradePair{"", ""},
			Price:     new(big.Float),
			Quantity:  new(big.Float),
		}
	}

	// Run the calculation process on start
	go func() {
		// As new trade arrive should push into queueChan for VWAP calculation
		for trade := range vwap.tradeChan {
			vwap.samplesChan <- trade
		}
	}()

	return vwap
}

func (vwap VWAP) Calculate() {
	go func() {
		for {
			select {
			case trade := <-vwap.samplesChan:
				_ = trade
				// oldPrice, oldQuantity := vwap.queue.AddNew(trade)

				// mulNewPriceAndQuantity := NewEmptyBigFloat()
				// mulNewPriceAndQuantity.Mul(trade.Price, trade.Quantity)

				// if oldPrice != nil && oldQuantity != nil {
				// 	mulOldPriceAndQuantity := NewEmptyBigFloat()
				// 	mulOldPriceAndQuantity.Mul(oldPrice, oldQuantity)

				// 	vwap.totalPriceAndQuantity.Sub(vwap.totalPriceAndQuantity, mulOldPriceAndQuantity)
				// 	vwap.totalPriceAndQuantity.Add(vwap.totalPriceAndQuantity, mulNewPriceAndQuantity)

				// 	vwap.totalQuantity.Sub(vwap.totalQuantity, oldQuantity).Add(vwap.totalQuantity, trade.Quantity)

				// } else {
				// 	vwap.totalPriceAndQuantity.Add(vwap.totalPriceAndQuantity, mulNewPriceAndQuantity)
				// 	vwap.totalQuantity.Add(vwap.totalQuantity, trade.Quantity)
				// }

				// vwapResult := NewEmptyBigFloat()
				// vwapResult.Quo(vwapResult.totalPriceAndQuantity, vwapResult.totalQuantity)

				// vwapResult.logger.Printf("VWAP for pair %s is %v", vc.pair, vwapResult)
			}
		}
	}()
}
