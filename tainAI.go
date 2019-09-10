package alpaca_game_ai

import (
	"fmt"
	"github.com/yaricom/goNEAT/experiments"
	"github.com/yaricom/goNEAT/neat"
	"github.com/yaricom/goNEAT/neat/genetics"
	"github.com/yaricom/goNEAT/neat/network"
	"time"

	"os"
	"strconv"
)

func (ev AlpacaGenerationEvaluator) gamestateToSensors(game Gamestate) []float64 {
	res := make([]float64, 7+(4*7)+1+1+1)

	//Top Card One Hot 7
	res[game.DiscardedCard.Type] = 1.0

	//Hand Cards One Hot 7*4
	for _, k := range game.Hand {
		i2 := 0
		for j := 0; j < 4; j++ {
			if res[7+(k.Type*4)+j] < 0.5 {
				i2 = j
				break
			}
		}
		res[7+(k.Type*4)+i2] = 1.0

	}

	res[7+(4*7)] = float64(game.PlayersLeft) / float64(ev.PlayerCount)
	res[7+(4*7)+1] = float64(game.CardpileLeft) / 56.0
	minCardCount := 6
	for _, k := range game.OtherPlayers {
		for _, v := range k {
			if minCardCount < v.CardCount {
				minCardCount = v.CardCount
			}
		}
	}

	res[7+(4*7)+1+1] = float64(minCardCount) / 6.0

	return res
}

func (ev AlpacaGenerationEvaluator) outputToAction(gamestate Gamestate, out []float64) GameAction {

	action := GameAction{Action: "LEAVE ROUND"}

	for j := 0; j < 9; j++ {

		max := 0.0
		idx := 0
		for i, v := range out {
			if max < v {
				max = v
				idx = i
			}
		}

		if idx < 7 {

			contains := false
			name := ""
			for _, v := range gamestate.Hand {
				if v.Type == idx {
					contains = true
					name = v.Name
					break
				}
			}

			if contains && (idx == gamestate.DiscardedCard.Type || idx == (gamestate.DiscardedCard.Type+1)%7) {
				action.Action = "DROP CARD"
				action.Card = name
			} else {
				//INVALID TURN
				out[idx] = 0.0
				//TRY AGAIN
			}

		} else if idx == 7 {
			//DRAW
			if gamestate.PlayersLeft == 1 || gamestate.CardpileLeft == 0 {
				//INVALID TURN
				out[idx] = 0.0
				//TRY AGAIN
			} else {
				action.Action = "DRAW CARD"
				break
			}
		} else if idx == 8 {
			//LEAVE
			action.Action = "LEAVE ROUND"
			break
		}
	}

	return action
}

type AlpacaGenerationEvaluator struct {
	OutputPath  string
	PlayerCount int
	selfPlay    bool
	baselineFnc TurnFunc
	seed        int64
	best        float64
}

