package internals

import (
	"container/list"
	"context"
	"math/big"
)

// Max number of samples to calculate
var QueueBufferSize = 200

// VWAP represesnts the Volume-Weighted Average Price calculation
type VWAP struct {
	ctx          context.Context
	tradeChan    <-chan Trade
	tradeSamples *list.List
}

type VWAPResult struct {
	Pair       TradePair
	VWAPResult *big.Float
}

// NewVWAP creates new Volume-Weighted Average Price from a TradeFeed
func NewVWAP(ctx context.Context, tradeFeed TradeFeed) VWAP {
	vwap := VWAP{
		ctx:          ctx,
		tradeChan:    tradeFeed.TradeChan(),
		tradeSamples: list.New(),
	}

	return vwap
}

func (vwap VWAP) Calculate(vwapResultChan chan<- VWAPResult) {
	go func() {
		for {
			select {
			case <-vwap.ctx.Done():
				return
			case trade := <-vwap.tradeChan:
				vwap.tradeSamples.PushBack(trade)

				for vwap.tradeSamples.Len() > QueueBufferSize {
					e := vwap.tradeSamples.Front()
					vwap.tradeSamples.Remove(e)
				}

				sumPriceAndVolume := new(big.Float)
				for e := vwap.tradeSamples.Front(); e != nil; e = e.Next() {
					priceAndVolume := new(big.Float).Mul(e.Value.(Trade).Price, e.Value.(Trade).Quantity)
					sumPriceAndVolume = new(big.Float).Add(sumPriceAndVolume, priceAndVolume)
				}

				sumVolume := new(big.Float)
				for e := vwap.tradeSamples.Front(); e != nil; e = e.Next() {
					sumVolume = new(big.Float).Add(sumVolume, e.Value.(Trade).Quantity)
				}

				vwapResultChan <- VWAPResult{
					Pair:       trade.TradePair,
					VWAPResult: new(big.Float).Quo(sumPriceAndVolume, sumVolume),
				}
			}
		}
	}()
}
