package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetRepoCommits(t *testing.T) {
	commits, cleanup, err := GetRepoCommits(10)
	if err != nil {
		t.Errorf("Failed to get repo commits: %v\n", err)
	}
	defer cleanup()

	if len(commits) != 10 {
		t.Errorf("Expected 10 commits, got %d\n", len(commits))
	}
}

func TestGetRepoCommitsFrom(t *testing.T) {
	lastWeek := time.Now().AddDate(0, 0, -7)
	commits, cleanup, err := GetRepoCommitsFrom(lastWeek)
	if err != nil {
		t.Errorf("Failed to get repo commits: %v\n", err)
	}
	defer cleanup()

	lastCommit := commits[0]
	firstCommit := commits[len(commits)-1]

	patch, err := firstCommit.Patch(lastCommit)
	if err != nil {
		t.Errorf("Failed to get patch: %v\n", err)
	}

	stats := patch.Stats()

	fmt.Printf("Stats: %v\n", stats)
}
