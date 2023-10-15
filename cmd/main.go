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

// Trade group with pair and provider
type tradeGroup struct {
	pair     internals.TradePair
	provider internals.TradeProvider
}

func createSignalChannel() chan os.Signal {
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, syscall.SIGINT, syscall.SIGTERM)
	return sigInt
}

func createResultChan() chan internals.VWAPResult {
	return make(chan internals.VWAPResult)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	vwapResultChan := createResultChan()

	// Create trade feeders
	tradeGroups := []tradeGroup{
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
	for _, tradeGroup := range tradeGroups {
		tradeChan, err := tradeGroup.provider.GetTradeChannel(tradeGroup.pair)
		if err != nil {
			panic(err)
		}

		// Run the calculations for each for the current trade feed
		internals.NewVWAP(ctx, tradeChan).Calculate(vwapResultChan)
	}

	for {
		select {
		case result := <-vwapResultChan:
			fmt.Printf("VWAP RESULT %s %f\n", result.Pair.String(), result.VWAPValue)
		case <-createSignalChannel():
			cancel()
			fmt.Println(" SIGINT: Closing the program")
			close(vwapResultChan)
			os.Exit(0)
		}
	}
}
