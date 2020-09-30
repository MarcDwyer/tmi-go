package tc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sorcix/irc"
)

const (
	privmsg = "PRIVMSG"
	join    = "JOIN"
	notice  = "NOTICE"
	auth    = "001"
)

type eventFunc func(msg *irc.Message)

type Channel struct {
	channel string
	send    chan string
	part    chan string
	events  map[string]eventFunc
}

func NewChannel(channel string, send chan string, part chan string) *Channel {
	events := []string{"PRIVMSG", "NOTICE", "JOIN", "001"}
	evt := make(map[string]eventFunc)
	for _, event := range events {
		evt[event] = nil
	}
	fmt.Println(evt)
	return &Channel{
		channel: channel,
		send:    send,
		part:    part,
		events:  evt,
	}
}

func (c Channel) SendMsg(msg string) {
	cmd := fmt.Sprintf("PRIVMSG %s :%s", c.channel, msg)
	c.send <- cmd
}

// Part the channel
func (c *Channel) Part() {
	c.part <- c.channel
}

/*
AddListener takes in an event which must match an IRC command
*/
func (c Channel) AddListener(event string, function eventFunc) error {
	event = strings.ToUpper(event)
	if _, ok := c.events[event]; ok {
		c.events[event] = function
		return nil
	}
	return errors.New("Event was not found. Must match valid IRC command")
}

func (c Channel) HandleMsg(event string, msg *irc.Message) error {
	if function, ok := c.events[event]; ok {
		if function != nil {
			function(msg)
			return nil
		}
		return fmt.Errorf("Function not set: %s", event)
	}
	return fmt.Errorf("Could not find event: %s", event)
}
