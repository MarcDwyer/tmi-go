package main

import (
	"fmt"
	"log"
	"os"
	tc "tmi-go/pkg"
)

func main() {
	oauth := os.Getenv("oauth")
	tc := tc.NewTwitchChat("roy_stang", oauth)
	fmt.Println(tc)
	err := tc.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	for {
	}
}
