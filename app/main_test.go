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

func TestCorruptStore(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gotocli")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	for _, tc := range []struct {
		name    string
		content string
	}{
		{"truncated json", `{"directories": {"work": "/home/u/work", "docs": "/home/u/docs"`},
		{"bare brace", `{`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testHome := t.TempDir()
			configPath := filepath.Join(testHome, ".goto.json")
			original := []byte(tc.content)
			if err := os.WriteFile(configPath, original, 0600); err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(binary, "goto", "add", "foo", "/tmp")
			cmd.Env = append(os.Environ(), "HOME="+testHome, "USERPROFILE="+testHome)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			if cmd.ProcessState == nil {
				t.Fatalf("start CLI: %v", err)
			}
			if got := cmd.ProcessState.ExitCode(); got != 1 {
				t.Fatalf("exit code = %d, want 1; stdout=%q stderr=%q", got, &stdout, &stderr)
			}
			if stdout.Len() != 0 {
				t.Errorf("unexpected stdout: %q", &stdout)
			}
			if !bytes.Contains(stderr.Bytes(), []byte("Error loading bookmarks:")) {
				t.Errorf("stderr missing error prefix: %q", &stderr)
			}
			data, readErr := os.ReadFile(configPath)
			if readErr != nil || !bytes.Equal(data, original) {
				t.Errorf("corrupt store was overwritten: %q, error: %v", data, readErr)
			}
		})
	}
}

func TestNullDirectoriesList(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gotocli")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	testHome := t.TempDir()
	configPath := filepath.Join(testHome, ".goto.json")
	if err := os.WriteFile(configPath, []byte(`{"directories": null}`), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binary, "goto", "list")
	cmd.Env = append(os.Environ(), "HOME="+testHome, "USERPROFILE="+testHome)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if cmd.ProcessState == nil {
		t.Fatalf("start CLI: %v", err)
	}
	// null directories is valid JSON; treat as empty store, no panic
	if got := cmd.ProcessState.ExitCode(); got != 0 {
		t.Fatalf("exit code = %d, want 0; stdout=%q stderr=%q", got, &stdout, &stderr)
	}
	if stdout.String() != "No directories saved.\n" {
		t.Errorf("stdout = %q", &stdout)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q", &stderr)
	}
}

func TestNullDirectoriesAdd(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gotocli")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	testHome := t.TempDir()
	configPath := filepath.Join(testHome, ".goto.json")
	if err := os.WriteFile(configPath, []byte(`{"directories": null}`), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binary, "goto", "add", "foo", "/tmp")
	cmd.Env = append(os.Environ(), "HOME="+testHome, "USERPROFILE="+testHome)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if cmd.ProcessState == nil {
		t.Fatalf("start CLI: %v", err)
	}
	if got := cmd.ProcessState.ExitCode(); got != 0 {
		t.Fatalf("exit code = %d, want 0; stdout=%q stderr=%q", got, &stdout, &stderr)
	}
	if stdout.String() != "Saved 'foo' -> /tmp\n" || stderr.Len() != 0 {
		t.Errorf("unexpected output: stdout=%q stderr=%q", &stdout, &stderr)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		t.Fatal(err)
	}
	if store.Directories["foo"] != "/tmp" {
		t.Errorf("store = %v", store.Directories)
	}
}
