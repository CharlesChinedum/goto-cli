package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Store is the structure for saved directories
type Store struct {
	Directories map[string]string `json:"directories"` // name -> path
}

// getConfigPath returns ~/.goto.json
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".goto.json") // filepath.Join handles / vs \ automatically
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
