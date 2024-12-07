package simulator

import (
	Blackjack "blackjack-lab/blackjack"
	Strategy "blackjack-lab/strategy"
	"fmt"
)

type BlackjackSimulator struct {
	BlackjackTable Blackjack.BlackjackTable
}

func GetPlayerHandStrategyType(hand Blackjack.Hand) Strategy.StrategyType {
	switch hand.GetHandType() {
	case Blackjack.PairSplitHandType:
		return Strategy.PairSplittingStrategy
	case Blackjack.SoftTotalsHandType:
		return Strategy.SoftTotalsStrategy
	case Blackjack.HardTotalsHandType:
		return Strategy.HardTotalsStrategy
	}
	return Strategy.HardTotalsStrategy
}

func (s *BlackjackSimulator) ShowStatics() {
	fmt.Println("Number of decks: ", s.BlackjackTable.NumOfDecks)
	fmt.Println("Min bet: ", s.BlackjackTable.MinBet)
	fmt.Println("running num of games: ")
}

func (s *BlackjackSimulator) StartSimulation() {
	maxGames := 500
	numOfPlayers := 1
	eachBet := 300
	numOfDecks := 6

	wallet := 6000.0

	s.BlackjackTable = Blackjack.BlackjackTable{
		NumOfDecks: numOfDecks,
		MinBet:     300,
	}

	strategy, loadStrategyErr := Strategy.GenerateStrategy(Strategy.LoadStrategyParams{
		HardTotalsFilePath: "../strategy/basic_hard.csv",
		SoftTotalsFilePath: "../strategy/basic_soft.csv",
		PairSplittingPath:  "../strategy/basic_pair_split.csv",
	})
	if loadStrategyErr != nil {
		fmt.Println("Error loading strategy: ", loadStrategyErr)
	}

	numOfEvPositive := 0
	numOfPairsJackpot := 0

	currentGame := 0
	for currentGame < maxGames {
		s.BlackjackTable.NewGame()

		for i := 0; i < numOfPlayers; i++ {
			s.BlackjackTable.AddPlayer(&Blackjack.Player{})
			s.BlackjackTable.Players[i].NewBet(eachBet)
		}

		s.BlackjackTable.InitialDeal()
		//s.BlackjackTable.ShowTableInfo()

		for i := 0; i < numOfPlayers; i++ {
			// Player's turn, check player whether is blackjack or not, scores = 21
			for playerHandIndex := 0; playerHandIndex < len(s.BlackjackTable.Players[i].Hands); playerHandIndex++ {
				for s.BlackjackTable.Players[i].CanMadeAction(playerHandIndex) {
					currentPlayerHand := s.BlackjackTable.Players[i].Hands[playerHandIndex]

					playerTotal := currentPlayerHand.Score()
					strategyType := GetPlayerHandStrategyType(currentPlayerHand)
					if (strategyType == Strategy.PairSplittingStrategy) && (len(currentPlayerHand.Cards) == 2) {
						playerTotal = currentPlayerHand.Cards[0].Rank
					}
					// Get the player's action based on the strategy
					playerAction, getActionErr := strategy.GetAction(playerTotal, s.BlackjackTable.Dealer.Cards[0].Rank, strategyType)
					if getActionErr != nil {
						fmt.Printf("Player's Hand %s\n", currentPlayerHand.String())
						fmt.Println("Error getting action: ", getActionErr)
						panic(getActionErr)
					}
					//fmt.Println("Player action: ", playerAction)

					if playerAction == Strategy.Double {
						//initialWallet -= currentPlayerHand.Bet
						currentPlayerHand.Bet *= 2
						s.BlackjackTable.DealToPlayer(i, playerHandIndex)
						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action:        Blackjack.DoubleDown,
							AdditionalBet: currentPlayerHand.Bet,
						},
							playerHandIndex,
						)
					} else if playerAction == Strategy.Hit {
						s.BlackjackTable.DealToPlayer(i, playerHandIndex)
						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action: Blackjack.Hit,
						},
							playerHandIndex,
						)
					} else if playerAction == Strategy.Split {
						// s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
						// 	Action: Blackjack.Stand,
						// }, playerHandIndex)
						s.BlackjackTable.Players[i].Split()
						s.BlackjackTable.DealToPlayer(i, playerHandIndex)
						s.BlackjackTable.DealToPlayer(i, playerHandIndex+1)

						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action:        Blackjack.Split,
							AdditionalBet: currentPlayerHand.Bet,
						},
							playerHandIndex,
						)

					} else if playerAction == Strategy.Stand {
						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action: Blackjack.Stand,
						}, playerHandIndex)
					}

				}
			}
		}

		// settle the dealer's hand
		for {
			s.BlackjackTable.DealToDealer()
			if s.BlackjackTable.Dealer.Score() >= 17 || len(s.BlackjackTable.Dealer.Cards) == 5 {
				break
			}
		}

		currentGameTotalBet := 0.0
		for i := 0; i < numOfPlayers; i++ {
			currentGameTotalBet += float64(s.BlackjackTable.Players[i].GetTotalBet())
		}

		ev := 0.0
		// settle the player's hand
		for i := 0; i < numOfPlayers; i++ {
			// if s.BlackjackTable.Players[i].Hands[0].HasPairsJackpot() {
			// 	numOfPairsJackpot++
			// 	wallet += float64(100) * 11
			// 	ev += float64(100) * 11
			// 	fmt.Println("Player ", i, " pairs jackpot", s.BlackjackTable.Players[i].Hands[0].String())
			// } else {
			// 	wallet -= float64(100)
			// 	ev -= float64(100)
			// }

			for _, playerHand := range s.BlackjackTable.Players[i].Hands {

				playerScore := playerHand.Score()
				if playerScore > 21 {
					wallet -= float64(playerHand.Bet)
					ev -= float64(playerHand.Bet)
					//fmt.Println("Player ", i, " busted")
				} else if playerScore == 21 {
					if i == 0 && playerHand.IsBlackjack() {
						wallet += float64(playerHand.Bet) * 1.5
						ev += float64(playerHand.Bet) * 1.5
					} else {
						wallet += float64(playerHand.Bet)
						ev += float64(playerHand.Bet)

					}

					//fmt.Println("Player ", i, " blackjack")
				} else if playerScore < 21 {
					dealerScore := s.BlackjackTable.Dealer.Score()
					if dealerScore > 21 {
						wallet += float64(playerHand.Bet)
						ev += float64(playerHand.Bet)
						//	fmt.Println("Player ", i, " win")
					} else if dealerScore == playerScore {
						//	fmt.Println("Player ", i, " tie")
						ev += float64(playerHand.Bet)
					} else if dealerScore > playerScore {
						wallet -= float64(playerHand.Bet)
						ev -= float64(playerHand.Bet)

						//	fmt.Println("Player ", i, " lose")
					} else {
						wallet += float64(playerHand.Bet)
						ev += float64(playerHand.Bet)
						//	fmt.Println("Player ", i, " win")
					}
				}
			}

		}

		if ev >= currentGameTotalBet {
			numOfEvPositive++
		}

		if wallet <= 0 {
			break
		}
		currentGame++
	}

	fmt.Println("Number of decks: ", s.BlackjackTable.NumOfDecks)
	fmt.Println("Min bet: ", s.BlackjackTable.MinBet)
	fmt.Println("running num of games: ", currentGame)

	//allWallets := 0

	// for i := 0; i < numOfPlayers; i++ {
	// 	fmt.Println("Player ", i, " wallet: ", s.BlackjackTable.Players[i].Wallet)
	// 	allWallets += s.BlackjackTable.Players[i].Wallet
	// }
	fmt.Println("All wallet: ", wallet)

	fmt.Printf("Number of positive ev: %d\n", numOfEvPositive)
	fmt.Printf("Number of win percent: %.2f\n", float32(numOfEvPositive)/float32(maxGames)*100.0)
	fmt.Printf("Number of pairs jackpot: %d\n", numOfPairsJackpot)
}

