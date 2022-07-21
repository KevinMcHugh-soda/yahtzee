package main

import (
	"fmt"

	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func main() {
	// hp := yahtzee.NewHumanPlayer()
	ai := yahtzee.NewAiPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{ai},
	}
	g.Play()
	fmt.Println("Game over! Goodbye!")
}
