package alpaca_game_ai

import (
	"fmt"
	"github.com/yaricom/goNEAT/experiments"
	"github.com/yaricom/goNEAT/neat"
	"github.com/yaricom/goNEAT/neat/genetics"
	"github.com/yaricom/goNEAT/neat/utils"
	"os"
	"testing"
)

func TestTrainAI(t *testing.T) {
	out_dir_path := "Out4"
	// Check if output dir exists
	if _, err := os.Stat(out_dir_path); err == nil {
		// clear it
		os.RemoveAll(out_dir_path)
	}
	// create output dir
	err := os.MkdirAll(out_dir_path, os.ModePerm)
	if err != nil {
		t.Errorf("Failed to create output directory, reason: %s", err)
		return
	}

	configFile, err := os.Open("context.neat")
	if err != nil {
		t.Error("Failed to load context", err)
		return
	}

	/*genomFile, err := os.Open("pole1_winner_200-930")
	if err != nil {
		t.Error("Failed to load context", err)
		return
	}*/

	context := neat.LoadContext(configFile)
	neat.LogLevel = neat.LogLevelDebug
	context.NodeActivatorsProb[0] = 0.25
	context.NodeActivators[0] = utils.SigmoidBipolarActivation

	context.NodeActivators = append(context.NodeActivators, utils.GaussianBipolarActivation)
	context.NodeActivatorsProb = append(context.NodeActivatorsProb, 0.25)

	context.NodeActivators = append(context.NodeActivators, utils.TanhActivation)
	context.NodeActivatorsProb = append(context.NodeActivatorsProb, 0.35)

	context.NodeActivators = append(context.NodeActivators, utils.LinearAbsActivation)
	context.NodeActivatorsProb = append(context.NodeActivatorsProb, 0.15)

	neat.LogLevel = neat.LogLevelInfo

	startGenome := genetics.NewGenomeRand(0, 7+(4*7)+1+1+1, 1+1+2, 10, 20, false, 0.25)
	//startGenome, err := genetics.ReadGenome(genomFile,0)

	//pop, err := genetics.NewPopulationRandom(7+(4*7)+1+1+1, 1+1+7, 50, true, 0.15, context)

	experiment := experiments.Experiment{
		Id:     0,
		Trials: make(experiments.Trials, context.NumRuns),
	}

	evaluator := AlpacaGenerationEvaluator{
		OutputPath:    "Out4",
		PlayerCount:   4,
		selfPlay:      true,
		selfCombiPlay: true,
		rounds:        1200,
		baselineFnc:   BaseBot,
	}

	err = experiment.Execute(context, startGenome, evaluator)
	//avg_nodes, avg_genes, avg_evals, _ := experiment.AvgWinner()
	fmt.Println(experiment.BestFitness())
	/*

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
	*/
}
