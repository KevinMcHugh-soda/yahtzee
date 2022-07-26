package yahtzee

import "math"

// This lies, kind of. It returns the greater probability to hit of the small or large straight
func (ls LargeStraight) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	coveredVals := valueCounts(hand)
	missingCountLow := 0
	for i := 0; i <= 4; i++ {
		if coveredVals[i] == 0 {
			missingCountLow += 1
		}
	}

	missingCountHigh := 0
	for i := 1; i <= 5; i++ {
		if coveredVals[i] == 0 {
			missingCountHigh += 1
		}
	}
	lowProb := (6 - float64(missingCountLow)) / (6.0 * float64(rollsRemaining))
	highProb := (6 - float64(missingCountHigh)) / (6.0 * float64(rollsRemaining))
	if lowProb > highProb {
		return lowProb
	}
	return highProb
}

func (ss SmallStraight) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	return 0.0
}

func (s Ones) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	missing := 1 - counts[1]
	return math.Pow(float64(rollsRemaining)/6.0, float64(missing))
}

func (s Twos) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	missing := 1 - counts[2]
	return math.Pow(float64(rollsRemaining)/6.0, float64(missing))
}

func (s Threes) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	missing := 1 - counts[3]
	return math.Pow(float64(rollsRemaining)/6.0, float64(missing))
}

func (s Fours) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	missing := 1 - counts[4]
	return math.Pow(float64(rollsRemaining)/6.0, float64(missing))
}

func (s Fives) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	missing := 1 - counts[5]
	return math.Pow(float64(rollsRemaining)/6.0, float64(missing))
}

func (s Sixes) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	missing := 1 - counts[6]
	return math.Pow(float64(rollsRemaining)/6.0, float64(missing))
}

func (s ThreeOfAKind) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	highestCount := 0
	for _, count := range counts {
		if count > highestCount {
			highestCount = count
		}
	}
	if highestCount >= 3 {
		return 1.0
	}
	// TODO
	return 0.0
}

func (s FourOfAKind) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	highestCount := 0
	for _, count := range counts {
		if count > highestCount {
			highestCount = count
		}
	}
	if highestCount >= 4 {
		return 1.0
	}
	// TODO
	return 0.0
}

func (s FullHouse) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	return 0.0
}

func (s Chance) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	return 0.0
}

func (s Yahtzee) ProbabilityToHit(hand Hand, rollsRemaining int) float64 {
	counts := valueCounts(hand)
	highestCount := 0
	for _, count := range counts {
		if count > highestCount {
			highestCount = count
		}
	}
	if highestCount == 5 {
		return 1.0
	}
	// TODO
	return 0.0
}
