package blackjack

import Deck "blackjack-lab/deck"

type Player struct {
	Hands   []Hand
	Actions [][]BlackJackAction // A slice of actions for each hand
	Wallet  int
}

func (p *Player) NewBet(bet int) {
	p.Hands = []Hand{{Cards: []Deck.Card{}, Bet: bet}}
	p.Actions = [][]BlackJackAction{}
}

func (p *Player) AddNewAction(action BlackJackAction, handIndex int) {
	for len(p.Actions) <= handIndex {
		p.Actions = append(p.Actions, []BlackJackAction{})
	}
	p.Actions[handIndex] = append(p.Actions[handIndex], action)
}

func (p *Player) GetLastAction(handIndex int) BlackJackAction {
	if len(p.Actions) > handIndex && len(p.Actions[handIndex]) > 0 {
		return p.Actions[handIndex][len(p.Actions[handIndex])-1]
	}
	return BlackJackAction{} // Return an empty action if none exists
}

func (p *Player) CanMadeAction(handIndex int) bool {
	if p.Hands[handIndex].IsBusted() || p.Hands[handIndex].Score() == 21 {
		return false
	}
	lastAction := p.GetLastAction(handIndex)
	if lastAction.Action == DoubleDown || lastAction.Action == Stand || lastAction.Action == Surrender {
		return false
	}
	return len(p.Hands) <= 4 && len(p.Hands[handIndex].Cards) <= 4
}

func (p *Player) GetTotalBet() int {
	totalBet := 0
	for _, hand := range p.Hands {
		totalBet += hand.Bet
	}
	return totalBet
}

// Implement the split method
func (p *Player) Split() {
	handIndex := len(p.Hands) - 1
	if !p.CanSplit(handIndex) {
		return
	}

	hand := p.Hands[handIndex]
	// Assume the first two cards are a pair and eligible for splitting
	newHands := []Hand{{Cards: []Deck.Card{hand.Cards[0]}, Bet: hand.Bet}, {Cards: []Deck.Card{hand.Cards[1]}, Bet: hand.Bet}}

	// Replace the original hand with one of the new hands and append the other
	p.Hands[handIndex] = newHands[0]
	p.Hands = append(p.Hands, newHands[1])

}

// Check if the player can split the hand
func (p *Player) CanSplit(handIndex int) bool {
	hand := p.Hands[handIndex]
	return len(hand.Cards) == 2 && hand.Cards[0].Value == hand.Cards[1].Value
}
