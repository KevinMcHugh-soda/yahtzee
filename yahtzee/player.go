package yahtzee

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Player interface {
	GetName() string
	GetScorecard() *Scorecard
	AssessRoll(hand Hand, rollsRemaining int) RollDecision
	PickScorable(hand Hand) Scoreable
}

type HumanPlayer struct {
	Scorecard *Scorecard
}

func (p HumanPlayer) GetName() string {
	return "Mr. Human!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
}

func (p HumanPlayer) GetScorecard() *Scorecard {
	return p.Scorecard
}

// TODO would be good to indicate roll no./rolls remaining
func (p HumanPlayer) AssessRoll(hand Hand, rollsRemaining int) RollDecision {
	fmt.Printf("Roll: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	// fmt.Println("Type y to keep, space to reroll:")
	reader := bufio.NewReader(os.Stdin)
	allInts := regexp.MustCompile("[1-6]{1,5}")
	var text string

	bools := make([]bool, 5)

OUTER:
	for {
		fmt.Println("please enter the values you want to keep.")
		text, _ = reader.ReadString('\n')
		if text != "\n" && !allInts.MatchString(text) {
			fmt.Println("only enter numbers between 1 and 6")
			continue OUTER
		}
		selectedValues := make([]int, 0)
		for idx, c := range text {
			if c == 10 {
				break
			}
			if int(c) < 49 || int(c) > 54 {
				// try again
				fmt.Println("all values must be between 1 and 6", int(c), c, string(c), idx)
				continue OUTER
			}
			selectedValues = append(selectedValues, int(c-'0'))

			// if idx < len(bools) {
			// 	bools[idx] = c == rune('y')
			// }
		}
		fmt.Println("values selected:", selectedValues)
		takenDieIndex := make(map[int]bool, 5)
		for _, value := range selectedValues {
			valueTaken := false
			for idx, die := range hand {
				if !takenDieIndex[idx] && !valueTaken {
					if die == value {
						valueTaken = true
						bools[idx] = true
						takenDieIndex[idx] = true
					}
				}
			}

			if !valueTaken {
				fmt.Println("only specify values you have, please, not", value)
				continue OUTER
			}
		}

		break
	}
	return RollDecision(bools)
}

func (p HumanPlayer) EnsureValidResponse(prompt string, isValid func(string) bool) string {
	var input string
	for {
		fmt.Println(prompt)
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if isValid(input) {
			break
		}
	}

	return input
}

func (p HumanPlayer) PickScorable(hand Hand) Scoreable {
	fmt.Printf("Hand: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	prompt := "Choose a row to score this roll\n"
	options := make(map[int]ScorableName)
	promptForName := make(map[ScorableName]string)
	for idx, name := range ScorableNames {
		if p.Scorecard.NameToScorePtr(name) == nil {
			options[idx+1] = name
			score := ScoreableByName(name).Score(hand, p.Scorecard.HadYahztee())
			promptForName[name] = fmt.Sprintf("(%2d points) [%d] to score %s;", score, idx+1, name)
		}
	}
	prompt += p.Scorecard.PrintWithDecorator(func(name ScorableName) string {
		return promptForName[name]
	})
	input := p.EnsureValidResponse(prompt, func(input string) bool {
		val, err := strconv.Atoi(input)

		return err == nil && val > 0 && val <= len(ScorableNames) && options[val] != ""
	})
	choice, _ := strconv.Atoi(input)

	return ScoreableByName(options[choice])
}

func NewHumanPlayer() *Player {
	z := 0
	scoreCard := Scorecard{YahtzeeBonusName: &z}
	hp := HumanPlayer{
		Scorecard: &scoreCard,
	}

	p := Player(hp)

	return &p
}
