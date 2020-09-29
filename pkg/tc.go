package tc

import (
	"fmt"

	"strings"

	"github.com/gorilla/websocket"
	"github.com/sorcix/irc"
)

type TwitchChat struct {
	credentials *TwitchCreds
	send        chan string
	part        chan *Channel
	channels    map[*Channel]bool
}
type TwitchCreds struct {
	Username string
	OAuth    string
}

func NewTwitchChat(credentials *TwitchCreds) *TwitchChat {
	return &TwitchChat{
		credentials: credentials,
		send:        make(chan string),
		part:        make(chan *Channel),
		channels:    make(map[*Channel]bool),
	}
}

type isAuth struct {
	succ   bool
	rawMsg string
}

func (tc TwitchChat) ConnectAndRun() string {
	tcURL := "ws://irc-ws.chat.twitch.tv:80"
	c, _, err := websocket.DefaultDialer.Dial(tcURL, nil)
	if err != nil {
		panic(err)
	}

	auth := make(chan isAuth)
	go tc.run(c)
	go func() {
		creds := tc.credentials
		tc.send <- fmt.Sprintf("PASS oauth:%s", creds.OAuth)
		tc.send <- fmt.Sprintf("NICK %s", creds.Username)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				continue
			}
			raw := string(msg[:])
			ircMsg := irc.ParseMessage(raw)
			fmt.Println(ircMsg)
			switch ircMsg.Command {
			case "001":
				fmt.Println("succ")
				pair := isAuth{
					rawMsg: raw,
					succ:   true,
				}
				auth <- pair
				break
			case "NOTICE":
				rawUnder := strings.ToLower(raw)
				isFail := strings.Contains(rawUnder, "failed")
				if isFail == true {
					auth <- isAuth{
						rawMsg: raw,
						succ:   false,
					}
				}
			case "PRIVMSG":
				fmt.Println(ircMsg.Trailing)
				break
			default:
				fmt.Println(ircMsg.Params)
			}
		}
	}()

	succ := <-auth
	if succ.succ == true {
		return succ.rawMsg
	}
	panic(succ.rawMsg)
}

func (tc TwitchChat) run(ws *websocket.Conn) {
	for {
		select {
		case msg := <-tc.send:
			fmt.Println(msg)
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Printf("Error sending: %s", msg)
			}
		case channel := <-tc.part:
			fmt.Println(channel)
			if _, ok := tc.channels[channel]; ok {
				fmt.Println("ok")
				tc.send <- fmt.Sprintf("PART %s", channel.channel)
				delete(tc.channels, channel)
			}
		}
	}
}

func (tc TwitchChat) JoinChannel(channel string) *Channel {
	fc := string(channel[0])
	if fc != "#" {
		channel = fmt.Sprintf("#%s", channel)
	}
	tc.send <- fmt.Sprintf("JOIN %s", channel)
	newchan := &Channel{
		channel: channel,
		send:    tc.send,
		part:    tc.part,
	}
	tc.channels[newchan] = true
	return newchan
}
