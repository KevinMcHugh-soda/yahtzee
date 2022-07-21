package main

import (
	"fmt"
	"time"

	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func main() {
	// hp := yahtzee.NewHumanPlayer()
	seed := time.Now().Unix()
	fmt.Printf("Playing a new game with seed: %d", seed)
	ai := yahtzee.NewAiPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{ai},
		Seed:    seed,
	}
	g.Play()
	fmt.Println("Game over! Goodbye!")
	fmt.Printf("Seed was %d", seed)
}
