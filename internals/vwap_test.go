package internals_test

import (
	"context"
	"testing"
	"time"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/eduardonunesp/rtvwap/internals/tradeproviders"
	"github.com/stretchr/testify/assert"
)

func TestVWAPCalc(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	vwapResultChan := make(chan internals.VWAPResult)

	// Testing with mock provider with few samples
	tradeFeeder, err := internals.NewTradeFeedWithPair(internals.NewTradePair("BTC", "USD"), tradeproviders.NewMockProvider(ctx))
	assert.Nil(err)

	internals.NewVWAP(ctx, tradeFeeder).Calculate(vwapResultChan)

	VWPALastResult := internals.VWAPResult{}
	go func() {
		for i := range vwapResultChan {
			VWPALastResult = i
		}
	}()

	time.Sleep(1 * time.Second)

	assert.Equal(VWPALastResult.Pair, internals.NewTradePair("BTC", "USD"))
	assert.Equal(VWPALastResult.VWAPValue.String(), "62246.8989")
}
