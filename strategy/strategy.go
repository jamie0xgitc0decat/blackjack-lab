package strategy

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Action string

const (
	Hit       Action = "H"
	Stand     Action = "S"
	Double    Action = "D"
	Split     Action = "SP"
	Surrender Action = "R"
)

type Strategy struct {
	HardTotals    map[int]map[int]Action // Maps player's total to dealer's card to action
	SoftTotals    map[int]map[int]Action // Similar structure can be used for soft totals
	PairSplitting map[int]map[int]Action // Maps pair value to dealer's card to action
}

type StrategyType string

const (
	HardTotalsStrategy    StrategyType = "HardTotals"
	SoftTotalsStrategy    StrategyType = "SoftTotals"
	PairSplittingStrategy StrategyType = "PairSplitting"
)

func LoadStrategy(filename string, strategyType StrategyType, strategy *Strategy) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Assume the first row is the dealer's upcard values (e.g., 2-10, A),
	// and subsequent rows are player totals or pairs with corresponding actions.
	for i, record := range records {
		if i == 0 {
			continue // Skip header row
		}

		// The first column in each row indicates the player's total or pair.
		playerVal := record[0]
		for j, actionStr := range record {
			if j == 0 {
				continue // Skip the player's total/pair column
			}

			//dealerCard, _ := strconv.Atoi(records[0][j]) // Convert dealer's card from string to int
			dealerCard := j + 1
			action := Action(actionStr)

			// Depending on the format (hard total, soft total, or pair), populate the appropriate map.
			// For simplicity, here we'll assume we're populating hard totals.
			// Similar logic can be used for soft totals and pair splitting by checking playerVal.
			if playerTotal, err := strconv.Atoi(playerVal); err == nil { // This is a hard total
				if strategyType == HardTotalsStrategy {
					if strategy.HardTotals[playerTotal] == nil {
						strategy.HardTotals[playerTotal] = make(map[int]Action)
					}
					strategy.HardTotals[playerTotal][dealerCard] = action
				} else if strategyType == SoftTotalsStrategy {
					if strategy.SoftTotals[playerTotal] == nil {
						strategy.SoftTotals[playerTotal] = make(map[int]Action)
					}
					strategy.SoftTotals[playerTotal][dealerCard] = action
				} else if strategyType == PairSplittingStrategy {
					if strategy.PairSplitting[playerTotal] == nil {
						strategy.PairSplitting[playerTotal] = make(map[int]Action)
					}
					strategy.PairSplitting[playerTotal][dealerCard] = action
				}

			} else { // This could be a soft total or pair, handle accordingly
				// Similar logic to populate SoftTotals or PairSplitting maps
				fmt.Printf("loadStrategy: failed to convert playerVal to int: %v\n", playerVal)
			}
		}
	}

	return nil
}

type LoadStrategyParams struct {
	HardTotalsFilePath string
	SoftTotalsFilePath string
	PairSplittingPath  string
}

// LoadStrategy loads the strategy from a given CSV file path.
func GenerateStrategy(params LoadStrategyParams) (Strategy, error) {

	strategy := Strategy{
		HardTotals:    make(map[int]map[int]Action),
		SoftTotals:    make(map[int]map[int]Action),
		PairSplitting: make(map[int]map[int]Action),
	}

	if err := LoadStrategy(params.HardTotalsFilePath, HardTotalsStrategy, &strategy); err != nil {
		return Strategy{}, err
	}
	if err := LoadStrategy(params.SoftTotalsFilePath, SoftTotalsStrategy, &strategy); err != nil {
		return Strategy{}, err
	}
	if err := LoadStrategy(params.PairSplittingPath, PairSplittingStrategy, &strategy); err != nil {
		return Strategy{}, err
	}
	return strategy, nil

}

// GetAction retrieves the recommended action based on the strategy tables.
// You would have separate functions or parameters to specify whether you're looking up a hard total, soft total, or pair.
func (s Strategy) GetAction(playerTotal, dealerUpCard int, stragetyType StrategyType) (Action, error) {
	// Simplified: this example assumes looking up a hard total strategy.
	// Implement similar logic for soft totals and pair splitting.
	if stragetyType == HardTotalsStrategy {

		if actions, exists := s.HardTotals[playerTotal]; exists {
			return actions[dealerUpCard], nil
		}
		if playerTotal > 17 {
			return Stand, nil
		}
		if playerTotal < 8 {
			return Hit, nil
		}
		return Action(""), fmt.Errorf("action not found for player total %d and dealer upcard %d and strategyType %s", playerTotal, dealerUpCard, stragetyType)
	} else if stragetyType == SoftTotalsStrategy {
		if actions, exists := s.SoftTotals[playerTotal]; exists {
			return actions[dealerUpCard], nil
		}
		if playerTotal < 13 {
			return Hit, nil
		}
		if playerTotal > 19 {
			return Stand, nil
		}
		return Action(""), fmt.Errorf("action not found for player total %d and dealer upcard %d and strategyType %s", playerTotal, dealerUpCard, stragetyType)

	} else {
		if actions, exists := s.PairSplitting[playerTotal]; exists {
			return actions[dealerUpCard], nil
		}
		return Action(""), fmt.Errorf("action not found for player total %d and dealer upcard %d and strategyType %s", playerTotal, dealerUpCard, stragetyType)
	}

}
