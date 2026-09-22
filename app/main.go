package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// version is stamped at release time via -ldflags.
var version = "dev"

type Store struct {
	Directories map[string]string `json:"directories"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
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

func usage(msg string) int {
	fmt.Fprintln(os.Stderr, msg)
	return 1
}

// sortedDirectoryNames returns bookmark names in lexicographic order so
// `list` output is stable across runs (map iteration order is randomized).
func sortedDirectoryNames(dirs map[string]string) []string {
	names := make([]string, 0, len(dirs))
	for name := range dirs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func main() {
	os.Exit(run(os.Args))
}

// run executes the CLI for the given argv (including program name) and
// returns the process exit code. Usage and unknown-command diagnostics go
// to stderr with a non-zero status so wrappers that treat stdout as a path
// never cd into a usage string.
func run(args []string) int {
	// Handle version before the argument guard so that both
	// "gotocli version" and "gotocli --version" work.
	if len(args) >= 2 {
		if args[1] == "version" || args[1] == "--version" || args[1] == "-v" {
			fmt.Println(version)
			return 0
		}
	}

	if len(args) < 3 {
		return usage("Usage: gotocli goto <command> [name] [path]")
	}

	mainCommand := args[1]
	command := args[2]

	if mainCommand != "goto" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s (expected 'goto')\n", mainCommand)
		return 1
	}

	store := loadStore()

	switch command {
	case "add":
		if len(args) < 5 {
			return usage("Usage: gotocli goto add <name> <path>")
		}
		name := args[3]
		path := strings.Join(args[4:], " ")
		store.Directories[name] = path
		saveStore(store)
		fmt.Printf("Saved '%s' -> %s\n", name, path)

	case "remove":
		if len(args) < 4 {
			return usage("Usage: gotocli goto remove <name>")
		}
		name := args[3]
		if _, exists := store.Directories[name]; !exists {
			fmt.Fprintf(os.Stderr, "No directory found for '%s'\n", name)
			return 1
		}
		delete(store.Directories, name)
		saveStore(store)
		fmt.Printf("Removed '%s'\n", name)

	case "list":
		if len(store.Directories) == 0 {
			fmt.Println("No directories saved.")
			return 0
		}
		for _, name := range sortedDirectoryNames(store.Directories) {
			fmt.Printf("  %s -> %s\n", name, store.Directories[name])
		}

	case "jump":
		if len(args) < 4 {
			return usage("Usage: gotocli goto jump <name>")
		}
		name := args[3]
		path, exists := store.Directories[name]
		if !exists {
			fmt.Fprintf(os.Stderr, "No directory found for '%s'\n", name)
			return 1
		}
		fmt.Println(path)

	case "edit":
		if len(args) < 5 {
			return usage("Usage: gotocli goto edit <name> <newpath>")
		}
		name := args[3]
		newPath := strings.Join(args[4:], " ")
		_, exists := store.Directories[name]
		if !exists {
			fmt.Fprintf(os.Stderr, "No directory found with name '%s'\n", name)
			return 1
		}
		store.Directories[name] = newPath
		saveStore(store)
		fmt.Printf("Updated '%s' -> %s\n", name, newPath)

	case "rename":
		if len(args) < 5 {
			return usage("Usage: gotocli goto rename <oldname> <newname>")
		}
		oldName := args[3]
		newName := args[4]
		path, exists := store.Directories[oldName]
		if !exists {
			fmt.Fprintf(os.Stderr, "No directory found with name '%s'\n", oldName)
			return 1
		}
		store.Directories[newName] = path
		delete(store.Directories, oldName)
		saveStore(store)
		fmt.Printf("Renamed '%s' -> '%s'\n", oldName, newName)

	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", command)
		return 1
	}

	return 0
}
