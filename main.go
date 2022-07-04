package main

import (
	"fmt"

	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func main() {
	hp := yahtzee.NewHumanPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{hp},
	}
	g.Play()
	fmt.Println("Game over! Goodbye!")
}
