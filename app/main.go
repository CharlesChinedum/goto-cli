package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Store is the structure for saved directories
type Store struct {
	Directories map[string]string `json:"directories"` // name -> path
}

// getConfigPath returns ~/.goto.json
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.goto.json"
}

// loadStore reads the saved directories from disk
func loadStore() Store {
	store := Store{Directories: make(map[string]string)}
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return store // return empty store if file doesn't exist yet
	}
	json.Unmarshal(data, &store)
	return store
}

// saveStore writes the directories to disk
func saveStore(store Store) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}

func main() {
	// fmt.Println("Hello World from Go!")
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Print("Enter a line: ")
	// scanner.Scan()

	// fmt.Println("You entered:", scanner.Text())

	// save directory path and name
	// remove directory path and name
	// list directory path and name
	// jump to directory path and name

	// save directory path and name
	// var splitCommand = strings.Split(scanner.Text(), " ")
	// var mainCommand = splitCommand[0]
	// var command = splitCommand[1]
	// var directoryPath = splitCommand[2]
	// var directoryName = splitCommand[3]

	// if mainCommand == "goto" {
	// 	if command == "add" {
	// 		// save directory path and name

	// 	}

	// 	if command == "remove" {
	// 		// remove directory path and name

	// 	}

	// 	if command == "list" {
	// 		// list directory path and name

	// 	}

	// 	if command == "jump" {
	// 		// jump to directory path and name

	// 	}
	// }

	// fmt.Println(mainCommand)
	// fmt.Println(command)
	// fmt.Println(directoryPath)
	// fmt.Println(directoryName)

	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: gotocli goto <command> [name] [path]")
		return
	}

	mainCommand := args[1]
	command := args[2]

	if mainCommand == "goto" {
		store := loadStore()

		switch command {
		case "add":
			if len(args) < 5 {
				fmt.Println("Usage: gotocli goto add <name> <path>")
				return
			}
			name := args[3]
			path := strings.Join(args[4:], " ") // handles paths with spaces
			store.Directories[name] = path
			saveStore(store)
			fmt.Printf("Saved '%s' -> %s\n", name, path)

		case "remove":
			if len(args) < 4 {
				fmt.Println("Usage: gotocli goto remove <name>")
				return
			}
			name := args[3]
			delete(store.Directories, name)
			saveStore(store)
			fmt.Printf("Removed '%s'\n", name)

		case "list":
			if len(store.Directories) == 0 {
				fmt.Println("No directories saved.")
				return
			}
			for name, path := range store.Directories {
				fmt.Printf("  %s -> %s\n", name, path)
			}

		case "jump":
			if len(args) < 4 {
				fmt.Println("Usage: gotocli goto jump <name>")
				return
			}
			name := args[3]
			path, exists := store.Directories[name]
			if !exists {
				fmt.Fprintf(os.Stderr, "No directory found for '%s'\n", name)
				os.Exit(1)
			}
			fmt.Println(path)

		default:
			fmt.Println("Unknown command:", command)
		}
	}

}
