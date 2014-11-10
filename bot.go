package main

import (
	"./twitch"
	"./wisdom"
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

func main() {
	channel := "#r/test"
	con := irc.IRC("Tresdin", "Tresdin")
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
		case "!matches", "!m", "!t", "!tournament":
			s := strings.Join(twitch.TournamentStreams(), " - ")
			if len(s) == 0 {
				con.Privmsg(channel, "No tournaments live.")
			} else {
				con.Privmsg(channel, s)
			}
		case "!all", "!a":
			con.Privmsg(channel, strings.Join(twitch.Dota2Streams(), " - "))
		case "!help", "!h":
			con.Privmsg(channel, "Use !s for filtered streams, !a for all streams.")
		case "!scores":
			con.Privmsg(channel, "no scores")
		case "!favorites", "!f":
			con.Privmsg(channel, strings.Join(twitch.FavoriteDota2Streams(), " - "))
		case "!streams", "!s":
			con.Privmsg(channel, strings.Join(twitch.TopDota2Streams(), " - "))
		case "!relax", "!wisdom":
			con.Privmsg(channel, wisdom.RandomWisdom())
		case "!joke":
			con.Privmsg(channel, "my MMR")
		}
	})

	go twitch.WatchFavorites(func(m string) {
		con.Privmsg(channel, m)
	})
	go twitch.WatchTournaments(func(m string) {
		con.Privmsg(channel, m)
	})

	con.Loop()
}
