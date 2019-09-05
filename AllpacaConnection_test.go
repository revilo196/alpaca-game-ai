package alpaca_game_ai

import (
	"fmt"
	"testing"
)

func printGame(gamestate * Gamestate) GameAction  {
	fmt.Println(gamestate)

	playableCards := filterPlayableCards(gamestate.Hand, gamestate.DiscardedCard)
	if len(playableCards) > 0 {
		return GameAction{
			Action: "DROP CARD",
			Card:   playableCards[0].Name,
		}
	}

	if calcHandcardValue(gamestate.Hand) < 3 {
		return GameAction{
			Action: "LEAVE ROUND",
			Card:   "",
		}
	}

	if gamestate.CardpileLeft > 0  {
		return GameAction{
			Action: "DRAW CARD",
			Card:   "",
		}
	}

	return GameAction{
		Action: "LEAVE ROUND",
		Card:   "",
	}
}

func TestTestServer(t *testing.T) {

	conn1 := AlpacaConnection{
		Url:        "http://localhost:3000",
		CallbackPort: 3001,
		CallbackHandle: "/call1",
		CallbackIP: "localhost",
		CallbackFunc:printGame,
	}
	conn2 := AlpacaConnection{
		Url:        "http://localhost:3000",
		CallbackPort: 3002,
		CallbackHandle: "/call2",
		CallbackIP: "localhost",
		CallbackFunc:printGame,
	}
	finished1 := make(chan bool)
	go conn1.ReceiveCallbacks(finished1)
	conn1.Login("P1")
	finished2 := make(chan bool)
	go conn2.ReceiveCallbacks(finished2)
	conn2.Login("P2")
	<- finished1
	<- finished2

}