package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

func buildCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "gotocli")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	return binary
}

func TestSaveStoreAtomicRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	store := Store{Directories: map[string]string{
		"projects": "/projects",
		"docs":     "/docs",
	}}
	if err := saveStore(store); err != nil {
		t.Fatal(err)
	}

	p := getConfigPath()
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	var got Store
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("saved file is not valid JSON: %v\n%s", err, data)
	}
	if got.Directories["projects"] != "/projects" || got.Directories["docs"] != "/docs" {
		t.Errorf("directories = %v", got.Directories)
	}

	if _, err := os.Stat(p + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("leftover temp file: %v", err)
	}

	loaded := loadStore()
	if loaded.Directories["projects"] != "/projects" || loaded.Directories["docs"] != "/docs" {
		t.Errorf("loadStore = %v", loaded.Directories)
	}
}

func TestAddSaveFailureDoesNotClaimSuccess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod cannot make the config dir unwritable on Windows")
	}

	binary := buildCLI(t)
	testHome := t.TempDir()
	configPath := filepath.Join(testHome, ".goto.json")
	original := []byte("{\"directories\":{\"work\":\"/w\"}}\n")
	if err := os.WriteFile(configPath, original, 0444); err != nil {
		t.Fatal(err)
	}
	// A 0444 file can still be replaced by rename if the directory is writable,
	// so lock down the directory so the temp file cannot be created.
	if err := os.Chmod(testHome, 0555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(testHome, 0700) })

	cmd := exec.Command(binary, "goto", "add", "foo", "/tmp")
	cmd.Env = append(os.Environ(), "HOME="+testHome, "USERPROFILE="+testHome)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if cmd.ProcessState == nil {
		t.Fatalf("start CLI: %v", err)
	}
	if got := cmd.ProcessState.ExitCode(); got != 1 {
		t.Errorf("exit code = %d, want 1; stdout=%q stderr=%q", got, &stdout, &stderr)
	}
	if stdout.Len() != 0 {
		t.Errorf("claimed success on write failure: stdout=%q", &stdout)
	}
	if !strings.Contains(stderr.String(), "Error saving store:") {
		t.Errorf("stderr = %q, want an error saving store", &stderr)
	}

	_ = os.Chmod(testHome, 0700)
	data, readErr := os.ReadFile(configPath)
	if readErr != nil || !bytes.Equal(data, original) {
		t.Errorf("failed save changed the store: %q, error: %v", data, readErr)
	}
}
