package tc

import (
	"errors"
	"fmt"

	"github.com/sorcix/irc"

	"github.com/gorilla/websocket"
)

type TwitchChat struct {
	username string
	oauth    string
	send     chan string
	ws       *websocket.Conn
}

func NewTwitchChat(username string, oauth string) *TwitchChat {
	return &TwitchChat{
		username: username,
		oauth:    oauth,
		send:     make(chan *string),
	}
}

func (tc TwitchChat) Connect() error {
	tcURL := "ws://irc-ws.chat.twitch.tv:80"
	c, _, err := websocket.DefaultDialer.Dial(tcURL, nil)
	if err != nil {
		return err
	}
	tc.ws = c
	go tc.run()
	go func() {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println("err")
		}
		ircMsg := irc.ParseMessage(string(msg[:]))
		fmt.Println(ircMsg.Command)
	}()
	tc.send <- fmt.Sprintf("PASS oauth:%s", tc.oauth)
	tc.send <- fmt.Sprintf("NICK %s", tc.username)
	return nil
}

func (tc TwitchChat) run() {
	for {
		select {
		case msg := <-tc.send:
			fmt.Println(msg)
			err := tc.ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error sending: %s", msg)
			}
		}
	}
}

// func (tc TwitchChat) send(msg string) error {
// 	if tc.ws == nil {
// 		return errors.New("Must connect before attempting to send message")
// 	}
// 	err := tc.ws.WriteMessage(websocket.TextMessage, []byte(msg))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (tc TwitchChat) JoinChannel(channel string) (*Channel, error) {
	fc := string(channel[0])
	if fc != "#" {
		channel = fmt.Sprintf("#%s", channel)
	}
	if tc.ws == nil {
		return nil, errors.New("Must connect before attempting to send message")
	}
	return &Channel{
			ws:      tc.ws,
			channel: channel,
		},
		nil
}
