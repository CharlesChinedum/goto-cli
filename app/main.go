package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	Directories map[string]string `json:"directories"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".goto.json")
}

func loadStore() Store {
	store := Store{Directories: make(map[string]string)}
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return store
	}
	json.Unmarshal(data, &store)
	return store
}

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
			path := strings.Join(args[4:], " ")
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

		case "edit":
			if len(args) < 5 {
				fmt.Println("Usage: gotocli goto edit <name> <newpath>")
				return
			}
			name := args[3]
			newPath := strings.Join(args[4:], " ")
			_, exists := store.Directories[name]
			if !exists {
				fmt.Fprintf(os.Stderr, "No directory found with name '%s'\n", name)
				os.Exit(1)
			}
			store.Directories[name] = newPath
			saveStore(store)
			fmt.Printf("Updated '%s' -> %s\n", name, newPath)

		case "rename":
			if len(args) < 5 {
				fmt.Println("Usage: gotocli goto rename <oldname> <newname>")
				return
			}
			oldName := args[3]
			newName := args[4]
			path, exists := store.Directories[oldName]
			if !exists {
				fmt.Fprintf(os.Stderr, "No directory found with name '%s'\n", oldName)
				os.Exit(1)
			}
			store.Directories[newName] = path
			delete(store.Directories, oldName)
			saveStore(store)
			fmt.Printf("Renamed '%s' -> '%s'\n", oldName, newName)

		default:
			fmt.Println("Unknown command:", command)
		}
	}

}
