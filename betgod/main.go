package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

// Notes
// 1.For GoColly, you need to set up the listeners first then run your visit() command.
// 2.

func main() {
	fmt.Println("System Online and Ready Sir")

	getPageLinks()

	fmt.Println("Scraping finished")
}

func getPageLinks() {
	// Scrape AFL season data
	c := colly.NewCollector()

	c.OnHTML("a", func(e *colly.HTMLElement) {
		// printing all URLs associated with the <a> links in the page
		fmt.Printf("%v\n", e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	baseUrl := "https://afltables.com/afl/seas/"
	// Root site, used to find URL addresses for all seasons
	err := c.Visit(fmt.Sprintf("%s%s", baseUrl, "season_idx.html"))
	if err != nil {
		log.Printf("Error occured bra: %+v", err)
		log.Fatal(err)
	}
}
