package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"

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
	signal.Notify(sigInt, os.Interrupt)

	// Create trades
	tradeFeeders := []tradePairs{
		{
			pair:     internals.NewTradePair("BTC", "USD"),
			provider: tradeproviders.NewCoinbaseProvider(),
		},
		{
			pair:     internals.NewTradePair("ETH", "USD"),
			provider: tradeproviders.NewCoinbaseProvider(),
		},
		{
			pair:     internals.NewTradePair("ETH", "BTC"),
			provider: tradeproviders.NewCoinbaseProvider(),
		},
	}

	// Run calculation for each trade pair
	for _, tradeFeed := range tradeFeeders {
		tradeFeeder, err := internals.NewTradeFeedWithPair(tradeFeed.pair, tradeproviders.NewCoinbaseProvider())
		if err != nil {
			log.Fatal(err)
		}

		internals.NewVWAP(tradeFeeder).Calculate()
	}

	select {
	case <-sigInt:
		fmt.Println(" SIGINT: Closing the program")
		os.Exit(0)
	}
}
