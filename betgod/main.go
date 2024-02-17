package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"log"
	"strings"
)

// Notes
// 1.For GoColly, you need to set up the listeners first then run your visit() command.
// 2.

func main() {
	fmt.Println("System Online and Ready Sir")

	//baseUrl := "https://afltables.com/afl/seas/"
	// afl links to each season
	//aflSeasonsList := getPageLinks(baseUrl)

	// Test navigation to one page + scrape
	//testUrl := "https://afltables.com/afl/seas/2023.html"
	//getPageStats(testUrl)

	// Test string mapping / interpolation
	handleStringMagic()

	fmt.Println("Scraping finished")
	//fmt.Println(aflSeasonsList)
}

func handleStringMagic() {
	// Struct to contain full match data
	var MatchResult = MatchStats{}
	teamOneSet := false
	//testStr := "Richmond  1.4   2.4   7.8  8.10  58Thu 16-Mar-2023 7:20 PM (6:20 PM) Att: 88,084 Venue: M.C.G.\nCarlton  3.1   4.6   6.9  8.10  58Match drawn [Match stats]\n"
	testStr := "St Kilda  3.3   5.4   6.6  10.7  67Sun 19-Mar-2023 4:40 PM (3:40 PM) Att: 23,429 Venue: Docklands\nFremantle  2.1   5.4   7.6  7.10  52St Kilda won by 15 pts [Match stats]\n"

	lines := strings.Split(testStr, "\n")

	for _, line := range lines {
		for team := range TeamNames {
			if strings.Contains(line, team) {
				stats := ExtractTeamStats(line, team)
				// First team is away
				if !teamOneSet {
					MatchResult.TeamOne = stats
					teamOneSet = true
				}
				// Second team is home
				MatchResult.TeamTwo = stats

				//fmt.Println("Game Stats:", stats)
				break
			}
		}
	}

	MatchResult.MatchID = uuid.New().String()
	// Find match winner
	tempTeamOneOutcome := MatchResult.TeamOne.FinalScore
	tempTeamTwoOutcome := MatchResult.TeamTwo.FinalScore

	if tempTeamOneOutcome > tempTeamTwoOutcome {
		MatchResult.TeamOne.MatchResult = "WIN"
		MatchResult.TeamTwo.MatchResult = "LOSS"
		MatchResult.WinningTeam = MatchResult.TeamOne.TeamName
	} else if tempTeamTwoOutcome < tempTeamOneOutcome {
		MatchResult.TeamOne.MatchResult = "LOSS"
		MatchResult.TeamTwo.MatchResult = "WIN"
		MatchResult.WinningTeam = MatchResult.TeamTwo.TeamName
	} else {
		// It's a draw
		MatchResult.TeamOne.MatchResult = "DRAW"
		MatchResult.TeamTwo.MatchResult = "DRAW"
	}

	fmt.Println(MatchResult)
}

func getPageStats(url string) {
	fmt.Println("Scraping: ")
	fmt.Println(url)
	count := 1

	c := colly.NewCollector()
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if count == 3 {
			return
		}

		// Every 2nd table on the page has the data we require
		// Ignore round number + we start at round 1.
		if count%2 == 0 {
			fmt.Println(e.Text)
		}
		count++
	})

	// Extract data
	// 1. Get team name
	// 2. Get quarter by quarter team scores
	// 3. Get quarter by quarter leader
	// 4. Game Result
	// 5. Game details (string)

	err := c.Visit(url)
	if err != nil {
		log.Printf("Error occured bra: %+v", err)
		log.Fatal(err)
	}
}

func getPageLinks(rootURL string) []string {
	// Scrape AFL season data
	c := colly.NewCollector()
	var aflSeasonsList []string
	// We only want the last 25 seasons (for now)
	// Because the season data starts back in 1897 through to 2023
	linkCount := 0
	totalGames := 126
	totalGamesToScrap := 25

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if linkCount < (totalGames - totalGamesToScrap) {
			linkCount++
			return
		}
		// Returning all <a> tag links on page
		aflSeasonsList = append(aflSeasonsList, fmt.Sprintf("%s%s", rootURL, e.Attr("href")))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	// Root site, used to find URL addresses for all seasons
	err := c.Visit(fmt.Sprintf("%s%s", rootURL, "season_idx.html"))
	if err != nil {
		log.Printf("Error occured bra: %+v", err)
		log.Fatal(err)
	}

	return aflSeasonsList
}
