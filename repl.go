package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/akashresides/pokedex/internal/pokeapi"
	"github.com/akashresides/pokedex/internal/pokecache"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
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
	"explore": {
		name:        "explore",
		description: "Displays a list of all Pokemon in a location area",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Attempt to catch a Pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "View details about a caught Pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "List all caught Pokemon",
		callback:    commandPokedex,
	},
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the next 20 location areas in the Pokemon world")
	fmt.Println("mapb: Displays the previous 20 location areas in the Pokemon world")
	fmt.Println("explore <location_area>: Displays a list of all Pokemon in a location area")
	fmt.Println("catch <pokemon>: Attempt to catch a Pokemon")
	fmt.Println("inspect <pokemon>: View details about a caught Pokemon")
	fmt.Println("pokedex: List all caught Pokemon")
	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	locationAreas, err := fetchLocationArea(url, cfg.Cache)
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

func fetchLocationArea(url string, cache *pokecache.Cache) (pokeapi.LocationAreasResponse, error) {
	data, found := cache.Get(url)
	if found {
		var locationAreas pokeapi.LocationAreasResponse
		err := json.Unmarshal(data, &locationAreas)
		if err != nil {
			return pokeapi.LocationAreasResponse{}, err
		}
		return locationAreas, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pokeapi.LocationAreasResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return pokeapi.LocationAreasResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return pokeapi.LocationAreasResponse{}, fmt.Errorf("bad status code: %v", res.StatusCode)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return pokeapi.LocationAreasResponse{}, err
	}

	cache.Add(url, data)

	var locationAreas pokeapi.LocationAreasResponse
	err = json.Unmarshal(data, &locationAreas)
	if err != nil {
		return pokeapi.LocationAreasResponse{}, err
	}

	return locationAreas, nil
}

func commandMapBack(cfg *config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := fetchLocationArea(*cfg.Previous, cfg.Cache)
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

func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a location area name")
	}

	locationName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationName)

	locationArea, err := fetchLocationAreaDetail(url, cfg.Cache)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", locationArea.Name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a Pokemon name")
	}

	pokemonName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := fetchPokemon(url, cfg.Cache)
	if err != nil {
		return err
	}

	baseExperience := pokemon.BaseExperience
	catchChance := float64(baseExperience) / 100.0
	if catchChance > 0.9 {
		catchChance = 0.9
	}
	catchChance = 1.0 - catchChance

	if rand.Float64() < catchChance {
		cfg.Pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func fetchLocationAreaDetail(url string, cache *pokecache.Cache) (pokeapi.LocationAreaDetail, error) {
	data, found := cache.Get(url)
	if found {
		var locationAreaDetail pokeapi.LocationAreaDetail
		err := json.Unmarshal(data, &locationAreaDetail)
		if err != nil {
			return pokeapi.LocationAreaDetail{}, err
		}
		return locationAreaDetail, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pokeapi.LocationAreaDetail{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return pokeapi.LocationAreaDetail{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return pokeapi.LocationAreaDetail{}, fmt.Errorf("bad status code: %v", res.StatusCode)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return pokeapi.LocationAreaDetail{}, err
	}

	cache.Add(url, data)

	var locationAreaDetail pokeapi.LocationAreaDetail
	err = json.Unmarshal(data, &locationAreaDetail)
	if err != nil {
		return pokeapi.LocationAreaDetail{}, err
	}

	return locationAreaDetail, nil
}

func fetchPokemon(url string, cache *pokecache.Cache) (pokeapi.Pokemon, error) {
	data, found := cache.Get(url)
	if found {
		var pokemon pokeapi.Pokemon
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return pokeapi.Pokemon{}, err
		}
		return pokemon, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pokeapi.Pokemon{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return pokeapi.Pokemon{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return pokeapi.Pokemon{}, fmt.Errorf("bad status code: %v", res.StatusCode)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return pokeapi.Pokemon{}, err
	}

	cache.Add(url, data)

	var pokemon pokeapi.Pokemon
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return pokeapi.Pokemon{}, err
	}

	return pokemon, nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a Pokemon name")
	}

	pokemonName := args[0]
	pokemon, exists := cfg.Pokedex[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
