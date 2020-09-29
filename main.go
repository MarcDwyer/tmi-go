package main

import (
	"os"
	"time"
	tc "tmi-go/pkg"
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
	tc.ConnectAndRun()
	name := "alisha12287"
	channel := tc.JoinChannel(name)

	time.Sleep(15 * time.Second)
	channel.Part()
	time.Sleep(15 * time.Second)

}
