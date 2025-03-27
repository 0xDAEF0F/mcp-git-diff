package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	repoURL = "https://github.com/0xDAEF0F/whistle.git"
)

func main() {
	numCommits := 1
	commits, err := GetRepoCommits(numCommits)
	if err != nil {
		fmt.Printf("Failed to get repo commits: %v\n", err)
		return
	}

	count := 0
	commits.ForEach(func(c *object.Commit) error {
		firstFiveChars := c.Hash.String()[:5]
		if count < numCommits {
			fmt.Printf("Commit %d: %s - %s\n", count+1, firstFiveChars, c.Message)

			// Get parent to compare with
			if parent, err := c.Parent(0); err == nil {
				// Get changes between commit and its parent
				patch, err := parent.Patch(c)
				if err == nil {
					stats := patch.Stats()
					fmt.Printf("  Files changed: %d\n", len(stats))
					for _, stat := range stats {
						fmt.Printf("  %s: +%d -%d\n", stat.Name, stat.Addition, stat.Deletion)
					}
					fmt.Println(patch.String()) // Print the full diff
				}
			} else if err == object.ErrParentNotFound {
				// Initial commit with no parent
				fmt.Println("  Initial commit - no diff available")
			}

			count++
		}
		return nil
	})
}

func GetRepoCommits(numCommits int) (object.CommitIter, error) {
	dir, err := os.MkdirTemp("", "repo-clone")
	if err != nil {
		fmt.Printf("Failed to create temp directory\n")
		return nil, err
	}
	defer os.RemoveAll(dir)

	fmt.Printf("Cloning repository\n")
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:          repoURL,
		SingleBranch: true,
		Depth:        numCommits + 1,
	})
	if err != nil {
		fmt.Printf("Failed to clone repository\n")
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		fmt.Printf("Failed to get HEAD\n")
		return nil, err
	}

	commits, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		fmt.Printf("Failed to get commit history: %v\n", err)
		return nil, err
	}

	return commits, nil
}

func GetRepoCommitsFrom(date time.Time) ([]*object.Commit, func(), error) {
	dir, err := os.MkdirTemp("", "repo-clone")
	if err != nil {
		fmt.Printf("Failed to create temp directory\n")
		return nil, nil, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return nil, cleanup, err
	}

	commits, err := repo.Log(&git.LogOptions{Since: &date})
	if err != nil {
		return nil, cleanup, err
	}

	commitsArray := []*object.Commit{}
	commits.ForEach(func(c *object.Commit) error {
		commitsArray = append(commitsArray, c)
		return nil
	})

	return commitsArray, cleanup, nil
}
