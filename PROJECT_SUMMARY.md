# Pokedex Project Summary

## Project Overview
A command-line Pokedex CLI application built in Go that explores Pokemon world location areas using the PokeAPI.

## Current State
- All tests passing
- Fully functional CLI with interactive REPL
- PokeAPI integration for location areas
- Pagination support for exploring location areas
- Pokemon exploration in location areas with caching

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
- Reduces API calls for repeated location area queries
- Configurable expiration interval (currently 5 minutes)

## Key Data Structures

### config struct
```go
type config struct {
    Next     *string              // URL for next page of results
    Previous *string              // URL for previous page of results
    Cache    *pokecache.Cache      // In-memory cache for API responses
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
LATEST: feat: add explore command to show Pokemon in location areas
0b20d19 feat: implement in-memory caching system for API responses
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
- List response includes pagination via "next" and "previous" fields

## Next Steps (Potential Future Features)
- Catch Pokemon command
- Inspect Pokemon command
- Battle mechanics
- Save/Load functionality
