package alpaca_game_ai

import (
	"math/rand"
	"time"
)

type AlpacaSimulation struct {
	players          []PlayerStat
	playerCards      [][]Cart
	callbacks        []TurnFunc
	cardPile         []Cart
	discardedPile    []Cart
	topDiscardedCard Cart
	nextPlayer       int
	over             bool
	Seed             int64
}

func NewAlpacaSimulation() AlpacaSimulation {
	return AlpacaSimulation{
		players:     make([]PlayerStat, 0),
		playerCards: make([][]Cart, 0),
		callbacks:   make([]TurnFunc, 0),
		Seed:        0,
	}
}

func (sim *AlpacaSimulation) AddPlayer(name string, turnFunc TurnFunc) {
	sim.players = append(sim.players, PlayerStat{
		PlayerName: name,
		CardCount:  0,
		Coins:      make([]Coin, 0),
		Score:      0,
		InGame:     false,
	})
	sim.callbacks = append(sim.callbacks, turnFunc)
}

func (sim *AlpacaSimulation) AddPlayers(names []string, turnFunctions []TurnFunc) {

	if len(names) == len(turnFunctions) {
		for i := 0; i < len(names); i++ {
			sim.AddPlayer(names[i], turnFunctions[i])
		}
	}
}

func (sim *AlpacaSimulation) RunSimulation(rounds int) []int {
	if sim.players == nil || sim.playerCards == nil || sim.callbacks == nil {
		panic("AlpacaSimulation not Init")
	}
	if len(sim.players) == 0 || len(sim.callbacks) == 0 {
		panic("No Player Added")
	}

	if sim.Seed == 0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(sim.Seed)
	}
	for i := 0; i < rounds; i++ {
		sim.round()
	}

	result := make([]int, len(sim.players))

	for i := 0; i < len(sim.players); i++ {
		result[i] = sim.players[i].Score
	}

	return result
}

func (sim *AlpacaSimulation) round() []int {
	sim.over = false
	sim.nextPlayer = 0
	sim.cardPile = make([]Cart, 0)
	sim.discardedPile = make([]Cart, 0)

	for i := 0; i < 8; i++ {
		for j := 1; j < 7; j++ {
			name := ""
			if j == 1 {
				name = "ONE"
			} else if j == 2 {
				name = "TWO"
			} else if j == 3 {
				name = "THREE"
			} else if j == 4 {
				name = "FOUR"
			} else if j == 5 {
				name = "FIVE"
			} else if j == 6 {
				name = "SIX"
			}

			sim.cardPile = append(sim.cardPile, Cart{
				Type:  j,
				Name:  name,
				Value: j,
			})
		}

		sim.cardPile = append(sim.cardPile, Cart{
			Type:  0,
			Name:  "ALPACA",
			Value: 10})

	}

	rand.Shuffle(len(sim.cardPile), func(i, j int) {
		sim.cardPile[i], sim.cardPile[j] = sim.cardPile[j], sim.cardPile[i]
	})

	sim.playerCards = make([][]Cart, len(sim.players))
	for i := 0; i < len(sim.players); i++ {
		sim.playerCards[i] = make([]Cart, 6)
		sim.players[i].CardCount = 6
		sim.players[i].InGame = true
	}

	for j := 0; j < 6; j++ {
		for i := 0; i < len(sim.players); i++ {
			//POP
			card := sim.cardPile[len(sim.cardPile)-1]
			sim.cardPile = sim.cardPile[:len(sim.cardPile)-1]
			sim.playerCards[i][j] = card
		}
	}

	//POP
	sim.topDiscardedCard = sim.cardPile[len(sim.cardPile)-1]
	sim.discardedPile = append(sim.discardedPile, sim.topDiscardedCard)
	sim.cardPile = sim.cardPile[:len(sim.cardPile)-1]

	for i := 0; sim.over == false; i++ {
		game := sim.getGame(sim.nextPlayer)
		if game.PlayersLeft == 0 {
			sim.over = true
			break
		}

		action := sim.callbacks[sim.nextPlayer](&game)
		sim.runAction(*action)
	}

	//Round over
	result := make([]int, len(sim.players))
	for i := 0; i < len(sim.players); i++ {
		//POP
		value := calcHandcardValue(sim.playerCards[i])
		if value == 0 {
			if sim.players[i].Score >= 10 {
				sim.players[i].Score -= 10
			} else if sim.players[i].Score >= 1 {
				sim.players[i].Score -= 1
			}
		} else {
			sim.players[i].Score += value
		}
		result[i] = value
	}

	return result
}

