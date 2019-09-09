package alpaca_game_ai

import (
	"fmt"
	"testing"
)

func TestTestRound(t *testing.T) {
	const playerCount = 4
	funcs := make([]TurnFunc, playerCount)
	for i := 0; i < playerCount; i++ {
		funcs[i] = BaseBot
	}

	control := AlpacaControl{
		Url:          "http://localhost:3000",
		CallbackIP:   "localhost",
		CallbackPort: 3001,
		PlayerCount:  playerCount,
		PlayerFunc:   funcs,
		TurnCount:    500,
	}

	control.Init()
	run := control.RunRound()
	fmt.Println(run)
}

func TestTestMultiRound(t *testing.T) {
	const playerCount = 4
	const roundCount = 1
	const roundLength = 500
	funcs := make([]TurnFunc, playerCount)
	for i := 0; i < playerCount; i++ {
		funcs[i] = BaseBot
	}
	funcs[0] = BadBot
	funcs[1] = BadBot
	control := AlpacaControl{
		Url:          "http://localhost:3000",
		CallbackIP:   "localhost",
		CallbackPort: 3001,
		PlayerCount:  playerCount,
		PlayerFunc:   funcs,
		TurnCount:    roundLength,
	}

	control.Init()
	sum := make([]int, playerCount)
	avg := make([]float64, playerCount)
	for i := 0; i < playerCount; i++ {
		sum[i] = 0
	}

	for i := 0; i < roundCount; i++ {
		run := control.RunRound()
		fmt.Println(run)
		for i := 0; i < playerCount; i++ {
			sum[i] += run[i]
		}
	}
	for i := 0; i < playerCount; i++ {
		avg[i] = float64(sum[i]) / float64(roundCount)
	}

	fmt.Println("Average:")
	fmt.Println(sum)

}
