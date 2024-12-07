package blackjack

import (
	deck "blackjack-lab/deck"
	"fmt"
	"strings"
)

type WinningResult int

const (
	// Define constants for the possible outcomes
	Loss WinningResult = iota // Loss is the default value with iota starting at 0
	Win
	Tie
)

type HandType int

const (
	HardTotalsHandType HandType = iota
	SoftTotalsHandType
	PairSplitHandType
)

type Hand struct {
	Cards         []deck.Card
	Bet           int
	WinningResult WinningResult
	HandType      HandType
}

func (h *Hand) Score() int {
	score := 0
	aces := 0

	for _, card := range h.Cards {
		if card.Rank == 11 {
			aces++
		}
		score += card.Rank
	}

	// Adjust for Aces after the initial calculation.
	for score > 21 && aces > 0 {
		score -= 10
		aces--
	}

	return score
}

func (h *Hand) IsBusted() bool {
	return h.Score() > 21
}

func (h *Hand) GetHandType() HandType {
	if len(h.Cards) == 2 && h.Cards[0].Value == h.Cards[1].Value {
		return PairSplitHandType
	}

	if h.HasAce() {
		return SoftTotalsHandType
	}

	return HardTotalsHandType
}

func (h *Hand) HasAce() bool {
	for _, card := range h.Cards {
		if card.Value == deck.Ace {
			return true
		}
	}
	return false
}

func (h *Hand) HasPairsJackpot() bool {
	return len(h.Cards) >= 2 && h.Cards[0].Value == h.Cards[1].Value
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Score() == 21
}

func (h *Hand) CanSplit() bool {
	return len(h.Cards) == 2 && h.Cards[0].Value == h.Cards[1].Value
}

func (h *Hand) String() string {
	var cards []string
	for _, card := range h.Cards {
		cards = append(cards, fmt.Sprintf("%s%s", card.Value, card.Suit))
	}
	return strings.Join(cards, ", ")
}