func (sim *AlpacaSimulation) getGame(playeridx int) Gamestate {

	game := Gamestate{
		MyTurn:        playeridx == sim.nextPlayer && sim.players[playeridx].InGame,
		OtherPlayers:  make([]map[string]PlayerStat, 0),
		Hand:          sim.playerCards[playeridx],
		DiscardedCard: sim.topDiscardedCard,
		Score:         sim.players[playeridx].Score,
		Coins:         sim.players[playeridx].Coins,
		CardpileLeft:  len(sim.cardPile),
		PlayersLeft:   0,
	}

	for i := 0; i < len(sim.players); i++ {
		if sim.players[i].InGame {
			game.PlayersLeft++
		}

		if playeridx != i {
			game.OtherPlayers = append(game.OtherPlayers, map[string]PlayerStat{sim.players[i].PlayerName: sim.players[i]})
		}
	}

	return game
}

func (sim *AlpacaSimulation) runAction(action GameAction) {

	if sim.players[sim.nextPlayer].InGame {
		if action.Action == "DROP CARD" {

			if action.Card == "ONE" {
				if 1 == sim.topDiscardedCard.Type || 0 == sim.topDiscardedCard.Type {
					card := Cart{
						Type:  1,
						Name:  "ONE",
						Value: 1,
					}
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.discardedPile = append(sim.discardedPile, card)
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			} else if action.Card == "TWO" {
				card := Cart{
					Type:  2,
					Name:  "TWO",
					Value: 2,
				}
				if 2 == sim.topDiscardedCard.Type || 1 == sim.topDiscardedCard.Type {
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			} else if action.Card == "THREE" {
				card := Cart{
					Type:  3,
					Name:  "THREE",
					Value: 3,
				}
				if 3 == sim.topDiscardedCard.Type || 2 == sim.topDiscardedCard.Type {
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			} else if action.Card == "FOUR" {
				card := Cart{
					Type:  4,
					Name:  "FOUR",
					Value: 4,
				}
				if 4 == sim.topDiscardedCard.Type || 3 == sim.topDiscardedCard.Type {
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			} else if action.Card == "FIVE" {
				card := Cart{
					Type:  5,
					Name:  "FIVE",
					Value: 5,
				}
				if 5 == sim.topDiscardedCard.Type || 4 == sim.topDiscardedCard.Type {
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			} else if action.Card == "SIX" {
				card := Cart{
					Type:  6,
					Name:  "SIX",
					Value: 6,
				}
				if 6 == sim.topDiscardedCard.Type || 5 == sim.topDiscardedCard.Type {
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			} else if action.Card == "ALPACA" {
				card := Cart{
					Type:  0,
					Name:  "ALPACA",
					Value: 10,
				}
				if 0 == sim.topDiscardedCard.Type || 6 == sim.topDiscardedCard.Type {
					sim.playerCards[sim.nextPlayer] = removeCard(sim.playerCards[sim.nextPlayer], card)
					sim.topDiscardedCard = card
					sim.players[sim.nextPlayer].CardCount--
				} else {
					panic("INVALID CARD")
				}
			}

			if len(sim.playerCards[sim.nextPlayer]) == 0 {
				sim.over = true
			}

		} else if action.Action == "DRAW CARD" {

			//POP
			card := sim.cardPile[len(sim.cardPile)-1]
			sim.cardPile = sim.cardPile[:len(sim.cardPile)-1]
			sim.playerCards[sim.nextPlayer] = append(sim.playerCards[sim.nextPlayer], card)
			sim.players[sim.nextPlayer].CardCount++

		} else if action.Action == "LEAVE ROUND" {
			sim.players[sim.nextPlayer].InGame = false
		}
	}

	sim.nextPlayer = (sim.nextPlayer + 1) % (len(sim.players))

}

func removeCard(cards []Cart, rm Cart) []Cart {
	idx := -1

	for i, k := range cards {
		if k.Type == rm.Type {
			idx = i
			break
		}
	}

	if idx >= 0 {
		//Remove Card
		cards[idx] = cards[len(cards)-1]
		return cards[:len(cards)-1]
	}
	return cards
}
