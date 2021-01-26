package cmd

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

func checkIfWsAvailable(websocketUrl *url.URL) error {
	// check if ws available on given port
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 2 * time.Second
	c, _, err := dialer.Dial(websocketUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("unable to connect to ws endpoint, because of %v", err)
	}
	err = c.Close()
	if err != nil {
		return err
	}
	return err
}