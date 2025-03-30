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
		Since: time.Now().AddDate(0, 0, -3),
	})
	if err != nil {
		t.Errorf("Failed to get diff: %v\n", err)
	}
}
