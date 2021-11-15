package tradeproviders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eduardonunesp/rtvwap/internals"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestCreateProvider(t *testing.T) {
	assert := assert.New(t)

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL = "ws" + strings.TrimPrefix(s.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var coinbaseProvider = NewCoinbaseProvider(ctx)
	pair := internals.NewTradePair("BTC", "USD")
	provider, err := coinbaseProvider.CreateTradeProvider(pair)
	assert.Nil(err)
	assert.NotNil(provider)

	go func() {
		time.Sleep(1 * time.Second)
		close(provider.TradeChan)
		cancel()
	}()

	<-provider.TradeChan
}
