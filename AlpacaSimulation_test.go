package alpaca_game_ai

import (
	"fmt"
	"testing"
)

func TestTestMultiSim(t *testing.T) {
	sim := NewAlpacaSimulation()

	sim.AddPlayer("P1", BaseBot)
	sim.AddPlayer("P2", BaseBot)
	sim.AddPlayer("P3", BadBot)
	sim.AddPlayer("P4", BadBot)

	fmt.Println(sim.RunSimulation(500))
}
