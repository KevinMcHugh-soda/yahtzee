package main

import (
	"fmt"
	"os"

	"kevinmchugh.me/yahtzee/m/v2/balatro"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "yahtzee" {
		fmt.Println("Yahtzee mode is deprecated in this Balatro build")
		os.Exit(1)
	}
	
	game := balatro.NewGame()
	game.Play()
}
