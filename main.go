package main

import (
	"fmt"
	"time"

	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func main() {
	// hp := yahtzee.NewHumanPlayer()
	seed := time.Now().Unix()
	defer func() {
		fmt.Println("Game over! Goodbye!")
		fmt.Printf("Seed was %d", seed)
	}()
	fmt.Printf("Playing a new game with seed: %d \n", seed)
	ai := yahtzee.NewAiPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{ai},
		Seed:    seed,
	}
	g.Play()
}
