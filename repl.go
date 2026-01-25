package main

import (
	"encoding/json"
	"fmt"
	"io"
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
