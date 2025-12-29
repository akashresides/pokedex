package main

import (
    "fmt"
    "os"
    "strings"
)

func cleanInput(text string) []string {
    return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

type cliCommand struct {
    name        string
    description string
    callback    func() error
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
}

func commandHelp() error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    fmt.Println()
    fmt.Println("help: Displays a help message")
    fmt.Println("exit: Exit the Pokedex")
    return nil
}