func (s *BlackjackSimulator) StartOneMonthSimulation() {
	totalWallet := 0.0
	for i := 0; i < 1; i++ {
		wallet := s.StartOneDaySimulation()
		totalWallet += wallet
		fmt.Println("Day ", i, " wallet: ", totalWallet)
	}

	fmt.Println("Total wallet: ", totalWallet-6000)
}

func (s *BlackjackSimulator) StartOneDaySimulation() float64 {
	maxGames := 100
	numOfPlayers := 1
	eachBet := 300
	numOfDecks := 6

	wallet := 3000.0

	s.BlackjackTable = Blackjack.BlackjackTable{
		NumOfDecks: numOfDecks,
		MinBet:     300,
	}

	strategy, loadStrategyErr := Strategy.GenerateStrategy(Strategy.LoadStrategyParams{
		HardTotalsFilePath: "../strategy/basic_hard.csv",
		SoftTotalsFilePath: "../strategy/basic_soft.csv",
		PairSplittingPath:  "../strategy/basic_pair_split.csv",
	})
	if loadStrategyErr != nil {
		fmt.Println("Error loading strategy: ", loadStrategyErr)
	}

	numOfEvPositive := 0
	numOfPairsJackpot := 0

	highestReward := 0

	currentGame := 0
	for currentGame < maxGames {
		s.BlackjackTable.NewGame()

		for i := 0; i < numOfPlayers; i++ {
			s.BlackjackTable.AddPlayer(&Blackjack.Player{})
			s.BlackjackTable.Players[i].NewBet(eachBet)
			wallet -= float64(eachBet)
		}

		s.BlackjackTable.InitialDeal()
		//s.BlackjackTable.ShowTableInfo()

		for i := 0; i < numOfPlayers; i++ {
			// Player's turn, check player whether is blackjack or not, scores = 21
			for playerHandIndex := 0; playerHandIndex < len(s.BlackjackTable.Players[i].Hands); playerHandIndex++ {
				for s.BlackjackTable.Players[i].CanMadeAction(playerHandIndex) {
					currentPlayerHand := s.BlackjackTable.Players[i].Hands[playerHandIndex]

					playerTotal := currentPlayerHand.Score()
					strategyType := GetPlayerHandStrategyType(currentPlayerHand)
					if (strategyType == Strategy.PairSplittingStrategy) && (len(currentPlayerHand.Cards) == 2) {
						playerTotal = currentPlayerHand.Cards[0].Rank
					}
					// Get the player's action based on the strategy
					playerAction, getActionErr := strategy.GetAction(playerTotal, s.BlackjackTable.Dealer.Cards[0].Rank, strategyType)
					if getActionErr != nil {
						fmt.Printf("Player's Hand %s\n", currentPlayerHand.String())
						fmt.Println("Error getting action: ", getActionErr)
						panic(getActionErr)
					}
					//fmt.Println("Player action: ", playerAction)

					if playerAction == Strategy.Double {
						//initialWallet -= currentPlayerHand.Bet
						currentPlayerHand.Bet *= 2
						s.BlackjackTable.DealToPlayer(i, playerHandIndex)
						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action:        Blackjack.DoubleDown,
							AdditionalBet: currentPlayerHand.Bet,
						},
							playerHandIndex,
						)

						wallet -= float64(currentPlayerHand.Bet)
					} else if playerAction == Strategy.Hit {
						s.BlackjackTable.DealToPlayer(i, playerHandIndex)
						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action: Blackjack.Hit,
						},
							playerHandIndex,
						)
					} else if playerAction == Strategy.Split {
						// s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
						// 	Action: Blackjack.Stand,
						// }, playerHandIndex)
						s.BlackjackTable.Players[i].Split()
						s.BlackjackTable.DealToPlayer(i, playerHandIndex)
						s.BlackjackTable.DealToPlayer(i, playerHandIndex+1)

						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action:        Blackjack.Split,
							AdditionalBet: currentPlayerHand.Bet,
						},
							playerHandIndex,
						)

						wallet -= float64(eachBet)

					} else if playerAction == Strategy.Stand {
						s.BlackjackTable.Players[i].AddNewAction(Blackjack.BlackJackAction{
							Action: Blackjack.Stand,
						}, playerHandIndex)
					}

				}
			}
		}

		// settle the dealer's hand
		for {
			s.BlackjackTable.DealToDealer()
			if s.BlackjackTable.Dealer.Score() >= 17 {
				break
			}
		}

		currentGameTotalBet := 0.0
		for i := 0; i < numOfPlayers; i++ {
			currentGameTotalBet += float64(s.BlackjackTable.Players[i].GetTotalBet())
		}

		ev := 0.0
		// settle the player's hand
		for i := 0; i < numOfPlayers; i++ {
			// if s.BlackjackTable.Players[i].Hands[0].HasPairsJackpot() {
			// 	numOfPairsJackpot++
			// 	wallet += float64(100) * 11
			// 	ev += float64(100) * 11
			// 	fmt.Println("Player ", i, " pairs jackpot", s.BlackjackTable.Players[i].Hands[0].String())
			// } else {
			// 	wallet -= float64(100)
			// 	ev -= float64(100)
			// }

			for _, playerHand := range s.BlackjackTable.Players[i].Hands {

				playerScore := playerHand.Score()
				if playerScore > 21 {
					//wallet -= float64(playerHand.Bet)
					ev -= float64(playerHand.Bet)
					//fmt.Println("Player ", i, " busted")
				} else if playerScore == 21 {
					if i == 0 && playerHand.IsBlackjack() {
						wallet += float64(playerHand.Bet)*1.5 + float64(playerHand.Bet)
						ev += float64(playerHand.Bet) * 2.5
					} else {
						wallet += float64(playerHand.Bet) + float64(playerHand.Bet)
						ev += float64(playerHand.Bet) * 2

					}

					//fmt.Println("Player ", i, " blackjack")
				} else if playerScore < 21 {
					dealerScore := s.BlackjackTable.Dealer.Score()
					if dealerScore > 21 {
						wallet += float64(playerHand.Bet) + float64(playerHand.Bet)
						ev += float64(playerHand.Bet) * 2
						//	fmt.Println("Player ", i, " win")
					} else if dealerScore == playerScore {
						//	fmt.Println("Player ", i, " tie")
						wallet += float64(playerHand.Bet)
						ev += float64(playerHand.Bet)
					} else if dealerScore > playerScore {
						//wallet -= float64(playerHand.Bet)
						ev -= float64(playerHand.Bet)

						//	fmt.Println("Player ", i, " lose")
					} else {
						wallet += float64(playerHand.Bet) * 2
						ev += float64(playerHand.Bet) * 2
						//	fmt.Println("Player ", i, " win")
					}
				}
			}

		}

		if ev >= currentGameTotalBet {
			numOfEvPositive++
		}

		if wallet > float64(highestReward) {
			highestReward = int(wallet)
		}

		if wallet <= 0 {
			break
		}
		currentGame++
	}

	fmt.Println("Number of decks: ", s.BlackjackTable.NumOfDecks)
	fmt.Println("Min bet: ", s.BlackjackTable.MinBet)
	fmt.Println("running num of games: ", currentGame)

	//allWallets := 0

	// for i := 0; i < numOfPlayers; i++ {
	// 	fmt.Println("Player ", i, " wallet: ", s.BlackjackTable.Players[i].Wallet)
	// 	allWallets += s.BlackjackTable.Players[i].Wallet
	// }
	fmt.Println("All wallet: ", wallet)

	fmt.Printf("Number of positive ev: %d\n", numOfEvPositive)
	fmt.Printf("Number of win percent: %.2f\n", float32(numOfEvPositive)/float32(currentGame)*100.0)
	fmt.Printf("Number of pairs jackpot: %d\n", numOfPairsJackpot)
	fmt.Printf("Highest reward: %d\n", highestReward)

	return wallet
}
