package internals_test

import (
	"context"
	"testing"
	"time"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/eduardonunesp/rtvwap/internals/tradeproviders"
	"github.com/stretchr/testify/require"
)

func TestVWAPCalc(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	vwapResultChan := make(chan internals.VWAPResult)

	// Testing with mock provider with few samples
	provider := tradeproviders.NewMockProvider(ctx)
	tradeChan, err := provider.GetTradeChannel(internals.NewTradePair("BTC", "USD"))
	require.NoError(err)
	require.NotNil(tradeChan)

	internals.NewVWAP(ctx, tradeChan).Calculate(vwapResultChan)

	VWPALastResult := internals.VWAPResult{}
	go func() {
		for i := range vwapResultChan {
			VWPALastResult = i
		}
	}()

	time.Sleep(1 * time.Second)

	require.Equal(VWPALastResult.Pair, internals.NewTradePair("BTC", "USD"))
	require.Equal(VWPALastResult.VWAPValue.String(), "62246.8989")
}
