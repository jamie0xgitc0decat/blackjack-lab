package blackjack

type BlackJackActionType string

const (
	Hit        BlackJackActionType = "Hit"
	Stand      BlackJackActionType = "Stand"
	DoubleDown BlackJackActionType = "DoubleDown"
	Split      BlackJackActionType = "Split"
	Insurance  BlackJackActionType = "Insurance"
	Surrender  BlackJackActionType = "Surrender"
)

type BlackJackAction struct {
	Action        BlackJackActionType
	AdditionalBet int // Used for actions like doule down or split where additional bets are placed.
}
