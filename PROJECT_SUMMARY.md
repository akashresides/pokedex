# Pokedex Project Summary

## Project Overview
A command-line Pokedex CLI application built in Go that explores Pokemon world location areas using the PokeAPI.

## Current State
 - All tests passing
 - Fully functional CLI with interactive REPL
 - PokeAPI integration for location areas
 - Pagination support for exploring location areas
 - Pokemon exploration in location areas with caching
 - Pokemon catching mechanics with base experience-based probability
 - Pokemon inspection to view caught Pokemon details
 - Pokedex listing to view all caught Pokemon

## Project Structure
```
pokedex/
├── main.go                    # Entry point with REPL loop
├── repl.go                    # Command definitions and handlers
├── repl_test.go               # Test suite
├── internal/
│   ├── pokeapi/
│   │   ├── client.go          # HTTP client for PokeAPI
│   │   └── types.go           # Response types
│   └── pokecache/
│       ├── cache.go           # Caching implementation
│       └── cache_test.go      # Cache tests
├── go.mod                     # Go module file
└── README.md                  # This summary
```

## Implemented Features

### 1. REPL Core (main.go)
- Interactive command-line interface with "Pokedex > " prompt
- Input cleaning (lowercase, trim whitespace, split by whitespace)
- Command lookup and execution
- Config struct for maintaining state across commands
- Cache initialization on startup (5-minute expiration interval)

### 2. Available Commands (repl.go)
- **help**: Displays available commands and descriptions
- **exit**: Exits the Pokedex gracefully
- **map**: Displays next 20 location areas from PokeAPI
- **mapb**: Displays previous 20 location areas (with "you're on the first page" message)
 - **explore <location_area>**: Displays all Pokemon found in a location area
 - **catch <pokemon>**: Attempts to catch a Pokemon using base experience to determine catch chance
 - **inspect <pokemon>**: View details about a caught Pokemon (name, height, weight, stats, types)
 - **pokedex**: List all caught Pokemon

### 3. PokeAPI Integration (internal/pokeapi/)
- HTTP client with 1-minute timeout
- Location area endpoint fetching
- JSON unmarshaling into Go structs
- Pagination support via next/previous URLs

### 4. Caching System (internal/pokecache/)
- Thread-safe in-memory cache using map[string]CacheEntry
- Expiration-based entry removal via reapLoop goroutine
- Cache entries stored with timestamp for expiration tracking
- Helper function fetchLocationArea() checks cache before HTTP requests
- Helper function fetchPokemon() checks cache before HTTP requests
- Reduces API calls for repeated location area and Pokemon queries
- Configurable expiration interval (currently 5 minutes)

### 5. Catch Command (repl.go)
- Fetches Pokemon data from PokeAPI using Pokemon endpoint
- Uses base experience to calculate catch probability (higher XP = harder to catch)
- Catch chance formula: 1.0 - min(base_experience / 100.0, 0.9)
- Random determination using math/rand.Float64()
- Successfully caught Pokemon stored in Pokedex map
- Prints "Throwing a Pokeball at <name>..." before attempting catch

## Key Data Structures

### config struct
```go
type config struct {
    Next     *string              // URL for next page of results
    Previous *string              // URL for previous page of results
    Cache    *pokecache.Cache      // In-memory cache for API responses
    Pokedex  map[string]Pokemon   // User's caught Pokemon collection
}
```

### LocationArea struct
```go
type LocationArea struct {
    Name string
    URL  string
}
```

### LocationAreasResponse struct
```go
type LocationAreasResponse struct {
    Count    int
    Next     *string
    Previous *string
    Results  []LocationArea
}
```

### PokemonEncounter struct
```go
type PokemonEncounter struct {
    Pokemon struct {
        Name string
        URL  string
    }
}
```

### LocationAreaDetail struct
```go
type LocationAreaDetail struct {
    Name              string
    PokemonEncounters []PokemonEncounter
}
```

### CacheEntry struct
```go
type CacheEntry struct {
    CreatedAt time.Time   // Timestamp when entry was added
    Val       []byte      // Cached response data
}
```

