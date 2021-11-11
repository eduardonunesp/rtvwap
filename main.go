package main

import (
	"log"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/eduardonunesp/rtvwap/internals/tradeproviders"
)

func main() {
	tradeFeeder, err := internals.NewTradeProvider("BTC", "USD", tradeproviders.NewCoinbaseProvider())
	if err != nil {
		log.Fatal(err)
	}

	for trade := range tradeFeeder.Listen() {
		log.Println(trade)
	}
}
