package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRemove(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gotocli")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	for _, tc := range []struct {
		name      string
		bookmark  string
		withStore bool
		wantCode  int
	}{
		{"missing bookmark", "projets", true, 1},
		{"missing store", "projects", false, 1},
		{"existing bookmark", "projects", true, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Run the CLI with a temporary home so real bookmarks are untouched.
			testHome := t.TempDir()
			configPath := filepath.Join(testHome, ".goto.json")
			original := []byte("{\"directories\":{\"projects\":\"/projects\",\"docs\":\"/docs\"}}\n")
			if tc.withStore {
				if err := os.WriteFile(configPath, original, 0600); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command(binary, "goto", "remove", tc.bookmark)
			cmd.Env = append(os.Environ(), "HOME="+testHome, "USERPROFILE="+testHome)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			if cmd.ProcessState == nil {
				t.Fatalf("start CLI: %v", err)
			}
			if got := cmd.ProcessState.ExitCode(); got != tc.wantCode {
				t.Errorf("exit code = %d, want %d; stdout=%q stderr=%q", got, tc.wantCode, &stdout, &stderr)
			}
			data, readErr := os.ReadFile(configPath)
			if tc.wantCode != 0 {
				if stdout.Len() != 0 {
					t.Errorf("unexpected success output: %q", &stdout)
				}
				if want := "No directory found for '" + tc.bookmark + "'\n"; stderr.String() != want {
					t.Errorf("stderr = %q, want %q", &stderr, want)
				}
				if tc.withStore {
					if readErr != nil || !bytes.Equal(data, original) {
						t.Errorf("failed removal changed the store: %q, error: %v", data, readErr)
					}
				} else if !os.IsNotExist(readErr) {
					t.Errorf("failed removal created a store; read error: %v", readErr)
				}
				return
			}
			if stdout.String() != "Removed 'projects'\n" || stderr.Len() != 0 {
				t.Errorf("unexpected output: stdout=%q stderr=%q", &stdout, &stderr)
			}
			var store Store
			if readErr != nil {
				t.Fatal(readErr)
			}
			if err := json.Unmarshal(data, &store); err != nil {
				t.Fatal(err)
			}
			if len(store.Directories) != 1 || store.Directories["docs"] != "/docs" {
				t.Errorf("unexpected remaining bookmarks: %v", store.Directories)
			}
		})
	}
}

func TestRename(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gotocli")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	for _, tc := range []struct {
		name     string
		oldName  string
		newName  string
		wantCode int
	}{
		{"onto existing name", "a", "b", 1},
		{"onto free name", "a", "c", 0},
		{"same name no-op", "a", "a", 0},
		{"missing old name", "missing", "c", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testHome := t.TempDir()
			configPath := filepath.Join(testHome, ".goto.json")
			original := []byte("{\"directories\":{\"a\":\"/a\",\"b\":\"/b\"}}\n")
			if err := os.WriteFile(configPath, original, 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(binary, "goto", "rename", tc.oldName, tc.newName)
			cmd.Env = append(os.Environ(), "HOME="+testHome, "USERPROFILE="+testHome)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			if cmd.ProcessState == nil {
				t.Fatalf("start CLI: %v", err)
			}
			if got := cmd.ProcessState.ExitCode(); got != tc.wantCode {
				t.Errorf("exit code = %d, want %d; stdout=%q stderr=%q", got, tc.wantCode, &stdout, &stderr)
			}
			data, readErr := os.ReadFile(configPath)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if tc.wantCode != 0 {
				if stdout.Len() != 0 {
					t.Errorf("unexpected success output: %q", &stdout)
				}
				if !bytes.Equal(data, original) {
					t.Errorf("failed rename changed the store: %q", data)
				}
				if tc.oldName == "missing" {
					if want := "No directory found with name 'missing'\n"; stderr.String() != want {
						t.Errorf("stderr = %q, want %q", &stderr, want)
					}
				} else if want := "'" + tc.newName + "' already exists; remove it first\n"; stderr.String() != want {
					t.Errorf("stderr = %q, want %q", &stderr, want)
				}
				return
			}
			if stderr.Len() != 0 {
				t.Errorf("unexpected stderr: %q", &stderr)
			}
			var store Store
			if err := json.Unmarshal(data, &store); err != nil {
				t.Fatal(err)
			}
			if tc.oldName == tc.newName {
				if len(store.Directories) != 2 || store.Directories["a"] != "/a" || store.Directories["b"] != "/b" {
					t.Errorf("same-name rename mutated store: %v", store.Directories)
				}
				return
			}
			if _, stillOld := store.Directories[tc.oldName]; stillOld {
				t.Errorf("old name still present: %v", store.Directories)
			}
			if store.Directories[tc.newName] != "/a" || store.Directories["b"] != "/b" {
				t.Errorf("unexpected store after rename: %v", store.Directories)
			}
		})
	}
}
