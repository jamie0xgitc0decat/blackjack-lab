// deck/deck.go

package deck

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Suit represents the suit of the card.
type Suit string

// Value represents the value of the card.
type Value string

type CardVal struct {
	Value string
	Rank  int
}

// Define suits
const (
	Spades   Suit = "S"
	Diamonds Suit = "D"
	Clubs    Suit = "C"
	Hearts   Suit = "H"
)

// Define card values
const (
	Ace   Value = "A"
	Two   Value = "2"
	Three Value = "3"
	Four  Value = "4"
	Five  Value = "5"
	Six   Value = "6"
	Seven Value = "7"
	Eight Value = "8"
	Nine  Value = "9"
	Ten   Value = "10"
	Jack  Value = "J"
	Queen Value = "Q"
	King  Value = "K"
)

func GetCardRank(v Value) int {
	switch v {
	case Ace:
		return 11
	case Jack, Queen, King:
		return 10
	default:
		num, err := strconv.Atoi(string(v))
		if err != nil {
			fmt.Println("Conversion error:", err)
		}
		return num
	}
}

// Card represents a playing card with a suit and a value.
type Card struct {
	Suit  Suit
	Value Value
	Rank  int
}

// New creates and returns a new deck of cards.
func New() []Card {
	suits := []Suit{Spades, Diamonds, Clubs, Hearts}
	values := []Value{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	var deck []Card
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value, Rank: GetCardRank(value)})
		}
	}

	return deck
}

func NewDeckCards(numOfDecks int) []Card {
	deck := New()
	var decks []Card
	for i := 0; i < numOfDecks; i++ {
		decks = append(decks, deck...)
	}
	rand.Shuffle(len(decks), func(i, j int) { decks[i], decks[j] = decks[j], decks[i] })
	return decks
}

// Shuffle randomizes the order of the cards in the deck.
func Shuffle(deck []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// Deal returns the first n cards from the deck and the remaining deck.
func Deal(deck []Card, n int) ([]Card, []Card) {
	return deck[:n], deck[n:]
}
