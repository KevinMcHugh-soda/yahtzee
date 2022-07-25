package yahtzee

import (
	"fmt"
)

type AIPlayer struct {
	Scorecard
}

func NewAiPlayer() *Player {
	z := 0
	scoreCard := Scorecard{YahtzeeBonusName: &z}
	ai := AIPlayer{
		Scorecard: scoreCard,
	}

	p := Player(ai)

	return &p
}

func (ai AIPlayer) GetName() string {
	return "🤖"
}

func (ai AIPlayer) GetScorecard() *Scorecard {
	return &ai.Scorecard
}

// TODO should be bonus-aware
func (ai AIPlayer) AssessRoll(hand Hand, rollsRemaining int) RollDecision {
	// calculate a targeted scorable, given incomplete scorables and probabilites of completion
	highestExpectedScore := 0.0
	var highestScorableName ScorableName
	// fmt.Println(hand)
	for _, name := range ScorableNames {
		scorable := ScoreableByName(name)
		if scorable == nil || ai.Scorecard.NameToScorePtr(name) != nil {
			continue
		}
		prob := scorable.ProbabilityToHit(hand, rollsRemaining)
		best := scorable.MaxPossible()
		// This doesn't work at all for chance, and says chance will always give you a 30 lol
		// I put in a hack to always give a 0 probability to hit chance
		expected := prob * float64(best)
		// fmt.Printf("	%s, %.2f, %d, %.2f\n", name, prob, best, expected)
		if expected >= highestExpectedScore {
			highestExpectedScore = expected
			highestScorableName = name
		}
		// TODO if expected == score then short circuit and return all keeps
	}

	strategy := StrategyForScorable(highestScorableName)
	// this seems to happen when all the unselected scorables have a 0 probability.
	// The >= on 46 should stop it.
	if strategy == nil {
		fmt.Println("picking a nil strategy for some reason", highestScorableName)
	}
	decision := strategy.PickKeepers(hand)
	// fmt.Println(hand, rollsRemaining, highestScorableName, decision)
	return decision
}

func (ai AIPlayer) PickScorable(hand Hand) Scoreable {
	highestScore := 0
	var highestScorable ScorableName
	for _, name := range ScorableNames {
		scorable := ScoreableByName(name)
		// TODO maybe a NullScorable?
		if scorable == nil || ai.Scorecard.NameToScorePtr(name) != nil {
			continue
		}
		score := scorable.Score(hand)

		// prefer harder ones, or maybe compare to best possible score
		if score >= int(highestScore) {
			highestScore = score
			highestScorable = name
		}
	}

	dec := ScoreableByName(highestScorable)
	// fmt.Printf("given %x, choosing %s\n", hand, highestScorable)
	return dec
}

func NewFaceValueStrategy(name ScorableName) ScorableVarietyStrategy {
	namesToNumbers := map[ScorableName]int{
		OnesName:   1,
		TwosName:   2,
		ThreesName: 3,
		FoursName:  4,
		FivesName:  5,
		SixesName:  6,
	}
	return FaceValueStrategy{namesToNumbers[name]}
}

func StrategyForScorable(name ScorableName) ScorableVarietyStrategy {
	variety := name.VarietyOfScorable()
	// TODO: we can probably get rid of the ScorableVariety type honestly.
	strategyMap := map[ScorableVariety]ScorableVarietyStrategy{
		FaceValueVariety: NewFaceValueStrategy(name),
		// TODO: Build a strategy for each ScorableVariety:
		// OfAKindStrategy is similar but might prefer 4,4,4,4,6 over trying for 4,4,4,4,4?
		// StraightStrategy keeps one of each (eventually handle the 1/6 thing)
		// FullHouseStrategy is a bespoke little snowflake
		OfAKindVariety:   FaceValueStrategy{},
		FullHouseVariety: FaceValueStrategy{},
		StraightVariety:  StraightStrategy{},
		ChanceVariety:    FaceValueStrategy{},
	}

	return strategyMap[variety]
}

type ScorableVarietyStrategy interface {
	PickKeepers(hand Hand) RollDecision
}

type FaceValueStrategy struct{ keptNumber int }

func (s FaceValueStrategy) PickKeepers(hand Hand) RollDecision {
	keep := make([]bool, 5)
	for idx, die := range hand {
		if die == s.keptNumber {
			keep[idx] = true
		}

	}

	return RollDecision(keep)
}

type StraightStrategy struct{}

func (s StraightStrategy) PickKeepers(hand Hand) RollDecision {
	keep := make([]bool, 5)
	// this is going to be wrong - it will keep both 1 and 6, but:
	last := 0
	for idx, cur := range hand {
		if last != cur {
			keep[idx] = true
		}
		last = cur
	}
	return RollDecision(keep)
}
