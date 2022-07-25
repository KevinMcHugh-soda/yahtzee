package main

import (
	"bufio"
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
		runManyGames()
		return
	} else if len(os.Args) > 1 && os.Args[1] == "regress" {
		fileName := os.Args[2]
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		oldScores := make([]int, 0)
		newScores := make([]int, 0)
		for scanner.Scan() {
			text := scanner.Text()
			vals := strings.Split(text, ":")
			seedStr, oldScoreStr := vals[0], vals[1]
			seed, err := strconv.Atoi(seedStr)
			if err != nil {
				panic(err)
			}
			oldScore, err := strconv.Atoi(oldScoreStr)
			if err != nil {
				panic(err)
			}
			newScore := runGame(int64(seed))
			oldScores = append(oldScores, oldScore)
			newScores = append(newScores, newScore)
			fmt.Printf("%s|%3d|%3d|%4d\n", strconv.Itoa(seed)[:5], oldScore, newScore, newScore-oldScore)
		}
		fmt.Println("OLD SCORES:")
		printHistogram(oldScores)
		fmt.Println("---------------------------------------------------------------------------------------------------")
		fmt.Println("NEW SCORES:")
		printHistogram(newScores)
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
	score := runGame(seed)
	fmt.Println("score was:", score)
}

func runGame(seed int64) int {
	// defer func() {
	// 	fmt.Println("Game over! Goodbye!")
	// 	fmt.Printf("Seed was %d\n", seed)
	// }()
	// fmt.Printf("Playing a new game with seed: %d \n", seed)
	ai := yahtzee.NewAiPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{ai},
		Seed:    seed,
	}
	g.Play()
	w := g.Winner[0]
	// fmt.Println(w.GetScorecard().Print())
	return w.GetScorecard().Total()
}

func runManyGames() {
	count := 5000
	if len(os.Args) >= 2 {
		count, _ = strconv.Atoi(os.Args[2])
	}
	scores := make(map[int64]int, count)
	for idx := 0; idx < count; idx++ {
		seed := int64(rand.Int())
		if idx%(count/10) == 0 {
			fmt.Println(idx)
		}
		score := runGame(seed)
		scores[seed] = score
	}

	vals := make([]int, 0, 1000)
	for _, score := range scores {
		vals = append(vals, score)
	}
	printHistogram(vals)

	f, err := os.Create(fmt.Sprintf("%d.games", time.Now().Unix()))
	if err != nil {
		fmt.Println(err)
		return
	}

	for idx, score := range scores {
		str := fmt.Sprintf("%d:%d\n", idx, score)
		f.Write([]byte(str))
	}
}

func printHistogram(scores []int) {
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
}
