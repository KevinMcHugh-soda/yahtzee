package balatro

import (
	"testing"
)

func TestCardValues(t *testing.T) {
	// Test face card values
	if Jack.GetValue() != 10 {
		t.Errorf("Jack should be worth 10, got %d", Jack.GetValue())
	}
	if Queen.GetValue() != 10 {
		t.Errorf("Queen should be worth 10, got %d", Queen.GetValue())
	}
	if King.GetValue() != 10 {
		t.Errorf("King should be worth 10, got %d", King.GetValue())
	}
	if Ace.GetValue() != 11 {
		t.Errorf("Ace should be worth 11, got %d", Ace.GetValue())
	}
	
	// Test number cards
	if Two.GetValue() != 2 {
		t.Errorf("Two should be worth 2, got %d", Two.GetValue())
	}
	if Ten.GetValue() != 10 {
		t.Errorf("Ten should be worth 10, got %d", Ten.GetValue())
	}
}

func TestHandEvaluation(t *testing.T) {
	// Test Pair
	pairHand := Hand{
		{Hearts, Ace},
		{Spades, Ace},
	}
	eval := EvaluateHand(pairHand)
	if eval.Type != Pair {
		t.Errorf("Expected Pair, got %s", eval.Type)
	}
	if eval.CardValue != 22 { // 11 + 11
		t.Errorf("Expected card value 22, got %d", eval.CardValue)
	}
	if eval.Multiplier != 2 {
		t.Errorf("Expected multiplier 2, got %d", eval.Multiplier)
	}
	if eval.TotalScore != 44 { // 22 * 2
		t.Errorf("Expected total score 44, got %d", eval.TotalScore)
	}

	// Test Three of a Kind
	threeKindHand := Hand{
		{Hearts, King},
		{Spades, King},
		{Clubs, King},
	}
	eval = EvaluateHand(threeKindHand)
	if eval.Type != ThreeOfAKind {
		t.Errorf("Expected Three of a Kind, got %s", eval.Type)
	}
	if eval.CardValue != 30 { // 10 + 10 + 10
		t.Errorf("Expected card value 30, got %d", eval.CardValue)
	}
	if eval.Multiplier != 4 {
		t.Errorf("Expected multiplier 4, got %d", eval.Multiplier)
	}

	// Test Flush
	flushHand := Hand{
		{Hearts, Two},
		{Hearts, Four},
		{Hearts, Six},
		{Hearts, Eight},
		{Hearts, Ten},
	}
	eval = EvaluateHand(flushHand)
	if eval.Type != Flush {
		t.Errorf("Expected Flush, got %s", eval.Type)
	}
	expectedValue := 2 + 4 + 6 + 8 + 10 // 30
	if eval.CardValue != expectedValue {
		t.Errorf("Expected card value %d, got %d", expectedValue, eval.CardValue)
	}
	if eval.Multiplier != 6 {
		t.Errorf("Expected multiplier 6, got %d", eval.Multiplier)
	}

	// Test Straight
	straightHand := Hand{
		{Hearts, Two},
		{Spades, Three},
		{Clubs, Four},
		{Diamonds, Five},
		{Hearts, Six},
	}
	eval = EvaluateHand(straightHand)
	if eval.Type != Straight {
		t.Errorf("Expected Straight, got %s", eval.Type)
	}
}

func TestDeckCreation(t *testing.T) {
	deck := NewDeck()
	if len(deck.Cards) != 52 {
		t.Errorf("Expected 52 cards in deck, got %d", len(deck.Cards))
	}

	// Count suits and ranks
	suitCounts := make(map[Suit]int)
	rankCounts := make(map[Rank]int)
	
	for _, card := range deck.Cards {
		suitCounts[card.Suit]++
		rankCounts[card.Rank]++
	}
	
	// Each suit should have 13 cards
	for suit := Hearts; suit <= Spades; suit++ {
		if suitCounts[suit] != 13 {
			t.Errorf("Expected 13 cards of suit %s, got %d", suit, suitCounts[suit])
		}
	}
	
	// Each rank should have 4 cards
	for rank := Two; rank <= Ace; rank++ {
		if rankCounts[rank] != 4 {
			t.Errorf("Expected 4 cards of rank %s, got %d", rank, rankCounts[rank])
		}
	}
}