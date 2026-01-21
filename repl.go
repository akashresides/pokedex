package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akashresides/pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

type config struct {
	Next     *string
	Previous *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var commands = map[string]cliCommand{
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "Displays the next 20 location areas in the Pokemon world",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the previous 20 location areas in the Pokemon world",
		callback:    commandMapBack,
	},
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the next 20 location areas in the Pokemon world")
	fmt.Println("mapb: Displays the previous 20 location areas in the Pokemon world")
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	client := pokeapi.NewClient()
	locationAreas, err := client.GetLocationAreas(cfg.Next)
	if err != nil {
		return err
	}

	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	client := pokeapi.NewClient()
	locationAreas, err := client.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}
