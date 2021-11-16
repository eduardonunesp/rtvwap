package tradeproviders

import (
	"context"
	"math/big"
	"time"

	"github.com/eduardonunesp/rtvwap/internals"
)

var valuesBTCUSD = [][]float64{
	{62246.89, 0.000161},
	{62246.9, 0.00086917},
	{62246.9, 0.01684996},
	{62246.9, 0.01574894},
	{62246.9, 0.0015504},
	{62246.9, 0.00314893},
	{62246.89, 0.00198374},
	{62246.9, 0.00234945},
	{62246.9, 0.00037572},
	{62246.89, 0.00290991},
	{62246.9, 0.00234945},
}

type mockProvider struct {
	ctx                context.Context
	errorTradeProvider error
}

func NewMockProvider(ctx context.Context) internals.TradeProviderCreator {
	return mockProvider{ctx, nil}
}

func (m mockProvider) CreateTradeProvider(pair internals.TradePair) (internals.TradeProvider, error) {
	tradeProvider := internals.TradeProvider{
		TradeChan: make(chan internals.Trade),
	}

	if m.errorTradeProvider != nil {
		return tradeProvider, m.errorTradeProvider
	}

	go func() {
		defer close(tradeProvider.TradeChan)

		count := 0
		loop := true

		for loop {
			select {
			case <-m.ctx.Done():
				loop = false
			default:
				newTrade := internals.Trade{
					TradePair: pair,
					Price:     big.NewFloat(valuesBTCUSD[count][0]),
					Quantity:  big.NewFloat(valuesBTCUSD[count][1]),
				}
				count++

				// New trade ready to push into go channel
				tradeProvider.TradeChan <- newTrade

				time.Sleep(100 * time.Millisecond)

				if count >= 10 {
					// return
					loop = false
				}
			}
		}
	}()

	return tradeProvider, nil
}
