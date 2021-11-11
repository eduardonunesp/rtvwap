package tradeproviders

import (
	"log"
	"math/big"
	"time"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	wsURL = "wss://ws-feed.exchange.coinbase.com"
)

type (
	wsConn *websocket.Conn

	subReq struct {
		Type       string   `json:"type"`
		ProductIDs []string `json:"product_ids"`
		Channels   []string `json:"channels"`
	}

	unsubSubRes struct {
		Type     string `json:"type"`
		Channels []struct {
			Type       string   `json:"type"`
			ProductIDs []string `json:"product_ids"`
		} `json:"channels"`
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
)

type CoinbaseProvider struct{}

func NewCoinbaseProvider() internals.TradeProviderInterface {
	return CoinbaseProvider{}
}

func newCoinbaseWS() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func subscribeToMatchChannel(wsConn *websocket.Conn, productID string) error {
	return wsConn.WriteJSON(subReq{
		Type:       "subscribe",
		ProductIDs: []string{productID},
		Channels:   []string{"matches"},
	})
}

func checkSubscription(wsConn *websocket.Conn) error {
	var subRes unsubSubRes
	if err := wsConn.ReadJSON(&subRes); err != nil {
		log.Fatal(err)
	}

	return nil
}

func matchResponse(wsConn *websocket.Conn) (string, *big.Float, *big.Float, error) {
	mRes := matchRes{}
	if err := wsConn.ReadJSON(&mRes); err != nil {
		return "", nil, nil, err
	}

	log.Printf("ProductID: %s Size: %s Price: %s\n", mRes.ProductID, mRes.Size, mRes.Price)

	price, _, err := big.NewFloat(0).Parse(mRes.Price, 10)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, "failed to parse price from string: %s", mRes.Price)
	}

	quantity, _, err := big.NewFloat(0).Parse(mRes.Size, 10)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, "failed to parse quantity from string: %s", mRes.Price)
	}

	return mRes.ProductID, price, quantity, nil
}

func (CoinbaseProvider) GetTradeProviderChan(pair internals.Pair) (internals.TradeProviderChan, error) {
	tradeProvider := internals.TradeProviderChan{
		Pair:      pair,
		TradeChan: make(chan internals.Trade),
	}

	wsConn, err := newCoinbaseWS()
	if err != nil {
		return tradeProvider, errors.Wrap(err, "failed to connect to coinbase websocket feed")
	}

	if err := subscribeToMatchChannel(wsConn, pair.Left+"-"+pair.Right); err != nil {
		return tradeProvider, errors.Wrap(err, "failed to subscribe to coinbase websocket feed")
	}

	go func() {
		for {
			productID, price, quantity, err := matchResponse(wsConn)
			if err != nil {
				// Maybe is a good idea to have a limit of failures before return as error
				log.Println("warn: failed to match response", err)
				continue
			}

			// Just to make sure that we get the right pair
			if productID != pair.Left+"-"+pair.Right {
				log.Printf("warn: product %s id and pair do not match %s\n", productID, pair.Left+"-"+pair.Right)
				continue
			}

			newTrade := internals.Trade{
				Pair:     pair,
				Price:    price,
				Quantity: quantity,
			}

			tradeProvider.TradeChan <- newTrade
		}
	}()

	return tradeProvider, nil
}
