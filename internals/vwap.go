package internals

import (
	"container/list"
	"context"
	"math/big"
)

// NewVWAP creates new Volume-Weighted Average Price from a TradeFeed
func NewVWAP(ctx context.Context, tradeChan TradeChannel) *VWAP {
	return &VWAP{
		ctx:          ctx,
		tradeChan:    tradeChan,
		tradeSamples: list.New(),
	}
}

// Calculate runs the VWAP calculation and send the resut to a chan of thep VWAP result
// The VWAP has the realtime calculations
func (vwap *VWAP) Calculate(vwapResultChan chan<- VWAPResult) {
	go func() {
	outloop:
		for {
			select {
			case <-vwap.ctx.Done():
				break outloop
			case trade, ok := <-vwap.tradeChan:
				if !ok {
					break outloop
				}
				vwap.tradeSamples.PushBack(trade)

				// Keep max of queue buffer size or 200 samples
				for vwap.tradeSamples.Len() > queueBufferSize {
					e := vwap.tradeSamples.Front()
					vwap.tradeSamples.Remove(e)
				}

				// PV = Σ(Price * Volume)
				sumPriceAndVolume := new(big.Float)
				for e := vwap.tradeSamples.Front(); e != nil; e = e.Next() {
					priceAndVolume := new(big.Float).Mul(e.Value.(Trade).Price, e.Value.(Trade).Quantity)
					sumPriceAndVolume = new(big.Float).Add(sumPriceAndVolume, priceAndVolume)
				}

				// SV = ΣVolume
				sumVolume := new(big.Float)
				for e := vwap.tradeSamples.Front(); e != nil; e = e.Next() {
					sumVolume = new(big.Float).Add(sumVolume, e.Value.(Trade).Quantity)
				}

				// VWP = PV / SV
				result := VWAPResult{
					Pair:      trade.TradePair,
					VWAPValue: new(big.Float).Quo(sumPriceAndVolume, sumVolume),
				}

				vwapResultChan <- result
			}
		}
	}()
}
