package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	tc "tmi-go/pkg"

	"github.com/sorcix/irc"
)

func main() {
	oauth := os.Getenv("oauth")
	clientID := os.Getenv("client_id")
	api := tc.TwitchAPI{
		OAuth:    oauth,
		ClientID: clientID,
	}
	followers, err := api.FetchV5Data("https://api.twitch.tv/kraken/streams/followed")
	if err != nil {
		panic(err)
	}
	// fmt.Println(followers)
	tc := tc.NewTwitchChat("roy_stang", oauth)
	tc.Connect()

	for _, stream := range followers.Streams {
		name := stream.Channel.Name
		c := tc.JoinChannel(name)
		c.SendMsg("is this stream vegan")
		c.AddListener("PRIVMSG", func(ircMsg *irc.Message) {
			msg := strings.ToLower(ircMsg.Trailing)
			dm := strings.Contains(msg, tc.Username)
			if dm == true {
				prnt := fmt.Sprintf("Channel: %s. User: %s. Says: %s", ircMsg.Params[0], ircMsg.User, ircMsg.Trailing)
				fmt.Println(prnt)
			}
		})
	}

	time.Sleep(60 * time.Second)

}
