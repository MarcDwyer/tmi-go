package tc

import (
	"fmt"
)

type Channel struct {
	channel string
	send    chan string
	part    chan *Channel
}

func NewChannel(channel string, send chan string, part chan *Channel) *Channel {
	return &Channel{
		channel,
		send,
		part,
	}
}

func (c Channel) SendMsg(msg string) {
	cmd := fmt.Sprintf("PRIVMSG %s :%s", c.channel, msg)
	c.send <- cmd
}

// Part the channel
func (c *Channel) Part() {
	c.part <- c
}
