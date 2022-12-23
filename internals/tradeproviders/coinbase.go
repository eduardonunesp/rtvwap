package tradeproviders

import (
	"context"
	"math/big"
	"time"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Type of the subscription expected
	typeSubscribe = "subscribe"

	// Channel params expected on the WS subscription
	channelParams = "matches"
)

var (
	// Websocket address to consume
	wsURL = "wss://ws-feed.exchange.coinbase.com"

	// Minor error, happens when a message from WS is a non expected message from channelParams
	errNonExpectedMessage = errors.New("Non expected message")

	// Max number of errors in row from the subscription
	ErrCounterThreshold = 3
)

type (
	subReq struct {
		Type       string   `json:"type"`
		ProductIDs []string `json:"product_ids"`
		Channels   []string `json:"channels"`
	}

	matchRes struct {
		Type         string    `json:"type"`
		TradeID      int       `json:"trade_id"`
		Sequence     int       `json:"sequence"`
		MakerOrderID string    `json:"maker_order_id"`
		TakerOrderID string    `json:"taker_order_id"`
		Time         time.Time `json:"time"`
		ProductID    string    `json:"product_id"`
		Size         string    `json:"size"`
		Price        string    `json:"price"`
		Side         string    `json:"side"`
	}

	coinbaseProvider struct {
		ctx context.Context
	}
)

// NewCoinbaseProvider returns the TradeProvider from Coinbase
func NewCoinbaseProvider(context context.Context) internals.TradeProviderCreator {
	return coinbaseProvider{context}
}

// CreateTradeProvider will return the TradeProvider ready to use with a go channel ready to consume
func (c coinbaseProvider) CreateTradeProvider(pair internals.TradePair) (internals.TradeProvider, error) {
	tradeProvider := internals.NewTradeProvider()

	wsConn, err := newCoinbaseWS()
	if err != nil {
		return tradeProvider, errors.Wrap(err, "failed to connect to coinbase websocket feed")
	}

	if err := subscribeToMatchChannel(wsConn, pair.From+"-"+pair.To); err != nil {
		return tradeProvider, errors.Wrap(err, "failed to subscribe to coinbase websocket feed")
	}

	go func() {
		defer close(tradeProvider.TradeChan)

		var errCounter int

		for {
			select {
			case <-c.ctx.Done():
				close(tradeProvider.TradeChan)
				wsConn.Close()
			default:
				price, quantity, err := matchResponse(wsConn)
				if err != nil && errors.Is(err, errNonExpectedMessage) {
					continue
				} else if err != nil {
					if errCounter >= ErrCounterThreshold {
						break
					}
					continue
				}
				errCounter = 0

				// Only accept with price and quantity ok
				if price == nil && quantity == nil {
					continue
				}

				// New trade ready to push into go channel
				tradeProvider.TradeChan <- internals.NewTrade(pair, price, quantity)
			}

		}
	}()

	return tradeProvider, nil
}

// Creates new Coinbase websocket
func newCoinbaseWS() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Subscribes to the channel to get the match orders
func subscribeToMatchChannel(wsConn *websocket.Conn, productID string) error {
	return wsConn.WriteJSON(subReq{
		Type:       typeSubscribe,
		ProductIDs: []string{productID},
		Channels:   []string{channelParams},
	})
}

// Get the orders that matched/closed with the price and quantity
func matchResponse(wsConn *websocket.Conn) (*big.Float, *big.Float, error) {
	mRes := matchRes{}
	if err := wsConn.ReadJSON(&mRes); err != nil {
		return nil, nil, err
	}

	if mRes.Type != "match" {
		return nil, nil, errNonExpectedMessage
	}

	price, _, err := big.NewFloat(0).Parse(mRes.Price, 10)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to parse price from string: %s", mRes.Price)
	}

	quantity, _, err := big.NewFloat(0).Parse(mRes.Size, 10)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to parse quantity from string: %s", mRes.Price)
	}

	return price, quantity, nil
}
