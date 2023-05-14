package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
    // siteUrl := "https://scrapeme.live/shop/"
    siteUrl := "https://www.betfair.com.au/exchange/plus/australian-rules/market/1.213699333"
    

    // initializing the slice of structs that will contain the scraped data 
    // var pokemonProducts []PokemonProduct

	c := colly.NewCollector()
    c.UserAgent = "Go program"
    
    c.OnHTML(".mv-runner-list", func(e *colly.HTMLElement) {
        println("e:", e)
    })

    c.Visit(siteUrl)
}

func loggingSiteDetails(c *colly.Collector) {
    c.OnRequest(func(r *colly.Request) {

        for key, value := range *r.Headers {
            fmt.Printf("%s: %s\n", key, value)
        }

        fmt.Println(r.Method)
    })

    c.OnHTML("title", func(e *colly.HTMLElement) {

        fmt.Println("-----------------------------")

        fmt.Println(e.Text)
    })

    c.OnResponse(func(r *colly.Response) {

        fmt.Println("-----------------------------")

        fmt.Println(r.StatusCode)

        for key, value := range *r.Headers {
            fmt.Printf("%s: %s\n", key, value)
        }
    })
}


// Types

// defining a data structure to store the scraped data 
type PokemonProduct struct { 
	url, image, name, price string 
}


// c.OnHTML("li.product", func(e *colly.HTMLElement) { 
//     // initializing a new PokemonProduct instance 
//     pokemonProduct := PokemonProduct{} 

//     // scraping the data of interest 
//     pokemonProduct.url = e.ChildAttr("a", "href") 
//     pokemonProduct.image = e.ChildAttr("img", "src") 
//     pokemonProduct.name = e.ChildText("h2") 
//     pokemonProduct.price = e.ChildText(".price") 

//     // adding the product instance with scraped data to the list of products 
//     pokemonProducts = append(pokemonProducts, pokemonProduct) 
    
// })