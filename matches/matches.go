// Package matches provides access to Dota 2 matches through the offical API.
package matches

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MajorScore returns the current score of the Dota 2 Majors match, if one is found.
func MajorScore() string {
	requestURL := "https://api.steampowered.com/IDOTA2Match_570/GetLiveLeagueGames/v0001/?league_id=4266&key=86F1ACC15C5F0A97465AA051D68122F6"
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	liveLeagueData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var games JSONLiveLeagueGamesRoot
	if err := json.Unmarshal(liveLeagueData, &games); err != nil {
		panic(err)
	}

	for _, g := range games.Result.Games {
		min := int(g.Scoreboard.Duration / 60)
		return fmt.Sprintf("%v (%d) vs %v (%d)%s. %d-%d kills %d minutes in.\n", g.TeamRadiant.TeamName, g.RadiantSeriesWin, g.TeamDire.TeamName, g.DireSeriesWin, nameForStage(g.StageName), g.Scoreboard.Radiant.Score, g.Scoreboard.Dire.Score, min)
	}
	return "No Major match found."
}

func nameForStage(s string) string {
	log.Println(s)
	bracket := ""
	switch s {
	case "#DOTA_TournamentBracket_UBQuarterFinals":
		bracket = " in the UB Quarter Finals (BO3)"
	case "#DOTA_TournamentBracket_UBSemiFinals":
		bracket = " in the UB Semi Finals (BO3)" // UB Quarter Semifinals actually
	case "#DOTA_TournamentBracket_UBFinalsMISSING":
		bracket = " in the UB Semi Finals (BO3)"
	case "#DOTA_TournamentBracket_FinalsMISSING":
		bracket = " in the Grand Final (BO5)"
	case "#DOTA_TournamentBracket_LBR1":
		bracket = " in the LB Round 1 (BO1)"
	case "#DOTA_TournamentBracket_LBR2":
		bracket = " in the LB Round 2 (BO1)"
	case "#DOTA_TournamentBracket_LBR3":
		bracket = " in the LB Round 3 (BO3)"
	case "#DOTA_TournamentBracket_LBR4":
		bracket = " in the LB Round 4 (BO3)"
	case "#DOTA_TournamentBracket_LBR5":
		bracket = " in the LB Round 5 (BO3)"
	case "#DOTA_TournamentBracket_LBFinals":
		bracket = " in the LB Finals (BO3)"
	}
	return bracket
}

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

type JSONLiveLeagueGamesRoot struct {
	Result JSONLiveLeagueGames `json:"result"`
}

type JSONLiveLeagueGames struct {
	Games []JSONGame `json:"games"`
}

type JSONGame struct {
	Players          []JSONPlayer `json:"players"`
	TeamRadiant      JSONTeam     `json:"radiant_team"`
	TeamDire         JSONTeam     `json:"dire_team"`
	LobbyID          int          `json:"lobby_id"`
	LeagueID         int          `json:"league_id"`
	Scoreboard       JSONScore    `json:"scoreboard"`
	DireSeriesWin    int          `json:"dire_series_wins"`
	RadiantSeriesWin int          `json:"radiant_series_wins"`
	StageName        string       `json:"stage_name"`
}

type JSONTeam struct {
	TeamName string `json:"team_name"`
	TeamID   int    `json:"team_id"`
}

type JSONPlayer struct {
	Name     string `json:"name"`
	AcountID int    `json:"account_id"`
}

type JSONScore struct {
	Radiant  JSONTeamScore `json:"radiant"`
	Dire     JSONTeamScore `json:"dire"`
	Duration float32       `json:"duration"`
}

type JSONTeamScore struct {
	Score int `json:"score"`
}

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

type JSONStreamRoot struct {
	Result []JSONStream `json:"streams"`
}

type JSONStream struct {
	LeagueID int    `json:"leagueid"`
	EN       string `json:"en"`
}
