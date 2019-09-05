package alpaca_game_ai

import "fmt"

func filterPlayableCards(allCards []Cart, topCard Cart) []Cart  {

	filterCard := make([]Cart,0)

	for _,k := range allCards {
		if k.Type == topCard.Type || (k.Type+1 % 7) == topCard.Value  {
			fmt.Print(topCard)
			fmt.Print("-->")
			fmt.Println(k)
			filterCard = append(filterCard, k)
		}
	}


	return filterCard
}

func calcHandcardValue(allCards []Cart) int {

	countedCard := make([]Cart,0)
	sum := 0

	for _,k := range allCards {
		contains:= false
		for _,c := range countedCard {
			if  c.Type == k.Type {
				contains = true
					break
			}
		}
		if contains {break}

		sum += k.Value
		countedCard = append(countedCard, k)
	}

	return sum
}