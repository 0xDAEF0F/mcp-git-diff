package gitops

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetRepoCommitsFrom retrieves commits from a repository since a given date
func GetRepoCommitsFrom(repoUrl string, date time.Time) ([]*object.Commit, func(), error) {
	repo, cleanup, err := GetRepo(repoUrl)
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

// GetRepoCommits retrieves a specific number of commits from the repository
func GetRepoCommits(repoUrl string, numCommits int) ([]*object.Commit, func(), error) {
	repo, cleanup, err := GetRepo(repoUrl, uint(numCommits))
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

func GetRepo(repoUrl string, depth ...uint) (*git.Repository, func(), error) {
	dir, err := os.MkdirTemp("", "repo-clone")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	cloneOpts := &git.CloneOptions{
		URL: repoUrl,
	}

	if len(depth) > 0 {
		cloneOpts.Depth = int(depth[0])
	}

	repo, err := git.PlainClone(dir, false, cloneOpts)
	if err != nil {
		return nil, cleanup, err
	}

	return repo, cleanup, nil
}
