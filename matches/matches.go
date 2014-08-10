package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	games := liveLeagueGames()
	for _, g := range games.Result.Games {
		title := (g.LeagueID)
		text := (g.LeagueID)
		fmt.Printf("* %v %v (%v) - %v vs %v - %v\n", g.LobbyID, title, g.LeagueID, g.TeamRadiant.TeamName, g.TeamDire.TeamName, text)
	}
}

func liveLeagueGames() JSONLiveLeagueGamesRoot {
	requestURL := "https://api.steampowered.com/IDOTA2Match_570/GetLiveLeagueGames/v0001/?key=" + apiKey()
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	liveLeagueData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dat JSONLiveLeagueGamesRoot
	if err := json.Unmarshal(liveLeagueData, &dat); err != nil {
		panic(err)
	}
	return dat
}

func apiKey() string {
	file, e := ioutil.ReadFile("./api.key")
	if e != nil {
		panic(e)
	}
	return string(file)
}

// JSON structs for LiveLeagueGames
type JSONLiveLeagueGamesRoot struct {
	Result JSONLiveLeagueGames `json:"result"`
}

type JSONLiveLeagueGames struct {
	Games []JSONGame `json:"games"`
}

type JSONGame struct {
	Players     []JSONPlayer `json:"players"`
	TeamRadiant JSONTeam     `json:"radiant_team"`
	TeamDire    JSONTeam     `json:"dire_team"`
	LobbyID     int          `json:"lobby_id"`
	LeagueID    int          `json:"league_id"`
}

type JSONTeam struct {
	TeamName string `json:"team_name"`
	TeamID   int    `json:"team_id"`
}

type JSONPlayer struct {
	Name     string `json:"name"`
	AcountID int    `json:"account_id"`
}

// JSON struct for GetLeagueListing
type JSONLeagueListingRoot struct {
	Result JSONLeagueListing `json:"result"`
}

type JSONLeagueListing struct {
	Leagues []JSONLeague `json:"leagues"`
}

type JSONLeague struct {
	Name     string `json:"name"`
	LeagueID int    `json:"leagueid"`
	URL      string `json:"tournament_url"`
}

// JSON struct for streams
type JSONStreamRoot struct {
	Result []JSONStream `json:"streams"`
}

type JSONStream struct {
	LeagueID int    `json:"leagueid"`
	EN       string `json:"en"`
}
