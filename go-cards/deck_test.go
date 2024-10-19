package main

import (
	"log"
	"os"
	"testing"
)

// 1. Assert a new deck should have 52 cards
func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(d))
	}

	if d[0] != "Ace of Spades" {
		t.Errorf("Expected the first card in deck to be 'Ace of Spades', but got %v", d[0])
	}

	if d[len(d)-1] != "King of Diamonds" {
		t.Errorf("Expected the first card in deck to be 'King of Diamonds', but got %v", d[len(d)-1])
	}
}

func TestSaveDeckToFileAndNewDeckFromFile(t *testing.T) {
	err := os.Remove("_deck-testing")
	if err != nil {
		// Safe to ignore
	}

	d := newDeck()
	err = d.saveToFile("_deck-testing")
	if err != nil {
		log.Print("Error trying to save test deck to '_deck-testing'file' file")
	}

	// New deck from file
	loadedDeck, loadErr := getSavedDeckFromFile("_deck-testing")
	if loadErr != nil {
		log.Print("Error trying to remove '_deck-testing' file.")
	}

	if len(loadedDeck) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(loadedDeck))
	}

	err = os.Remove("_deck-testing")
	if err != nil {
		log.Print("Error trying to remove '_deck-testing' file.")
	}
}