### Cache struct
```go
type Cache struct {
    entries  map[string]CacheEntry  // Thread-safe map for cached data
    mu       sync.RWMutex            // Read/write mutex for thread safety
    interval time.Duration          // Expiration interval for entries
}
```

### Pokemon struct
```go
type Pokemon struct {
    Name           string
    BaseExperience int
    Height         int
    Weight         int
    Stats          []Stat
    Types          []PokemonType
}

type Stat struct {
    BaseStat int
    Stat     struct {
        Name string
    }
}

type PokemonType struct {
    Type struct {
        Name string
    }
}
```

## Testing

### REPL Tests (repl_test.go)
- TestCleanInput: Validates input cleaning functionality
- TestCommandsMapContainsExpectedCommands: Ensures commands are registered
- TestCommandStructureValidation: Validates command structure
- TestInvalidCommandLookup: Tests invalid command handling
- TestCommandNamesAreUnique: Checks for duplicate commands
- TestCommandDescriptionsArePopulated: Validates descriptions

### Cache Tests (internal/pokecache/cache_test.go)
- TestAddGet: Tests basic add and get functionality
- TestReapLoop: Tests that expired entries are removed
- TestReapLoopNotExpired: Tests that non-expired entries are kept
- TestAddOverwrite: Tests that adding the same key overwrites the value
- TestGetNonExistent: Tests that getting a non-existent key returns false

Run tests with: `go test -v` or `go test ./...`

## Git History
```
LATEST: feat: add catch command to catch Pokemon with base experience-based probability
0b20d19 feat: add explore command to show Pokemon in location areas
0b20d18 feat: implement in-memory caching system for API responses
571e3b6 test: add command registry validation tests
7bdc340 test: improve CleanInput test coverage with edge cases
6698775 chore: stop tracking and ignore pokedex binary
336e684 Bootdev pokedex project part 1
```

## Running the Project
```bash
# Build and run
go run .

# Or build binary
go build
./pokedex
```

## Example Usage
```
Pokedex > help
Welcome to the Pokedex!
Usage:

  help: Displays a help message
  exit: Exit the Pokedex
  map: Displays the next 20 location areas in the Pokemon world
  mapb: Displays the previous 20 location areas in the Pokemon world
  explore <location_area>: Displays a list of all Pokemon in a location area
  catch <pokemon>: Attempt to catch a Pokemon
  inspect <pokemon>: View details about a caught Pokemon
  pokedex: List all caught Pokemon

Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
... (17 more)

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
 - gyarados
 - remoraid
 - octillery
 - wingull
 - pelipper
 - shellos
 - gastrodon

Pokedex > map
mt-coronet-1f-route-216
mt-coronet-1f-route-211
... (18 more)

Pokedex > mapb
canalave-city-area
eterna-city-area
pastoria-city-area
... (17 more)

Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu escaped!

 Pokedex > catch pikachu
 Throwing a Pokeball at pikachu...
 pikachu was caught!

 Pokedex > inspect pikachu
 Name: pikachu
 Height: 4
 Weight: 60
 Stats:
  -hp: 35
  -attack: 55
  -defense: 40
  -special-attack: 50
  -special-defense: 50
  -speed: 90
 Types:
  - electric

  Pokedex > inspect pidgey
  you have not caught that pokemon

  Pokedex > catch pidgey
  Throwing a Pokeball at pidgey...
  pidgey was caught!
  You may now inspect it with the inspect command.

  Pokedex > catch caterpie
  Throwing a Pokeball at caterpie...
  caterpie was caught!
  You may now inspect it with the inspect command.

  Pokedex > pokedex
  Your Pokedex:
   - pidgey
   - caterpie

  Pokedex > exit
  Closing the Pokedex... Goodbye!
  ```

## Go Module
```
module github.com/akashresides/pokedex
go 1.22.2
```

## API Reference
- PokeAPI Base URL: https://pokeapi.co/api/v2
- Location Area Endpoint: /location-area/ (list)
- Location Area Detail Endpoint: /location-area/{name}/ (with Pokemon encounters)
- Pokemon Endpoint: /pokemon/{name}/ (Pokemon details including base_experience)
- List response includes pagination via "next" and "previous" fields

 ## Next Steps (Potential Future Features)
 - Battle mechanics
 - Save/Load functionality
 