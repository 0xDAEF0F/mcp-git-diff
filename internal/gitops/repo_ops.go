package gitops

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitOpts struct {
	Since       time.Time
	Branch      *string
	AuthorEmail *string
}

func GetDiffWithOpts(repoUrl string, opts *CommitOpts) (string, error) {
	repo, cleanup, err := getGitRepo(repoUrl)
	defer cleanup()
	if err != nil {
		return "", err
	}

	logOpts := &git.LogOptions{
		Since: &opts.Since,
	}

	if opts.Branch != nil {
		ref, err := repo.Reference(plumbing.NewBranchReferenceName(*opts.Branch), true)
		if err != nil {
			return "", err
		}
		logOpts.From = ref.Hash()
	}

	commits, err := repo.Log(logOpts)
	if err != nil {
		return "", err
	}

	commitsArray := []*object.Commit{}
	commits.ForEach(func(c *object.Commit) error {
		if opts.AuthorEmail == nil || c.Author.Email == *opts.AuthorEmail {
			commitsArray = append(commitsArray, c)
		}
		return nil
	})

	lastCommit := commitsArray[0]
	firstCommit := commitsArray[len(commitsArray)-1]

	patch, err := firstCommit.Patch(lastCommit)
	if err != nil {
		return "", err
	}

	return patch.String(), nil
}

func getGitRepo(repoUrl string, depth ...uint) (*git.Repository, func(), error) {
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
