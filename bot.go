package main

import (
	"fmt"
	"strings"

	"./matches"
	"./tips"
	"./twitch"
	"github.com/thoj/go-ircevent"
)

func main() {
	channel := "#dota2"

	con := irc.IRC("Tresdin", "Tresdin")
	err := con.Connect("irc.euirc.net:6667")
	if err != nil {
		fmt.Println("Failed connecting")
		return
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		switch e.Message() {
		case "!major", "!m":
			con.Privmsg(channel, twitch.Major())
		case "!streams":
			con.Privmsg(channel, strings.Join(twitch.FilteredDota2Streams(), " - "))
		case "!tip":
			con.Privmsg(channel, tips.RandomTip())
		case "!credit":
			con.Privmsg(channel, "https://github.com/yene/gobot")
		case "!score", "!s":
			con.Privmsg(channel, matches.MajorScore())
		}
	})

	con.Loop()
}
