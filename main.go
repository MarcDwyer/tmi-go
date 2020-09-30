package main

import (
	"fmt"
	"os"
	"time"
	tc "tmi-go/pkg"

	"github.com/sorcix/irc"
)

func main() {
	oauth := os.Getenv("oauth")
	// clientID := os.Getenv("client_id")
	// api := tc.TwitchAPI{
	// 	OAuth:    oauth,
	// 	ClientID: clientID,
	// }
	// followers, err := api.FetchV5Data("https://api.twitch.tv/kraken/streams/followed")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(followers)
	tc := tc.NewTwitchChat(&tc.TwitchCreds{
		Username: "roy_stang",
		OAuth:    oauth,
	})
	tc.Connect()
	name := "alisha12287"
	channel := tc.JoinChannel(name)

	channel.AddListener("PRIVMSG", func(ircMsg *irc.Message) {
		msg := fmt.Sprintf("%s: %s", ircMsg.User, ircMsg.Trailing)
		fmt.Println(msg)
	})
	time.Sleep(60 * time.Second)

}
