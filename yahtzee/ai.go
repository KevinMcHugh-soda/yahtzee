package yahtzee

import "fmt"

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
	var highestScorable ScorableName
	for _, name := range ScorableNames {
		scorable := ScoreableByName(name)
		if scorable == nil {
			continue
		}
		prob := scorable.ProbabilityToHit(hand, rollsRemaining)
		best := scorable.MaxPossible()
		expected := prob * float64(best)
		if expected > highestExpectedScore {
			highestExpectedScore = expected
			highestScorable = name
		}
	}
	// TODO: Build a strategy for each ScorableVariety:
	// FaceValueStrategy keeps whichever value we have most of
	// OfAKindStrategy is similar but might prefer 4,4,4,4,6 over trying for 4,4,4,4,4?
	// StraightStrategy keeps one of each (eventually handle the 1/6 thing)
	// FullHouseStrategy is a bespoke little snowflake

	strategy := StrategyForScorable(highestScorable)
	if strategy == nil {
		fmt.Println(highestScorable)
	}
	decision := strategy.PickKeepers(hand)
	fmt.Println(hand, rollsRemaining, highestScorable, decision)
	return decision
}

func (ai AIPlayer) PickScorable(hand Hand) Scoreable {
	highestScore := 0
	var highestScorable ScorableName
	fmt.Println(hand)
	for _, name := range ScorableNames {
		scorable := ScoreableByName(name)
		// TODO maybe a NullScorable?
		if scorable == nil {
			continue
		}
		score := scorable.Score(hand)

		// prefer harder ones
		if score >= int(highestScore) {
			highestScore = score
			highestScorable = name
		}
	}

	dec := ScoreableByName(highestScorable)
	fmt.Println(dec)
	return dec
}

func StrategyForScorable(name ScorableName) ScorableVarietyStrategy {
	variety := name.VarietyOfScorable()
	// TODO: we can probably get rid of the ScorableVariety type honestly.
	strategyMap := map[ScorableVariety]ScorableVarietyStrategy{
		FaceValueVariety: FaceValueStrategy{},
		// TODO:
		OfAKindVariety:   FaceValueStrategy{},
		FullHouseVariety: FaceValueStrategy{},
		StraightVariety:  FaceValueStrategy{},
		ChanceVariety:    FaceValueStrategy{},
	}

	return strategyMap[variety]
}

type ScorableVarietyStrategy interface {
	PickKeepers(hand Hand) RollDecision
}

type FaceValueStrategy struct{}

func (s FaceValueStrategy) PickKeepers(hand Hand) RollDecision {
	keep := make([]bool, 5)
	counts := valueCounts(hand)
	mostPresentValue, mostPresentCount := -1, 0
	for idx, count := range counts {
		if count > mostPresentCount {
			mostPresentValue = idx
			mostPresentCount = count
		}
	}

	for idx, die := range hand {
		if die == mostPresentValue+1 {
			keep[idx] = true
		}

	}

	return RollDecision(keep)
}
