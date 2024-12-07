package main

import (
	Simulator "blackjack-lab/simulator"
)


func main() {

	simulator := Simulator.BlackjackSimulator{}
	//simulator.StartSimulation()
	simulator.StartOneMonthSimulation()

}
