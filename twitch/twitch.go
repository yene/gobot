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

type Channel struct {
	Name     string
	Status   string
	URL      string
	Language string
	Viewers  int
}

var favorites []string
var tournaments []string
var all []string

func UpdateStreams(streams chan []Channel) {
	for {
		time.Sleep(time.Second * 30)

		// get all dota streams, even russians oO
		requestURL := "https://api.twitch.tv/kraken/streams?game=Dota+2"
		res, err := http.Get(requestURL)
		if err != nil {
			log.Print(err)
			continue
		}
		streamsData, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Print(err)
			continue
		}

		var dat JSONResult
		if err := json.Unmarshal(streamsData, &dat); err != nil {
			fmt.Println("error unmarshalling Dota2streamsData dota2 stream")
			ioutil.WriteFile("streamsData.json", streamsData, 0644)
			log.Print(err)
			continue
		}

		cslice := make([]Channel, 0)
		for _, g := range dat.Streams {
			if len(g.Channel.URL) == 0 { // skip channels without url, a twitch bug
				continue
			}
			// what do i need url, display status
			c := Channel{g.Channel.DisplayName, g.Channel.Status, g.Channel.URL, g.Channel.Language, g.Viewers}
			cslice = append(cslice, c)
		}

		streams <- cslice
	}
}

func WatchFavorites(streams chan []Channel, callback func(m string)) {
	for s := range streams {
		bla := ""
		for _, c := range s {
			bla += fmt.Sprintf("%s %s %s %s", c.Name, c.Status, c.URL, c.Language)

		}
		fmt.Println(bla)
	}
	/*
		favorites = FavoriteDota2Streams()
		for {
			time.Sleep(time.Second * 30)
			newFavorites := FavoriteDota2Streams()
			if len(newFavorites) == 0 {
				continue // sometimes the api delivers no results
			}

			for _, g := range newFavorites {
				if !inside(favorites, g) {
					callback(g + " started streaming.")
				}
			}
			favorites = newFavorites
		}*/
}

func WatchTournaments(callback func(m string)) {
	tournaments = TournamentStreams()
	for {
		time.Sleep(time.Second * 30)
		newTournaments := TournamentStreams()
		if len(newTournaments) == 0 {
			continue // sometimes the api delivers no results
		}

		for _, g := range newTournaments {
			if !inside(tournaments, g) {
				callback(g)
			}
		}
		tournaments = newTournaments
	}
}

func WatchAll(callback func(m string)) {
	all = Dota2Streams()
	for {
		time.Sleep(time.Second * 10)
		newAll := Dota2Streams()
		if len(newAll) == 0 {
			continue // sometimes the api delivers no results
		}

		for _, g := range newAll {
			if !inside(all, g) {
				callback(g)
			}
		}
		all = newAll
	}
}

func inside(haystack []string, needle string) bool {
	for _, g := range haystack {
		if g == needle {
			return true
		}
	}
	return false
}

func FavoriteDota2Streams() []string {
	f := favoriteList()
	concatenated := strings.Replace(f, "\n", ",", -1)
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
		fmt.Println("error unmarshalling favorite dota2 stream")
		ioutil.WriteFile("favorites.json", streams, 0644)
		log.Fatal(err)
	}

	sslice := make([]string, 0)
	for _, g := range dat.Streams {
		s := fmt.Sprintf("\u0002%s\u000F %s", g.Channel.DisplayName, g.Channel.URL)
		if len(g.Channel.URL) == 0 {
			continue
		}
		sslice = append(sslice, s)
	}

	return sslice
}

func TournamentStreams() []string {
	t := tournamentsList()
	concatenated := strings.Replace(t, "\n", ",", -1)
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
		fmt.Println("error unmarshalling tournament dota2 stream")
		ioutil.WriteFile("favorites.json", streams, 0644)
		log.Fatal(err)
	}

	sslice := make([]string, 0)
	for _, g := range dat.Streams {
		if len(g.Channel.URL) == 0 {
			continue
		}

		if isRebroadcast(g.Channel.Status) {
			continue
		}

		if containsVersus(g.Channel.Status) || containsLive(g.Channel.Status) {
			s := fmt.Sprintf("%s %s", g.Channel.Status, g.Channel.URL)
			sslice = append(sslice, s)
		}
	}

	return sslice
}

func FilteredDota2Streams() []string {
	requestURL := "https://api.twitch.tv/kraken/streams?game=Dota+2&broadcaster_language=en&limit=15"
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
		fmt.Println("error unmarshalling filtered dota2 stream")
		ioutil.WriteFile("favorites.json", streams, 0644)
		log.Fatal(err)
	}

	limitOfStreams := 5
	c := 0

	sslice := make([]string, 0)
	for _, g := range dat.Streams {
		if c == limitOfStreams {
			break
		}

		if len(g.Channel.URL) == 0 {
			continue
		}

		if !isBlacklisted(g.Channel.Name) && g.Viewers > 100 && !isRebroadcast(g.Channel.Status) {
			s := fmt.Sprintf("\u0002%s\u000F (%d) %s", g.Channel.DisplayName, g.Viewers, g.Channel.URL)
			sslice = append(sslice, s)
			c++
		}
	}

	return sslice
}

func Dota2Streams() []string {
	// get all dota streams, even russians oO
	requestURL := "https://api.twitch.tv/kraken/streams?game=Dota+2"
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
		fmt.Println("error unmarshalling Dota2Streams dota2 stream")
		ioutil.WriteFile("favorites.json", streams, 0644)
		log.Fatal(err)
	}

	sslice := make([]string, 0)
	for _, g := range dat.Streams {
		if len(g.Channel.URL) == 0 {
			continue
		}
		s := fmt.Sprintf("\u0002%s\u000F %s", g.Channel.DisplayName, g.Channel.URL)
		sslice = append(sslice, s)
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

func favoriteList() string {
	file, e := ioutil.ReadFile("./favorites.txt")
	if e != nil {
		panic(e)
	}
	return string(file)
}

func tournamentsList() string {
	file, e := ioutil.ReadFile("./tournaments.txt")
	if e != nil {
		panic(e)
	}
	return string(file)
}

func isRebroadcast(stream string) bool {
	s := strings.ToLower(stream)
	return strings.Contains(s, "rebroadcast")
}

func containsVersus(stream string) bool {
	s := strings.ToLower(stream)
	return strings.Contains(s, " vs ")
}

func containsLive(stream string) bool {
	s := strings.ToLower(stream)
	return strings.Contains(s, "live")
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

func whitelistStreams() []string {
	file, e := ioutil.ReadFile("./whitelist.txt")
	if e != nil {
		panic(e)
	}
	return strings.Split(string(file), "\n")
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
	Language    string `json:"broadcaster_language"`
}
