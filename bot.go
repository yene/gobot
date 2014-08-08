package main

import (
	"./streamer"
	"fmt"
	"github.com/thoj/go-ircevent"
)

func main() {
	channel := "#r/dota2" //#r/dota2"
	con := irc.IRC("Lina", "Lina")
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
		//case "!matches", "!m":
		//con.Privmsg(channel, "no matches")
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
		case "!relax":
			con.Privmsg(channel, "People forget that they also win games because the enemies are feeding and are shit.")
			// con.Privmsg(channel, "MMR is not bugged or rigged or anything. Get better, you'll end up with better MMR.")
			// con.Privmsg(channel, "Stop crying that "I ended up with 2500 because my team mates suck". No. You ended up 2,5k because you suck. ")
		}
		case "!joke":
			con.Privmsg(channel, "my MMR")
		}
	})
	con.Loop()
}
