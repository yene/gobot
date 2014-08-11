package main

import (
	"./twitch"
	"./wisdom"
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
)

var favorites []string

func main() {
	favorites = twitch.FavoriteDota2Streams()
	channel := "#test" //#r/dota2"
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
		case "!matches", "!m":
			con.Privmsg(channel, "no matches")
		case "!help":
			con.Privmsg(channel, "no help")
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
	go func() {
		for {
			time.Sleep(time.Second * 30)
			newFavorites := twitch.FavoriteDota2Streams()
			for _, g := range newFavorites {
				if !inisdeFavorites(g) {
					con.Privmsg(channel, g+" started streaming.")
				}
			}
			favorites = newFavorites
		}
	}()
	con.Loop()
}

func inisdeFavorites(a string) bool {
	for _, g := range favorites {
		if g == a {
			return true
		}
	}
	return false
}