func (ex AlpacaGenerationEvaluator) GenerationEvaluate(pop *genetics.Population, epoch *experiments.Generation, context *neat.NeatContext) (err error) {
	const CORES = 4
	fin := make(chan bool)
	ex.seed = time.Now().UnixNano()
	if !ex.selfPlay {
		for i, org := range pop.Organisms {
			if i < CORES {
				go ex.orgEvaluate(org, fin) //winner
			} else {
				<-fin
				go ex.orgEvaluate(org, fin)
			}
		}
		for i := 0; i < CORES; i++ {
			<-fin
		}
	} else {
		for i := 0; i < len(pop.Organisms); i += ex.PlayerCount {
			var orgs []*genetics.Organism

			//Unterteilen in PlayerCount Teile
			if i+ex.PlayerCount < len(pop.Organisms) {
				orgs = pop.Organisms[i : i+ex.PlayerCount]
			} else {
				orgs = pop.Organisms[i:]
			}

			if i/4 < CORES {
				go ex.orgsEvaluate(orgs, fin) //winner
			} else {
				<-fin
				go ex.orgsEvaluate(orgs, fin)
			}
		}
		for i := 0; i < CORES; i++ {
			<-fin
		}
	}

	epoch.FillPopulationStatistics(pop)

	// Only print to file every print_every generations
	if epoch.Solved || epoch.Id%context.PrintEvery == 0 {
		pop_path := fmt.Sprintf("%s/gen_%d", experiments.OutDirForTrial(ex.OutputPath, epoch.TrialId), epoch.Id)
		file, err := os.Create(pop_path)
		if err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
		} else {
			pop.WriteBySpecies(file)
		}
	}

	// print best organisms
	for _, org := range pop.Organisms {
		if org.Fitness > ex.best {
			ex.best = org.Fitness
			// Prints the winner organism to file!
			org_path := fmt.Sprintf("%s/%s_%d-%d", experiments.OutDirForTrial(ex.OutputPath, epoch.TrialId),
				"pole1_winner", org.Phenotype.NodeCount(), org.Phenotype.LinkCount())
			file, err := os.Create(org_path)
			if err != nil {
				neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism genome, reason: %s\n", err))
			} else {
				org.Genotype.Write(file)
				neat.InfoLog(fmt.Sprintf("Generation #%d winner dumped to: %s\n", epoch.Id, org_path))
			}
			break
		}
	}

	return err
}

func (ex *AlpacaGenerationEvaluator) orgEvaluate(organism *genetics.Organism, fin chan bool) (isWinner bool) {

	result := ex.runGame(organism.Phenotype)
	const BADEST_GAME = 10000.0 //?
	organism.Error = float64(result) / BADEST_GAME
	organism.Fitness = 1.0 - organism.Error
	fin <- true
	return false // IsWinner ...
}

func (ex *AlpacaGenerationEvaluator) runGame(net *network.Network) (score int) {
	sim := NewAlpacaSimulation()
	sim.Seed = ex.seed
	for i := 1; i < ex.PlayerCount; i++ {
		sim.AddPlayer("P"+strconv.Itoa(i), ex.baselineFnc)
	}

	sim.AddPlayer("EvoBot", ex.makeNetworkRunFunc(net))

	result := sim.RunSimulation(1000)
	//mfmt.Println(result)
	return result[ex.PlayerCount-1]
}

func (ex *AlpacaGenerationEvaluator) orgsEvaluate(organisms []*genetics.Organism, fin chan bool) (isWinner bool) {
	const BADEST_GAME = 10000.0 //?

	nets := make([]*network.Network, len(organisms))

	//Prepare Networks
	for i := 0; i < len(organisms); i++ {
		nets[i] = organisms[i].Phenotype
	}

	//Run The Game
	result := ex.runGameMult(nets)

	//Evaluatle Results
	for i := 0; i < len(result); i++ {
		organisms[i].Error = float64(result[i]) / BADEST_GAME
		organisms[i].Fitness = 1.0 - organisms[i].Error
	}

	fin <- true
	return false // IsWinner ...
}

func (ex *AlpacaGenerationEvaluator) makeNetworkRunFunc(net *network.Network) func(gamestate *Gamestate) *GameAction {
	return func(gamestate *Gamestate) *GameAction {
		err := net.LoadSensors(ex.gamestateToSensors(*gamestate))
		if err != nil {
			panic(err)
		}

		if res, err := net.Activate(); !res {
			//If it loops, exit returning only fitness of 1 step
			neat.DebugLog(fmt.Sprintf("Failed to activate Network, reason: %s", err))
			return nil
		}

		out := net.ReadOutputs()
		action := ex.outputToAction(*gamestate, out)
		return &action
	}
}

func (ex *AlpacaGenerationEvaluator) runGameMult(nets []*network.Network) (score []int) {
	sim := NewAlpacaSimulation()
	sim.Seed = ex.seed

	for i := 0; i < ex.PlayerCount; i++ {
		sim.AddPlayer("EvoBot"+strconv.Itoa(i), ex.makeNetworkRunFunc(nets[i]))
	}

	result := sim.RunSimulation(1000)
	//mfmt.Println(result)
	return result

}
