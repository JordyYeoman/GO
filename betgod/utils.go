package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// Remember to use uppercase declaration if you want to export
var TeamNames = map[string]bool{
	"Richmond":               true,
	"Carlton":                true,
	"Geelong":                true,
	"Collingwood":            true,
	"North Melbourne":        true,
	"West Coast":             true,
	"Port Adelaide":          true,
	"Brisbane Lions":         true,
	"Melbourne":              true,
	"Western Bulldogs":       true,
	"Gold Coast":             true,
	"Sydney":                 true,
	"Greater Western Sydney": true,
	"Adelaide":               true,
	"Hawthorn":               true,
	"Essendon":               true,
	"St Kilda":               true,
	"Fremantle":              true,
}

func StripDigitsFromString(s string) string {
	var tempVar string

	for _, char := range s {
		if !unicode.IsDigit(char) {
			tempVar += string(char)
		}
	}

	return tempVar
}

func GetMatchData(sliceOfStrings []string) string {
	var tempStr string
	// Part 5 is int + day, we just want the day
	for i, item := range sliceOfStrings {
		if i == 5 {
			tempStr += StripDigitsFromString(item)
		}
		if i > 5 {
			tempStr += " "
			tempStr += item
		}
	}

	return tempStr
}

func GetFinalScore(str string) int {
	var tempScore string
	for _, char := range str {
		if unicode.IsDigit(char) {
			tempScore += string(char)
		} else {
			break // Stop iteration if non-digit character encountered
		}
	}

	s, err := strconv.Atoi(tempScore)
	if err != nil {
		fmt.Println("Can't convert this to an int!")
	}

	return s
}

func (s *TeamStats) SetQuarterScore(quarter, score int) {
	switch quarter {
	case 1:
		s.QuarterOneScore += score
	case 2:
		s.QuarterTwoScore += score
	case 3:
		s.QuarterThreeScore += score
	case 4:
		s.QuarterFourScore += score
	}
}

func RemoveTeamName(line, team string) string {
	return strings.TrimPrefix(line, team)
}

func ExtractTeamStats(line, team string) TeamStats {
	var stats TeamStats
	stats.TeamName = team
	parts := strings.Fields(line) // Split the line by spaces
	endOfTeamScoresInStringSplit := 4

	fmt.Printf("parts: %+v", parts)
	// TODO: Check for special team names that are not just 1 word eg - 'st kilda'

	// Final Score
	stats.FinalScore = GetFinalScore(parts[endOfTeamScoresInStringSplit])

	// Match data
	stats.MatchData = GetMatchData(parts)

	// Quarters
	for i := 0; i < endOfTeamScoresInStringSplit; i++ {
		score := parts[i]
		scoreParts := strings.Split(score, ".")
		if len(scoreParts) != 2 {
			fmt.Println("Invalid score format:", score)
			continue
		}
		score1, err1 := strconv.Atoi(scoreParts[0])
		score2, err2 := strconv.Atoi(scoreParts[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting score to int:", score)
			continue
		}
		quarter := i + 1 // Quarter 1 corresponds to index 1, Quarter 2 to index 2, and so on
		switch quarter {
		case 1:
			stats.QuarterOneData = score
			stats.QuarterOneScore = score1*6 + score2
		case 2:
			stats.QuarterTwoData = score
			stats.QuarterTwoScore = score1*6 + score2
		case 3:
			stats.QuarterThreeData = score
			stats.QuarterThreeScore = score1*6 + score2
		case 4:
			stats.QuarterFourData = score
			stats.QuarterFourScore = score1*6 + score2
		}
	}

	return stats
}

// Used to extract page links from the root page.
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
