package test

import (
	"os"
	"slices"
	"testing"
)

func CreateEmptyDir(t testing.TB, path string) {
	t.Helper()

	err := os.MkdirAll(path, 0o755)
	RequireNoError(t, err)
}

func CreateEmptyFile(t testing.TB, path string) {
	t.Helper()

	err := os.WriteFile(path, nil, 0o644)
	RequireNoError(t, err)
}

func AssertStringSlicesEqual(t testing.TB, got, want []string) {
	t.Helper()

	slices.Sort(got)
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

func RequireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error; got %v", err)
	}
}
