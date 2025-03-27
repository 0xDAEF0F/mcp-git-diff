package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetRepoCommits(t *testing.T) {
	commits, err := GetRepoCommits(1)
	if err != nil {
		t.Errorf("Failed to get repo commits: %v\n", err)
	}

	commit, err := commits.Next()

	fmt.Printf("Author email: %s\n", commit.Author.Email)

	if len(commit.Hash.String()) != 40 {
		t.Errorf("Commit hash is not 40 characters: %s\n", commit.Hash.String())
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
