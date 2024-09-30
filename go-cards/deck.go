package main

import (
	"fmt"
	"os"
	"strings"
)

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

func (d Deck) toString() string {
	return strings.Join(d, ",")
}

func (d Deck) saveToFile(fileName string) error {
	return os.WriteFile(fileName, []byte(d.toString()), 0666)
}

func getSavedDeckFromFile(fileName string) (Deck, error) {
	bs, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	s := strings.Split(string(bs), ",")
	return s, nil
}
