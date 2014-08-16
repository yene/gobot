package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var favorites []string

func WatchFavorites(callback func(m string)) {
	favorites = FavoriteDota2Streams()
	for {
		time.Sleep(time.Second * 30)
		newFavorites := FavoriteDota2Streams()
		if len(newFavorites) == 0 {
			continue // sometimes the api delivers no results
		}

		for _, g := range newFavorites {
			if !inisdeFavorites(g) {
				callback(g + " started streaming.")
			}
		}
		favorites = newFavorites
	}
}

func inisdeFavorites(a string) bool {
	for _, g := range favorites {
		if g == a {
			return true
		}
	}
	return false
}

func FavoriteDota2Streams() []string {
	favorites := favoriteStreams()
	concatenated := strings.Replace(favorites, "\n", ",", -1)
	requestURL := "https://api.twitch.tv/kraken/streams?game=Dota+2&channel=" + concatenated
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
		s := fmt.Sprintf("\u0002%s\u000F %s", g.Channel.DisplayName, g.Channel.URL)
		sslice = append(sslice, s)
	}

	return sslice
}

func TopDota2Streams() []string {
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

	limitOfStreams := 5
	c := 0

	sslice := make([]string, 0)
	for _, g := range dat.Streams {
		if c == limitOfStreams {
			break
		}
		if !isBlacklisted(g.Channel.Name) && g.Viewers > 800 && !isRebroadcast(g.Channel.Status) {
			s := fmt.Sprintf("\u0002%s\u000F (%d) %s", g.Channel.DisplayName, g.Viewers, g.Channel.URL)
			sslice = append(sslice, s)
			c++
		}
	}

	return sslice
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

func isRebroadcast(stream string) bool {
	s := strings.ToLower(stream)
	return strings.Contains(s, "rebroadcast")
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
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Status      string `json:"status"`
}
