package AlpacaGameAI

func filterPlayableCards(allCards []Cart, topCard Cart) []Cart {

	filterCard := make([]Cart, 0)

	for _, k := range allCards {
		if k.Type == topCard.Type || k.Type == (topCard.Type+1)%7 {
			filterCard = append(filterCard, k)
		}
	}

	return filterCard
}

func filterSamePlayableCards(allCards []Cart, topCard Cart) []Cart {

	filterCard := make([]Cart, 0)

	for _, k := range allCards {
		if k.Type == topCard.Type {
			filterCard = append(filterCard, k)
		}
	}

	return filterCard
}

func calcHandcardValue(allCards []Cart) int {

	countedCard := make([]Cart, 0)
	sum := 0

	for _, k := range allCards {
		contains := false
		for _, c := range countedCard {
			if c.Type == k.Type {
				contains = true
				break
			}
		}
		if contains {
			break
		}

		sum += k.Value
		countedCard = append(countedCard, k)
	}

	return sum
}

func BaseBot(gamestate *Gamestate) *GameAction {
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

func BadBot(gamestate *Gamestate) *GameAction {
	action := new(GameAction)

	playableCards := filterPlayableCards(gamestate.Hand, gamestate.DiscardedCard)
	if len(playableCards) > 0 {
		action.Action = "DROP CARD"
		action.Card = playableCards[0].Name
		return action
	}

	if calcHandcardValue(gamestate.Hand) < 3 {
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
