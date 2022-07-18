package yahtzee

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
