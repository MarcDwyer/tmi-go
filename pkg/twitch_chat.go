package tc

import (
	"fmt"

	"strings"

	"github.com/gorilla/websocket"
	"github.com/sorcix/irc"
)

type TwitchChat struct {
	Username string
	OAuth    string
	send     chan string
	part     chan string
	channels map[string]*Channel
}

// Trailing is the message, USER is the user origin
//
func NewTwitchChat(username string, oauth string) *TwitchChat {
	return &TwitchChat{
		Username: username,
		OAuth:    oauth,
		send:     make(chan string),
		part:     make(chan string),
		channels: make(map[string]*Channel),
	}
}

type isAuth struct {
	succ   bool
	rawMsg string
}

//Connect connects to Twitch's IRC server
// will panic if it fails
func (tc TwitchChat) Connect() string {
	tcURL := "ws://irc-ws.chat.twitch.tv:80"
	c, _, err := websocket.DefaultDialer.Dial(tcURL, nil)
	if err != nil {
		panic(err)
	}
	auth := make(chan isAuth)
	go tc.run(c)
	go func() {
		tc.send <- fmt.Sprintf("PASS oauth:%s", tc.OAuth)
		tc.send <- fmt.Sprintf("NICK %s", tc.Username)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				continue
			}
			raw := string(msg[:])
			ircMsg := irc.ParseMessage(raw)
			// fmt.Println(ircMsg.Params)
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
				break
			default:
				{
					name := ircMsg.Params[0]
					if channel, ok := tc.channels[name]; ok {
						err := channel.HandleMsg(ircMsg.Command, ircMsg)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
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
		case chanName := <-tc.part:
			if _, ok := tc.channels[chanName]; ok {
				tc.send <- fmt.Sprintf("PART %s", chanName)
				delete(tc.channels, chanName)
			}
		}
	}
}

//JoinChannel joins a twitch channel example: "ninja" or "xqcow"
func (tc TwitchChat) JoinChannel(channel string) *Channel {
	channel = strings.ToLower(channel)
	fc := string(channel[0])
	if fc != "#" {
		channel = fmt.Sprintf("#%s", channel)
	}
	tc.send <- fmt.Sprintf("JOIN %s", channel)
	newchan := NewChannel(channel, tc.send, tc.part)
	tc.channels[channel] = newchan
	return newchan
}
