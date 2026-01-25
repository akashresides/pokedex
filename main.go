package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/akashresides/pokedex/internal/pokecache"
)

func main() {
	cfg := &config{
		Cache: pokecache.NewCache(5 * time.Minute),
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := cleanInput(scanner.Text())

		if len(text) == 0 {
			continue
		}

		cmdName := text[0]
		cmd, exists := commands[cmdName]

		if !exists || cmd.callback == nil {
			fmt.Println("Unknown command.")
			continue
		}

		// Run the command
		err := cmd.callback(cfg, text[1:])
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
