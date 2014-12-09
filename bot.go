package main

import (
	"./tips"
	"./twitch"
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

func main() {
	channel := "#test2"
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
		case "!favorites", "!f":
			//con.Privmsg(channel, strings.Join(twitch.FavoriteDota2Streams(), " - "))
		case "!streams", "!s":
			con.Privmsg(channel, strings.Join(twitch.FilteredDota2Streams(), " - "))
		case "!tip":
			con.Privmsg(channel, tips.RandomTip())
		case "!credit":
			con.Privmsg(channel, "https://github.com/yene/gobot")
		}
	})

	streams := make(chan []twitch.Channel)

	go twitch.UpdateStreams(streams)

	go twitch.WatchFavorites(streams, func(m string) {
		con.Privmsg(channel, m)
	})
	/*
		go twitch.WatchTournaments(func(m string) {
			con.Privmsg(channel, m)
		})

		go twitch.WatchAll(func(m string) {
			con.Privmsg(channel, m)
		})
	*/

	con.Loop()
}
