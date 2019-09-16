package main

import (
	"alpaca-game-ai/AlpacaGameAI"
	"flag"
	"fmt"
	"github.com/yaricom/goNEAT/neat/genetics"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	test := flag.Bool("test", false, "testing without a server")
	url := flag.String("url", "localhost:3000", "Url To the Alpaca Server")
	myIp := flag.String("myip", "192.168.1.25", "IP-Address of this Computer")
	playerCount := flag.Int("pCnt", 4, "The Number of players that are playing")
	base := flag.Bool("baseline", false, "Use The Baseline AI")
	gen := flag.String("gen", "best_winner.gen", "Genome File for building the Network")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	botFunction := AlpacaGameAI.BaseBot

	if !(*base) {
		genomFile, err := os.Open(*gen)
		if err != nil {
			log.Fatal("Failed to load genome ", err)
			return
		}
		genome, err := genetics.ReadGenome(genomFile, 0)
		ex := AlpacaGameAI.AlpacaGenerationEvaluator{
			OutputPath:  "None",
			PlayerCount: *playerCount,
		}
		net, err := genome.Genesis(0)
		botFunction = ex.MakeNetworkRunFunc(net, new(int))
	}
	if !(*test) {
		control := AlpacaGameAI.AlpacaControl{
			Url:          *url,
			CallbackIP:   *myIp,
			CallbackPort: 3001 + rand.Intn(999),
			PlayerCount:  1,
			PlayerFunc:   []AlpacaGameAI.TurnFunc{botFunction},
			TurnCount:    999999999,
		}
		control.Init()
		control.RunRound()
	} else {

		sim := AlpacaGameAI.NewAlpacaSimulation()

		for i := 0; i < (*playerCount)-1; i++ {
			sim.AddPlayer("Base"+strconv.Itoa(i), AlpacaGameAI.BaseBot)
		}
		sim.AddPlayer("BOT", botFunction)

		fmt.Println("Running Test")
		result := sim.RunSimulation(50000)
		fmt.Println("Result:")
		fmt.Println(result)
	}
}
