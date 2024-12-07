package blackjack

import (
	Deck "blackjack-lab/deck"
	"fmt"
)

type BlackjackTable struct {
	Deck       []Deck.Card
	Dealer     Hand
	Players    []Player
	NumOfDecks int
	MinBet     int
}

type BlackjackTableConfigParams struct {
	NumOfDecks int
	MinBet     int
}

func (t *BlackjackTable) NewTable(params BlackjackTableConfigParams) {
	t.NumOfDecks = params.NumOfDecks
	t.MinBet = params.MinBet
}

func (t *BlackjackTable) NewGame() {
	t.Deck = Deck.NewDeckCards(t.NumOfDecks)
	t.Dealer = Hand{}
	t.Players = []Player{}
}

func (t *BlackjackTable) InitialDeal() {
	for i := 0; i < 2; i++ {
		for j := range t.Players {
			drawCard, remainingDeck := Deck.Deal(t.Deck, 1)
			t.Players[j].Hands[0].Cards = append(t.Players[j].Hands[0].Cards, drawCard[0])
			t.Deck = remainingDeck
		}
		if i == 0 {
			drawCard, remainingDeck := Deck.Deal(t.Deck, 1)
			t.Dealer.Cards = append(t.Dealer.Cards, drawCard[0])
			t.Deck = remainingDeck
		}
	}
}

func (t *BlackjackTable) DealToDealer() {
	for t.Dealer.Score() < 17 {
		drawCard, remainingDeck := Deck.Deal(t.Deck, 1)
		t.Dealer.Cards = append(t.Dealer.Cards, drawCard[0])
		t.Deck = remainingDeck
	}
}

func (t *BlackjackTable) DealToPlayer(playerIndex int, handIndex int) {
	drawCard, remainingDeck := Deck.Deal(t.Deck, 1)
	t.Players[playerIndex].Hands[handIndex].Cards = append(t.Players[playerIndex].Hands[handIndex].Cards, drawCard[0])
	t.Deck = remainingDeck
}

func (t *BlackjackTable) ShowTableInfo() string {
	var info string
	info += "Dealer: " + t.Dealer.String() + "\n"

	// for i, player := range t.Players {
	// 	info += fmt.Sprintf("Player %d: %s, score: %d\n", i+1, player.Hand.String(), player.Hand.Score())
	// }
	fmt.Println(info)
	return info
}

func (t *BlackjackTable) DealerScore() int {
	return t.Dealer.Score()
}

func (t *BlackjackTable) PlayerScore(i int, handIndex int) int {
	return t.Players[i].Hands[handIndex].Score()
}

func (t *BlackjackTable) PlayerBet(i int, handIndex int) int {
	return t.Players[i].Hands[handIndex].Bet
}

func (t *BlackjackTable) AddPlayer(Player *Player) {
	t.Players = append(t.Players, *Player)
}
