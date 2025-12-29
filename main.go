package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
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
        err := cmd.callback()
        if err != nil {
            fmt.Println("Error:", err)
        }
    }
}

