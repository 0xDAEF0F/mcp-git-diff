package gitops

import (
	"fmt"
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

func GetDiffWithOpts(repo *git.Repository, opts *CommitOpts) (string, error) {
	fmt.Printf("Getting commits\n")

	logOpts := &git.LogOptions{
		Since: &opts.Since,
	}

	if opts.Branch != nil {
		ref, err := repo.Reference(plumbing.NewBranchReferenceName(*opts.Branch), true)
		if err != nil {
			fmt.Printf("Error getting branch reference: %v\n", err)
			return "", err
		}
		logOpts.From = ref.Hash()
	}

	commits, err := repo.Log(logOpts)
	if err != nil {
		fmt.Printf("Error getting commits: %v\n", err)
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

	// Only one commit in this range
	if lastCommit.Hash == firstCommit.Hash {
		prevCommit, err := lastCommit.Parent(0)
		// No last commit
		if err != nil {
			tree, _ := lastCommit.Tree()
			patch, err := (&object.Tree{}).Patch(tree)
			if err != nil {
				fmt.Printf("Error patching tree: %v\n", err)
				return "", err
			}
			return patch.String(), nil
		}
		firstCommit = prevCommit
	}

	patch, err := firstCommit.Patch(lastCommit)
	if err != nil {
		return "", err
	}

	return patch.String(), nil
}

// GetGitRepo clones a git repository into a temporary directory and returns the repository object and a cleanup function.
func GetGitRepo(repoUrl string, depth ...uint) (*git.Repository, func(), error) {
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

type Contributor struct {
	Name  string
	Email string
}

func GetContributors(repo *git.Repository) ([]Contributor, error) {
	commits, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, err
	}

	contributorsMap := make(map[string]Contributor)
	err = commits.ForEach(func(c *object.Commit) error {
		email := c.Author.Email
		if _, exists := contributorsMap[email]; !exists {
			contributorsMap[email] = Contributor{
				Name:  c.Author.Name,
				Email: email,
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	contributors := make([]Contributor, 0, len(contributorsMap))
	for _, contributor := range contributorsMap {
		contributors = append(contributors, contributor)
	}

	return contributors, nil
}
