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

type expectedTradeFeed struct {
	pair     internals.TradePair
	provider internals.TradeProviderCreator
}

func main() {
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)

	tradeFeeders := []expectedTradeFeed{
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

	for _, tradeFeed := range tradeFeeders {
		tradeFeeder, err := internals.NewTradeFeedWithPair(tradeFeed.pair, tradeproviders.NewCoinbaseProvider())
		if err != nil {
			log.Fatal(err)
		}

		vwap := internals.NewVWAP(tradeFeeder)

		vwap.Calculate()
	}

	select {
	case <-sigInt:
		fmt.Println(" SIGINT: Closing the program")
		os.Exit(0)
	}
}
