package alpaca_game_ai

import (
	"context"
	"net/http"
	"strconv"
)

type TurnFunc func(gamestate *Gamestate) *GameAction

type AlpacaControl struct {
	Url          string
	CallbackIP   string
	CallbackPort int
	PlayerCount  int
	PlayerFunc   []TurnFunc
	TurnCount    int
	connections  []AlpacaConnection
	roundEnd     chan bool
}

func (a *AlpacaControl) Init() {
	conn := make([]AlpacaConnection, a.PlayerCount)
	end := make(chan bool)

	makeFnc := func(idx int) func(gamestate *Gamestate, turnNr int) *GameAction {
		return func(gamestate *Gamestate, turnNr int) *GameAction {
			if turnNr > a.TurnCount {
				_, _ = http.Get(a.Url + "/reset")
			}
			//if turnNr % 100 == 0 {
			//}
			return a.PlayerFunc[idx](gamestate)
		}
	}

	for i := 0; i < a.PlayerCount; i++ {
		conn[i].Url = a.Url
		conn[i].CallbackPort = a.CallbackPort
		conn[i].CallbackIP = a.CallbackIP
		conn[i].CallbackHandle = "/call" + strconv.Itoa(i)

		conn[i].CallbackFunc = makeFnc(i)
		conn[i].GameEnd = end
		conn[i].ReceiveCallbacks()

	}
	a.connections = conn
	a.roundEnd = end
}

func (a *AlpacaControl) RunRound() []int {
	for i := 0; i < a.PlayerCount; i++ {
		a.connections[i].Reset()
	}

	svr := &http.Server{Addr: ":" + strconv.Itoa(a.CallbackPort)}
	go func() {
		svr.ListenAndServe()
	}()

	for i := 0; i < a.PlayerCount; i++ {
		a.connections[i].Login("P" + strconv.Itoa(i))
	}
	<-a.roundEnd
	svr.Shutdown(context.TODO())

	score := make([]int, a.PlayerCount)
	for i := 0; i < a.PlayerCount; i++ {
		score[i] = a.connections[i].LastScore
	}
	return score
}
