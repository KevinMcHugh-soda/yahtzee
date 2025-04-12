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

var hard = "🟨🟨🟨🟨🟩🟩🟩🟩🟦🟦\n🟨🟨🟨🟨🟩🟩🟩🟩🟦🟦\n🟨🟨🟥🟥🟩🟩🟩🟦🟦🟦\n🟨🟨🟥🟥🟩🟦🟦🟦🟦🟧\n⬛🟥🟥🟪🟪🟪🟪🟧🟧🟧\n⬛🟪🟪🟪🟪🟪🟪⬜⬜⬜\n⬛⬛⬛🟪🟪🟪🟪🟪⬜⬜\n⬛⬛⬛🟪🟫🟫🟪⬜⬜🟫\n⬛⬛🍺🍺🍺🟫🟫🟫🟫🟫\n⬛⬛⬛🍺🍺🍺🍺🍺🍺🍺"
var april11 = "🟧🟧🟧🟪🟪🟦🟦\n🟧🟩🟧⬜🟪🟥🟦\n🟩🟩🟧⬜🟥🟥🟦\n🟨🟩🟧🟧🟧🟥🟦\n🟨🟩🟧🟧🟧🟥🟦\n🟩🟩🟩🟧🟥🟥🟥\n🟧🟧🟧🟧🟧🟧🟧"

func main() {
	seed := time.Now().Unix()
	fmt.Println("hello", os.Args)
	if len(os.Args) > 1 && os.Args[1] == "star" {
		var p *yahtzee.Puzzle
		var err error

		if len(os.Args) > 2 && os.Args[2] == "hard" {
			p, err = yahtzee.ParsePuzzle(strings.Split(hard, "\n"), 2)
			if err != nil {
				fmt.Println("Error in parsing puzzle!", err)
			}
		} else if len(os.Args) > 2 && os.Args[2] == "411" {
			p, err = yahtzee.ParsePuzzle(strings.Split(april11, "\n"), 1)
			if err != nil {
				fmt.Println("Error in parsing puzzle!", err)
			}
		} else {
			y := yahtzee.MakeEasyPuzzle()
			p = &y
		}
		p2, err := p.Deduce()
		if err != nil {
			fmt.Println(err)
			return
		}
		p2.Print("state before solving, after deduction")
		x, solved := yahtzee.Solve(*p2)
		if solved {
			x.Print("solution!!!")
		} else {
			x.Print("last answer")
		}

		return
	} else if len(os.Args) > 1 && os.Args[1] == "mass" {
		runManyGames()
		return
	} else if len(os.Args) > 1 && os.Args[1] == "regress" {
		fileName := os.Args[2]
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		deltas := make([]int, 0, 1000)
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
			newScore, _ := runGame(int64(seed))
			oldScores = append(oldScores, oldScore)
			newScores = append(newScores, newScore)
			if len(seedStr) > 5 {
				seedStr = seedStr[:5]
			}
			fmt.Printf("%s|%3d|%3d|%4d\n", seedStr, oldScore, newScore, newScore-oldScore)
			deltas = append(deltas, newScore-oldScore)
		}
		sum := 0.0
		for _, d := range deltas {
			sum += float64(d)
		}
		fmt.Println("delta:", sum/float64(len(deltas)))
		fmt.Println("---------------------------------------------------------------------------------------------------")
		printComparativeHistogram(oldScores, newScores)
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
	_, g := runGame(seed)
	fmt.Println(g.Winner[0].GetScorecard().Print())
}

func runGame(seed int64) (int, yahtzee.Game) {
	defer func() {
		fmt.Println("Game over! Goodbye!")
		fmt.Printf("Seed was %d\n", seed)
	}()
	fmt.Printf("Playing a new game with seed: %d \n", seed)
	// ai := yahtzee.NewAiPlayer()
	me := yahtzee.NewHumanPlayer()
	g := yahtzee.Game{
		Players: []*yahtzee.Player{
			me,
			// ai,
		},
		Seed: seed,
	}
	g.Play()
	w := g.Winner[0]
	fmt.Println(w.GetScorecard().Print())
	return w.GetScorecard().Total(), g
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
		score, _ := runGame(seed)
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

func printComparativeHistogram(oldScores, newScores []int) {
	oldScoresByDecile := make(map[int]int)
	newScoresByDecile := make(map[int]int)
	maxDecile := 0

	for idx, oldScore := range oldScores {
		newScore := newScores[idx]

		if oldScore/10 > maxDecile {
			maxDecile = oldScore / 10
		}
		if newScore/10 > maxDecile {
			maxDecile = newScore / 10
		}

		oldScoresByDecile[oldScore/10] += 1
		newScoresByDecile[newScore/10] += 1
	}
	highestCount := 0
	for decile, oldCount := range oldScoresByDecile {
		newCount := newScoresByDecile[decile]

		if oldCount > highestCount {
			highestCount = oldCount
		}
		if newCount > highestCount {
			highestCount = newCount
		}
	}
	width := intFromEnvVarOr("WIDTH", 121)

	// "%3d|(%3d)%s|%s(%3d)\n"
	usableForHistogram := width - 3 - 1 - 1 - 3 - 1 - 1 - 1 - 3 - 1
	if usableForHistogram < 1 {
		panic("Width must be at least 16")
	}
	scalingFactor := 2 * float64(highestCount) / float64(usableForHistogram)
	fmt.Println("usable for histogram", usableForHistogram, "highest count", highestCount, "scaling factor", scalingFactor)
	if scalingFactor == 0 {
		scalingFactor = 1
	}
	colorReset := "\033[0m"

	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	for idx := 0; idx < maxDecile; idx++ {
		oldScoreCount := oldScoresByDecile[idx]
		newScoreCount := newScoresByDecile[idx]
		paddingNeeded := (usableForHistogram / 2) + 5
		oldScoreString := fmt.Sprintf("%s%*s%s", colorRed, paddingNeeded, strings.Repeat("=", int(float64(oldScoreCount)/scalingFactor)), colorReset)
		newScoreString := fmt.Sprintf("%s%-*s%s", colorGreen, paddingNeeded, strings.Repeat("=", int(float64(newScoreCount)/scalingFactor)), colorReset)
		fmt.Printf("%3d|(%3d)%s|%s(%3d)\n", idx*10, oldScoreCount, oldScoreString, newScoreString, newScoreCount)
	}
}

func intFromEnvVarOr(varName string, def int) int {
	valStr := os.Getenv(varName)
	val := def
	if len(valStr) > 0 {
		var err error
		val, err = strconv.Atoi(valStr)
		if err != nil {
			panic(err)
		}
	}

	return val
}
