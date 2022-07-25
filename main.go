package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func main() {
	seed := time.Now().Unix()
	if len(os.Args) > 1 && os.Args[1] == "mass" {
		// If a scoremap file doesnt
		count := 5000
		if len(os.Args) >= 2 {
			count, _ = strconv.Atoi(os.Args[2])
		}
		scores := make([]int, count)
		for idx := 0; idx < count; idx++ {
			seed := int64(rand.Int())
			if idx%(count/10) == 0 {
				fmt.Println(idx)
			}
			score := runGame(seed)
			scores[idx] = score
		}
		scoresByDecile := make(map[int]int)
		maxDecile := 0
		for _, score := range scores {
			if score/10 > maxDecile {
				maxDecile = score / 10
			}
			scoresByDecile[score/10] += 1
		}
		for idx := 0; idx < maxDecile; idx++ {
			fmt.Printf("%3d,%4d|%s\n", idx*10, scoresByDecile[idx], strings.Repeat("=", scoresByDecile[idx]/5))
		}
		return
	} else if len(os.Args) > 1 {
		fmt.Println(os.Args)
		arg := os.Args[1]
		seed64, err := strconv.Atoi(arg)
		seed = int64(seed64)
		if err != nil {
			panic(err)
		}
	}
	runGame(seed)
	// TODO: build a harness to generate pair/score seeds and save them,
	// and also to run from those seeds
}

func runGame(seed int64) int {
	// defer func() {
	// 	fmt.Println("Game over! Goodbye!")
	// 	fmt.Printf("Seed was %d", seed)
	// }()
	// fmt.Printf("Playing a new game with seed: %d \n", seed)
	ai := yahtzee.NewAiPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{ai},
		Seed:    seed,
	}
	g.Play()
	return g.Winner[0].GetScorecard().Total()
}
