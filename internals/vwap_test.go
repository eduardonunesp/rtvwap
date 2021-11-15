package internals

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	ctx                context.Context
	errorTradeProvider error
}

func newMockProvider(ctx context.Context) TradeProviderCreator {
	return mockProvider{ctx, nil}
}

func (m mockProvider) CreateTradeProvider(pair TradePair) (TradeProvider, error) {
	tradeProvider := TradeProvider{
		TradeChan: make(chan Trade),
	}

	// if m.errorTradeProvider != nil {
	// 	return tradeProvider, m.errorTradeProvider
	// }

	// count := 0

	go func() {
		// 	defer close(tradeProvider.TradeChan)

		// 	for {
		// 		select {
		// 		case <-m.ctx.Done():
		// 			break
		// 		default:
		// 			count++

		// New trade ready to push into go channel
		newTrade := Trade{
			TradePair: pair,
			Price:     big.NewFloat(65.0),
			Quantity:  big.NewFloat(0.0001),
		}

		log.Println("SD")
		tradeProvider.TradeChan <- newTrade

		time.Sleep(1 * time.Millisecond)

		// 	if count > 10 {
		// 		break
		// 	}
		// }
		// }
	}()

	return tradeProvider, nil
}

func TestVWAPCalc(t *testing.T) {
	assert := assert.New(t)
	// _ = assert
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	// vwapResultChan := make(chan VWAPResult)

	tradeFeeder, err := NewTradeFeedWithPair(NewTradePair("BTC", "USD"), newMockProvider(ctx))
	assert.Nil(err)
	_ = tradeFeeder

	time.Sleep(1 * time.Second)

	for i := range tradeFeeder.TradeChan() {
		log.Println(i)
	}

	// _ = vwapResultChan
	// _ = tradeFeeder
	// _ = cancel

	// NewVWAP(ctx, tradeFeeder).Calculate(vwapResultChan)

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	cancel()
	// }()

	// for v := range vwapResultChan {
	// 	log.Println(v)
	// }
}
