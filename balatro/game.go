package balatro

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	Deck       *Deck
	PlayerHand []Card
	Score      int
	Round      int
}

func NewGame() *Game {
	deck := NewDeck()
	deck.Shuffle()
	
	return &Game{
		Deck:       deck,
		PlayerHand: make([]Card, 0),
		Score:      0,
		Round:      1,
	}
}

func (g *Game) Play() {
	fmt.Println("=== Welcome to Balatro CLI ===")
	fmt.Println("Select up to 5 cards to form a poker hand and score points!")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("=== Round %d ===\n", g.Round)
		fmt.Printf("Current Score: %d\n", g.Score)
		fmt.Println()

		// Deal 8 cards for the player to choose from
		availableCards := g.Deck.Draw(8)
		if len(availableCards) == 0 {
			fmt.Println("No more cards in the deck! Game Over!")
			break
		}

		fmt.Println("Available cards:")
		for i, card := range availableCards {
			fmt.Printf("%d: %s ", i+1, card)
		}
		fmt.Println()
		fmt.Println()

		// Let player select up to 5 cards
		selectedCards := g.selectCards(availableCards, reader)
		
		if len(selectedCards) == 0 {
			fmt.Println("No cards selected. Ending game.")
			break
		}

		// Evaluate and score the hand
		evaluation := EvaluateHand(selectedCards)
		g.Score += evaluation.TotalScore

		fmt.Println()
		fmt.Printf("Selected hand: %s\n", Hand(selectedCards))
		fmt.Printf("Hand type: %s\n", evaluation.Type)
		fmt.Printf("Card value total: %d\n", evaluation.CardValue)
		fmt.Printf("Multiplier: %dx\n", evaluation.Multiplier)
		fmt.Printf("Round score: %d\n", evaluation.TotalScore)
		fmt.Printf("Total score: %d\n", g.Score)
		fmt.Println()

		// Ask if player wants to continue
		fmt.Print("Continue to next round? (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input != "y" && input != "yes" {
			break
		}

		g.Round++
		fmt.Println()
	}

	fmt.Printf("\nGame finished! Final score: %d\n", g.Score)
}

func (g *Game) selectCards(availableCards []Card, reader *bufio.Reader) []Card {
	selectedCards := make([]Card, 0, 5)
	
	for len(selectedCards) < 5 {
		if len(selectedCards) > 0 {
			fmt.Printf("Selected cards (%d/5): %s\n", len(selectedCards), Hand(selectedCards))
		}
		
		fmt.Printf("Select a card (1-%d) or 'done' to finish selection: ", len(availableCards))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input == "done" || input == "d" {
			break
		}
		
		cardIndex, err := strconv.Atoi(input)
		if err != nil || cardIndex < 1 || cardIndex > len(availableCards) {
			fmt.Println("Invalid selection. Please enter a number between 1 and", len(availableCards))
			continue
		}
		
		selectedCard := availableCards[cardIndex-1]
		
		// Check if card is already selected
		alreadySelected := false
		for _, card := range selectedCards {
			if card.Suit == selectedCard.Suit && card.Rank == selectedCard.Rank {
				alreadySelected = true
				break
			}
		}
		
		if alreadySelected {
			fmt.Println("Card already selected!")
			continue
		}
		
		selectedCards = append(selectedCards, selectedCard)
		fmt.Printf("Added %s to your hand\n", selectedCard)
		
		if len(selectedCards) == 5 {
			fmt.Println("Hand is full (5 cards)")
			break
		}
		fmt.Println()
	}
	
	return selectedCards
}