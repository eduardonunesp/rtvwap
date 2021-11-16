package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/eduardonunesp/rtvwap/internals/tradeproviders"
)

// Trade pairs to create
type tradePairs struct {
	pair     internals.TradePair
	provider internals.TradeProviderCreator
}

func main() {
	// interruption signal
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	vwapResultChan := make(chan internals.VWAPResult)

	// Create trades
	tradeFeeders := []tradePairs{
		{
			pair:     internals.NewTradePair("BTC", "USD"),
			provider: tradeproviders.NewCoinbaseProvider(ctx),
		},
		{
			pair:     internals.NewTradePair("ETH", "USD"),
			provider: tradeproviders.NewCoinbaseProvider(ctx),
		},
		{
			pair:     internals.NewTradePair("ETH", "BTC"),
			provider: tradeproviders.NewCoinbaseProvider(ctx),
		},
	}

	// Run calculation for each trade pair
	for _, tradeFeed := range tradeFeeders {
		tradeFeeder, err := internals.NewTradeFeedWithPair(tradeFeed.pair, tradeFeed.provider)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Run the calculations for each for the current trade feeder
		internals.NewVWAP(ctx, tradeFeeder).Calculate(vwapResultChan)
	}

	for {
		select {
		case result := <-vwapResultChan:
			fmt.Printf("VWAP RESULT %s %f\n", result.Pair.Left+"-"+result.Pair.Right, result.VWAPValue)
		case <-sigInt:
			cancel()
			fmt.Println(" SIGINT: Closing the program")
			os.Exit(0)
		}
	}
}
