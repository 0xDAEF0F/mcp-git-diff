package gitops

import (
	"fmt"
	"testing"
	"time"
)

const (
	repoURL = "https://github.com/0xDAEF0F/whistle.git"
)

func TestGetRepo(t *testing.T) {
	_, cleanup, err := getGitRepo(repoURL)
	if err != nil {
		t.Errorf("Failed to get repo: %v\n", err)
	}
	defer cleanup()

	fmt.Printf("Repo url: %v\n", repoURL)
}

func TestGetDiffWithOpts(t *testing.T) {
	_, err := GetDiffWithOpts(repoURL, &CommitOpts{
		Since: time.Now().AddDate(-1, 0, 0),
	})
	if err != nil {
		t.Errorf("Failed to get diff: %v\n", err)
	}
}

func TestGetContributors(t *testing.T) {
	contributors, err := GetContributors("https://github.com/0xDAEF0F/physical-nft.git")
	if err != nil {
		t.Fatalf("Failed to get contributors: %v", err)
	}

	fmt.Printf("Contributors: %v\n", contributors)

	if len(contributors) != 3 {
		t.Fatalf("Expected 3 contributors, got %d", len(contributors))
	}

	expectedContributors := []Contributor{
		{
			Name:  "daemon",
			Email: "0xDAEF0F@proton.me",
		},
		{
			Name:  "dmojarrot",
			Email: "a01638460@itesm.mx",
		},
		{
			Name:  "dmojarrot",
			Email: "62604183+dmojarrot@users.noreply.github.com",
		},
	}

	found := false
	for _, c := range contributors {
		for _, expected := range expectedContributors {
			if c.Name == expected.Name && c.Email == expected.Email {
				found = true
				break
			}
		}
	}

	if !found {
		t.Errorf("Expected contributor %v not found in %v", expectedContributors, contributors)
	}
}
