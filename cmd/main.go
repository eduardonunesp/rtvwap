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
	provider internals.TradeProviderCreator
}

func createSignalChannel() chan os.Signal {
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, syscall.SIGINT, syscall.SIGTERM)
	return sigInt
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	vwapResultChan := internals.CreateResultChan()

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
		tradeFeed, err := internals.NewTradeFeed(tradeGroup.pair, tradeGroup.provider)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Run the calculations for each for the current trade feed
		internals.NewVWAP(ctx, tradeFeed).Calculate(vwapResultChan)
	}

	for {
		select {
		case result := <-vwapResultChan:
			fmt.Printf("VWAP RESULT %s %f\n", result.Pair.From+"-"+result.Pair.To, result.VWAPValue)
		case <-createSignalChannel():
			cancel()
			fmt.Println(" SIGINT: Closing the program")
			os.Exit(0)
		}
	}
}
