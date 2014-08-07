package streamer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func FavoriteDota2Streams() []string {
	favorites := favoriteStreams()
	concatenated := strings.Replace(favorites, "\n", ",", -1)
	requestURL := "https://api.twitch.tv/kraken/streams?game=Dota+2&limit=10&channel=" + concatenated
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	streams, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dat JSONResult
	if err := json.Unmarshal(streams, &dat); err != nil {
		panic(err)
	}

	sslice := make([]string, 0)
	for _, g := range dat.Streams {
		s := fmt.Sprintf("%s (%d) - %s - %s", g.Channel.Name, g.Viewers, g.Channel.Status, g.Channel.URL)
		sslice = append(sslice, s)
	}

	return sslice
}

func TopDota2Streams() {
	requestURL := "https://api.twitch.tv/kraken/streams?game=Dota+2&limit=10"
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	streams, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dat JSONResult
	if err := json.Unmarshal(streams, &dat); err != nil {
		panic(err)
	}

	for _, g := range dat.Streams {
		if !isBlacklisted(g.Channel.Name) {
			fmt.Println("Stream: " + g.Channel.Name + " - " + g.Channel.Status + " - " + g.Channel.URL)
		}
	}
}

func clientID() string {
	file, e := ioutil.ReadFile("./client.id")
	if e != nil {
		panic(e)
	}
	return string(file)
}

func favoriteStreams() string {
	file, e := ioutil.ReadFile("./favorites.txt")
	if e != nil {
		panic(e)
	}
	return string(file)
}

func blacklistStreams() []string {
	file, e := ioutil.ReadFile("./blacklist.txt")
	if e != nil {
		panic(e)
	}
	return strings.Split(string(file), "\n")
}

func isBlacklisted(stream string) bool {
	blacklist := blacklistStreams()
	for _, b := range blacklist {
		if b == stream {
			return true
		}
	}
	return false
}

// JSON structs
type JSONResult struct {
	Streams []JSONStreams `json:"streams"`
}

type JSONStreams struct {
	Channel JSONChannel `json:"channel"`
	Viewers int         `json:"viewers"`
}

type JSONChannel struct {
	Name   string `json:"display_name"`
	URL    string `json:"url"`
	Status string `json:"status"`
}
