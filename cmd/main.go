package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	repoURL = "https://github.com/0xDAEF0F/whistle.git"
)

func main() {
	lastWeek := time.Now().AddDate(0, 0, -7)
	commits, cleanup, err := GetRepoCommitsFrom(lastWeek)
	if err != nil {
		fmt.Printf("Failed to get repo commits: %v\n", err)
		return
	}
	defer cleanup()

	lastCommit := commits[0]
	firstCommit := commits[len(commits)-1]

	patch, err := firstCommit.Patch(lastCommit)
	if err != nil {
		fmt.Printf("Failed to get patch: %v\n", err)
		return
	}

	// Filter out lock files from the patch output
	patchLines := strings.Split(patch.String(), "\n")
	var filteredLines []string
	isLockFile := false

	for _, line := range patchLines {
		if strings.HasPrefix(line, "diff --git") {
			isLockFile = strings.Contains(line, "lock")
		}
		if !isLockFile {
			filteredLines = append(filteredLines, line)
		}
	}

	fmt.Print(strings.Join(filteredLines, "\n"))
}

func getRepo(repoUrl string) (*git.Repository, func(), error) {
	dir, err := os.MkdirTemp("", "repo-clone")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoUrl,
	})
	if err != nil {
		return nil, cleanup, err
	}

	return repo, cleanup, nil
}

func GetRepoCommits(numCommits int) ([]*object.Commit, func(), error) {
	repo, cleanup, err := getRepo(repoURL)
	if err != nil {
		return nil, cleanup, err
	}

	commits, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, cleanup, err
	}

	commitsArray := []*object.Commit{}

	for i := 0; i < numCommits; i++ {
		commit, err := commits.Next()
		if err != nil {
			return nil, cleanup, err
		}
		commitsArray = append(commitsArray, commit)
	}

	return commitsArray, cleanup, nil
}

func GetRepoCommitsFrom(date time.Time) ([]*object.Commit, func(), error) {
	repo, cleanup, err := getRepo(repoURL)
	if err != nil {
		return nil, cleanup, err
	}

	commits, err := repo.Log(&git.LogOptions{Since: &date})
	if err != nil {
		fmt.Printf("Failed to get commit history: %v\n", err)
		return nil, cleanup, err
	}

	commitsArray := []*object.Commit{}
	commits.ForEach(func(c *object.Commit) error {
		commitsArray = append(commitsArray, c)
		return nil
	})

	return commitsArray, cleanup, nil
}
