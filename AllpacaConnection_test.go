package alpaca_game_ai

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func printGame(gamestate *Gamestate, turnCnt int) *GameAction {
	action := new(GameAction)

	samePlayableCards := filterSamePlayableCards(gamestate.Hand, gamestate.DiscardedCard)
	if len(samePlayableCards) > 0 {
		action.Action = "DROP CARD"
		action.Card = samePlayableCards[0].Name
		return action
	}

	playableCards := filterPlayableCards(gamestate.Hand, gamestate.DiscardedCard)
	if len(playableCards) > 0 {
		action.Action = "DROP CARD"
		action.Card = playableCards[0].Name
		return action
	}

	if calcHandcardValue(gamestate.Hand) < 5 {
		action.Action = "LEAVE ROUND"
		action.Card = ""
		return action
	}

	if gamestate.CardpileLeft > 0 && gamestate.PlayersLeft > 1 {
		action.Action = "DRAW CARD"
		action.Card = ""
		return action
	}

	action.Action = "LEAVE ROUND"
	action.Card = ""
	return action
}

func TestTestServer(t *testing.T) {

	end := make(chan bool)

	conn1 := AlpacaConnection{
		Url:            "http://localhost:3000",
		CallbackPort:   3001,
		CallbackHandle: "/call1",
		CallbackIP:     "localhost",
		CallbackFunc:   printGame,
		GameEnd:        end,
	}
	conn2 := AlpacaConnection{
		Url:            "http://localhost:3000",
		CallbackPort:   3001,
		CallbackHandle: "/call2",
		CallbackIP:     "localhost",
		CallbackFunc:   printGame,
		GameEnd:        end,
	}
	conn3 := AlpacaConnection{
		Url:            "http://localhost:3000",
		CallbackPort:   3001,
		CallbackHandle: "/call3",
		CallbackIP:     "localhost",
		CallbackFunc:   printGame,
		GameEnd:        end,
	}
	conn4 := AlpacaConnection{
		Url:            "http://localhost:3000",
		CallbackPort:   3001,
		CallbackHandle: "/call4",
		CallbackIP:     "localhost",
		CallbackFunc:   printGame,
		GameEnd:        end,
	}

	conn1.ReceiveCallbacks()
	conn2.ReceiveCallbacks()
	conn3.ReceiveCallbacks()
	conn4.ReceiveCallbacks()

	for i := 0; i < 5; i++ {

		conn1.Reset()
		conn2.Reset()
		conn3.Reset()
		conn4.Reset()

		svr := &http.Server{Addr: ":3001"}
		go func() {
			svr.ListenAndServe()
		}()

		conn1.Login("P1")
		conn2.Login("P2")
		conn3.Login("P3")
		conn4.Login("P4")

		fmt.Println(<-end)
		fmt.Println(conn1.LastScore, conn2.LastScore, conn3.LastScore, conn4.LastScore)

		svr.Shutdown(context.TODO())
	}

}
