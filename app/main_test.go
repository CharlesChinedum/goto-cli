package main

import (
	"bytes"
	"encoding/json"
	"io"
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

func captureOutput(t *testing.T, fn func() int) (stdout, stderr string, code int) {
	t.Helper()
	oldOut, oldErr := os.Stdout, os.Stderr
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = wOut, wErr

	doneOut := make(chan string)
	doneErr := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		doneOut <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rErr)
		doneErr <- buf.String()
	}()

	code = fn()

	wOut.Close()
	wErr.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	stdout = <-doneOut
	stderr = <-doneErr
	return stdout, stderr, code
}

func TestRun_UsageAndUnknownGoToStderr(t *testing.T) {
	// Not parallel: captureOutput swaps process-wide os.Stdout/os.Stderr.
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStderr string
		wantStdout string
	}{
		{
			name:       "missing args",
			args:       []string{"gotocli"},
			wantCode:   1,
			wantStderr: "Usage: gotocli goto <command> [name] [path]",
		},
		{
			name:       "jump missing name",
			args:       []string{"gotocli", "goto", "jump"},
			wantCode:   1,
			wantStderr: "Usage: gotocli goto jump <name>",
		},
		{
			name:       "unknown subcommand",
			args:       []string{"gotocli", "goto", "bogus"},
			wantCode:   1,
			wantStderr: "Unknown command: bogus",
		},
		{
			name:       "non-goto first arg with enough args",
			args:       []string{"gotocli", "notgoto", "list"},
			wantCode:   1,
			wantStderr: "Unknown command: notgoto (expected 'goto')",
		},
		{
			name:       "short non-goto still usage not silent",
			args:       []string{"gotocli", "list"},
			wantCode:   1,
			wantStderr: "Usage: gotocli goto <command> [name] [path]",
		},
		{
			name:       "version still on stdout",
			args:       []string{"gotocli", "version"},
			wantCode:   0,
			wantStdout: "dev",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := captureOutput(t, func() int { return run(tc.args) })
			if code != tc.wantCode {
				t.Fatalf("exit = %d, want %d\nstdout=%q\nstderr=%q", code, tc.wantCode, stdout, stderr)
			}
			if tc.wantStderr != "" && !strings.Contains(stderr, tc.wantStderr) {
				t.Fatalf("stderr = %q, want substring %q", stderr, tc.wantStderr)
			}
			if tc.wantStderr != "" && strings.Contains(stdout, strings.TrimSpace(tc.wantStderr)) {
				t.Fatalf("diagnostic leaked to stdout: %q", stdout)
			}
			if tc.wantStdout != "" && !strings.Contains(stdout, tc.wantStdout) {
				t.Fatalf("stdout = %q, want substring %q", stdout, tc.wantStdout)
			}
			if code != 0 && strings.TrimSpace(stdout) != "" && tc.wantStdout == "" {
				// Usage/unknown must leave stdout empty for jump wrappers.
				if strings.Contains(stdout, "Usage:") || strings.Contains(stdout, "Unknown command") {
					t.Fatalf("diagnostic on stdout: %q", stdout)
				}
			}
		})
	}
}
