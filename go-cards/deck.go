package main

import "fmt"

type Deck []string

func newDeck() Deck {
	cards := Deck{}

	cardSuits := []string{"Spades", "Hearts", "Clubs", "Diamonds"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}

	for _, suit := range cardSuits {
		for _, card := range cardValues {
			cards = append(cards, card+" of "+suit)
		}
	}

	return cards
}

// receiver function
func (d Deck) logAllCards() {
	for i, card := range d {
		fmt.Println("========")
		fmt.Println("i: ", i)
		fmt.Println("card: ", card)
	}
}

func (d Deck) deal(handSize int) (Deck, Deck) {
	return d[:handSize], d[handSize:]
}
