# rtvwap

The project rtvwap is a sample project for calculate the volume-weighted average price for trading pairs, it's using the `coinbase` as provider but it's extensible for different providers if needed

## Sample

The sample code supports the following pairs:
- BTC-USD
- ETH-USD
- ETH-BTC

## Description

The `rtvwap` is composed of some components/packages which are the `tradeproviders` and `vwap (calculation)`. Also, it needs the `TradePair` struct which
are responsible for the definition of trading to compute the average. The `TradeProvider` struct is the wrapper for the `TradeChan` which streams information of each 
match `Trade` with the trading pairs and the Price and Quantity (important variables for the VWAP calculation).

The interface `TradeProviderCreator` has one single function `CreateTradeProvider` which all `tradeproviders` must agree in order to be considered an trading provider,
the function has the following signature `CreateTradeProvider(pair TradePair) (TradeProvider, error)`.

In order to calculate the `TradeProvider` you should use the `TradeFeeder` which will receive any `TradeProvider` that respects the interface `TraderProviderCreator` and 
and return the `TradeFeed` object that should be passed to the `VWAP` calculation

And finally the `VWAP` calculation is created with the `Context` and `TradeFeed` (that should contains the ready `TradeChan`), using the function `NewVWAP(ctx context.Context, tradeFeed TradeFeed) VWAP`, 
next you should call the method `Calculate` of the struct `VWAP` and passing the result channel using the struct `VWAPResult` to get in the realtime the calculations for the `VWAP` pair selected.

## CoinbaseProvider

The coinbase provider accepts a trading pair to subscribe to coinbase socket in order to start listen the channel data for trading information, here we're using the 
 subscription for `channel matches` to get the trading matches that happen on the plataform. It's important to note that the channels can have missing sequences of matches
which for the realtime propose this implementation don't try to track and recover, later the information provider can be improved for retrieve all missing sequences if needed, and 
because many sequences are missed can be a waste of "realtime" to try to always retrieve all missed sequences.

Because all providers follows the `TradeProviderCreator` interface should be straightforward to add others providers in order to improve the market values using many
different plataforms like coinmarketcap which incorporate about 200 exchanges as source of trading information

## VWAP

Finally the [`VWAP`](https://en.wikipedia.org/wiki/Volume-weighted_average_price) which is implemented as a consume of the `TradeFeed`, which does the calculation and returns the data stream over the `chan VWAPResult` that can be used as the stream of real time data for WVAP results

> The VWAP calculating the sum of the latest 200 samples it can be updated by changing the var `internals.QueueBufferSize`

## Running 

```bash
go get -u github.com/eduardonunesp/rtvwap/...
```

> Should install into your GOPATH/bin the binary `rtvwap` 

## Testing

```bash
go test ./...
```

## Improvements

There're room for many improvements specially for calculation which can be transformed into some interface for serve as other calculations or indexes. Also the tests
can be improved adding more mocked that and improving the `Mock` provider which is already used for testing the `TradeFeed` and `VWAP` calculations. With a proper backend
the data can be stored for more analysis as well. As mentioned before the sequece numbers of matches sometimes have some gaps this can be improved if some store or 
observer can watch for gaps on the sequences and retrieve the trading information for more complete information.






## VWAP 

The calcula
