package alpaca_game_ai

import (
	"fmt"
	"github.com/yaricom/goNEAT/experiments"
	"github.com/yaricom/goNEAT/neat"
	"github.com/yaricom/goNEAT/neat/genetics"
	"os"
	"testing"
	"time"
)

func TestTrainAI(t *testing.T) {

	// Input Count
	// 7 Top Card
	// 4*7 für karten
	// 1 Cards Left
	// 1 Players Left
	// 1 Player LowestCardCount

	configFile, err := os.Open("context.neat")
	if err != nil {
		t.Error("Failed to load context", err)
		return
	}

	context := neat.LoadContext(configFile)
	neat.LogLevel = neat.LogLevelInfo

	pop, err := genetics.NewPopulationRandom(7+(4*7)+1+1+1, 1+1+7, 50, true, 0.15, context)

	evaluator := AlpacaGenerationEvaluator{
		OutputPath:  "None",
		PlayerCount: 4,
		selfPlay:    false,
		baselineFnc: BaseBot,
	}

	epoch_exec := genetics.SequentialPopulationEpochExecutor{}

	for generation_id := 0; generation_id < context.NumGenerations; generation_id++ {
		neat.InfoLog(fmt.Sprintf(">>>>> Generation:%3d\tRun: %d\n", generation_id, 0))
		generation := experiments.Generation{
			Id:      generation_id,
			TrialId: 0,
		}
		gen_start_time := time.Now()
		err = evaluator.GenerationEvaluate(pop, &generation, context)
		if err != nil {
			neat.InfoLog(fmt.Sprintf("!!!!! Generation [%d] evaluation failed !!!!!\n", generation_id))
			return
		}
		generation.Executed = time.Now()
		fmt.Println(generation.Average())
		fmt.Println(generation.Best.Fitness)

		// Turnover population of organisms to the next epoch if appropriate
		if !generation.Solved {
			neat.DebugLog(">>>>> start next generation")
			err = epoch_exec.NextEpoch(generation_id, pop, context)
			if err != nil {
				neat.InfoLog(fmt.Sprintf("!!!!! Epoch execution failed in generation [%d] !!!!!\n", generation_id))
				return
			}
		}

		// Set generation duration, which also includes preparation for the next epoch
		generation.Duration = generation.Executed.Sub(gen_start_time)
	}

}
