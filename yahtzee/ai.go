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
	bestProportion := 0.0
	var bestScorableName ScorableName
	// fmt.Println(hand)
	for _, name := range ScorableNames {
		scorable := ScoreableByName(name)
		if scorable == nil || ai.Scorecard.NameToScorePtr(name) != nil || name == ChanceName {
			continue
		}
		prob := scorable.ProbabilityToHit(hand, rollsRemaining)
		max := scorable.MaxPossible()

		expected := prob * float64(max)
		proportion := expected / float64(max)
		if name == LargeStraightName && prob < 1.0 {
			// arbitrary but decent
			proportion -= 0.5
		}
		if name.VarietyOfScorable() == FaceValueVariety {
			if prob >= 1.0 { // I cheated probability and called it out of 3, to prioritize bonus
				proportion += 0.25
			}
		}
		fmt.Printf("	%s, %.2f, %d, %.2f\n", name, prob, max, proportion)
		if proportion >= bestProportion {
			bestProportion = proportion
			bestScorableName = name
		}
		// TODO if expected == score then short circuit and return all keeps
	}

	if bestScorableName == "" {
		bestScorableName = ChanceName
	}

	strategy := StrategyForScorable(bestScorableName)
	// this seems to happen when all the unselected scorables have a 0 probability.
	// The >= on 46 should stop it.
	if strategy == nil {
		fmt.Println("picking a nil strategy for some reason", bestScorableName)
	}
	decision := strategy.PickKeepers(hand)
	fmt.Printf("roll: %v; hand: %d; chasing: %s; holding: %v\n", hand, rollsRemaining, bestScorableName, decision)
	return decision
}

func (ai AIPlayer) PickScorable(hand Hand) Scoreable {
	highestScore := 0
	var highestScorable ScorableName
	for _, name := range ScorableNames {
		scorable := ScoreableByName(name)
		// TODO maybe a NullScorable?
		if scorable == nil || ai.Scorecard.NameToScorePtr(name) != nil || name == ChanceName {
			continue
		}
		score := scorable.Score(hand, ai.Scorecard.HadYahztee())
		if name.VarietyOfScorable() == FaceValueVariety {
			if scorable.ProbabilityToHit(hand, 0) > 1.0 { // I cheated probability and called it out of 3, to prioritize bonus
				score += 10
			}
		}
		// prefer harder ones, or maybe compare to best possible score
		if score >= int(highestScore) {
			highestScore = score
			highestScorable = name
		}
	}

	// We didn't find anything worth scoring, throw it in Chance.
	if highestScorable == "" {
		highestScorable = ChanceName
	}

	dec := ScoreableByName(highestScorable)
	fmt.Printf("given %x, choosing %s\n", hand, highestScorable)
	// fmt.Println(hand, "-", highestScorable)

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
		// StraightStrategy keeps one of each (eventually handle the 1/6 thing)
		// FullHouseStrategy is a bespoke little snowflake
		OfAKindVariety:   OfAKindValueStrategy{},
		FullHouseVariety: FaceValueStrategy{},
		StraightVariety:  StraightStrategy{},
		ChanceVariety:    OfAKindValueStrategy{},
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

type OfAKindValueStrategy struct{}

func (s OfAKindValueStrategy) PickKeepers(hand Hand) RollDecision {
	keep := make([]bool, 5)
	counts := valueCounts(hand)
	mostPresentValue, mostPresentCount := 1, 0
	for idx, count := range counts {
		if count >= mostPresentCount {
			mostPresentValue = idx
			mostPresentCount = count
		}
	}

	for idx, die := range hand {
		if die == mostPresentValue {
			keep[idx] = true
		}

	}

	return RollDecision(keep)
}

type StraightStrategy struct{}

func (s StraightStrategy) PickKeepers(hand Hand) RollDecision {
	keep := make([]bool, 5)
	// this is going to be wrong - it will keep both 1 and 6, but:
	// also wrong for small straight
	last := 0
	for idx, cur := range hand {
		if last != cur {
			keep[idx] = true
		}
		last = cur
	}
	return RollDecision(keep)
}
