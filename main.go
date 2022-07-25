package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func main() {
	// hp := yahtzee.NewHumanPlayer()
	// TODO: build a harness to generate pair/score seeds and save them,
	// and also to run from those seeds
	seed := time.Now().Unix()
	if len(os.Args) > 1 {
		fmt.Println(os.Args)
		arg := os.Args[1]
		seed64, err := strconv.Atoi(arg)
		seed = int64(seed64)
		if err != nil {
			panic(err)
		}
	}
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
