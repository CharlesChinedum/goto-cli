package main

import (
	"reflect"
	"testing"
)

func TestSortedDirectoryNamesStableOrder(t *testing.T) {
	dirs := map[string]string{
		"gamma":   "/g",
		"alpha":   "/a",
		"epsilon": "/e",
		"beta":    "/b",
		"delta":   "/d",
	}
	want := []string{"alpha", "beta", "delta", "epsilon", "gamma"}

	for i := 0; i < 50; i++ {
		got := sortedDirectoryNames(dirs)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("iteration %d: got %v, want %v", i, got, want)
		}
	}
}

func TestSortedDirectoryNamesEmpty(t *testing.T) {
	got := sortedDirectoryNames(map[string]string{})
	if len(got) != 0 {
		t.Fatalf("got %v, want empty slice", got)
	}
}
