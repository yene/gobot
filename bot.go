package main

import (
	"./streamer"
	"fmt"
	"github.com/thoj/go-ircevent"
)

func main() {
	channel := "#test" //#r/dota2"
	con := irc.IRC("7kMMR", "7kMMR")
	err := con.Connect("irc.quakenet.org:6667")
	if err != nil {
		fmt.Println("Failed connecting")
		return
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {

		switch e.Message() {
		case "!matches", "!m":
			con.Privmsg(channel, "no matches")
		case "!help":
			con.Privmsg(channel, "no help")
		case "!scores":
			con.Privmsg(channel, "no scores")
		case "!favorites", "!f":
			for _, g := range streamer.FavoriteDota2Streams() {
				con.Privmsg(channel, g)
			}
		case "!streams", "!s":
			for _, g := range streamer.TopDota2Streams() {
				con.Privmsg(channel, g)
			}
		case "!joke":
			con.Privmsg(channel, "my mmr")
		}
	})
	con.Loop()
}
