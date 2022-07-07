package yahtzee

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Player interface {
	GetName() string
	GetScorecard() *Scorecard
	AssessRoll(hand Hand) RollDecision
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
func (p HumanPlayer) AssessRoll(hand Hand) RollDecision {
	fmt.Printf("Roll: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	fmt.Println("Type y to keep, space to reroll:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	bools := make([]bool, 5)
	for idx, c := range text {
		if idx < len(bools) {
			bools[idx] = c == rune('y')
		}
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
			score := ScoreableByName(name).Score(hand)
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
