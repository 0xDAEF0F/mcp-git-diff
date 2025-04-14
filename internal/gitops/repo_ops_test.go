package gitops

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	repoURL = "https://github.com/0xDAEF0F/whistle.git"
)

func TestGetRepo(t *testing.T) {
	_, cleanup, err := GetGitRepo(repoURL)
	if err != nil {
		t.Errorf("Failed to get repo: %v\n", err)
	}
	defer cleanup()

	fmt.Printf("Repo url: %v\n", repoURL)
}

func TestGetDiffWithOpts(t *testing.T) {
	repo, cleanup, err := GetGitRepo(repoURL)
	if err != nil {
		t.Fatalf("Failed to get repo: %v", err)
	}
	defer cleanup()

	_, err = GetDiffWithOpts(repo, &CommitOpts{
		Since: time.Now().AddDate(-1, 0, 0),
	})
	if err != nil {
		t.Errorf("Failed to get diff: %v\n", err)
	}
}

func TestGetContributors(t *testing.T) {
	repo, cleanup, err := GetGitRepo("https://github.com/0xDAEF0F/physical-nft.git")
	if err != nil {
		t.Fatalf("Failed to get repo: %v", err)
	}
	defer cleanup()

	contributors, err := GetContributors(repo)
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

func TestGetDiffWithSingleCommit(t *testing.T) {
	// 1. Create in-memory storage
	storer := memory.NewStorage()

	// 2. Create and store the file content blob
	blob := storer.NewEncodedObject()
	blob.SetType(plumbing.BlobObject)
	writer, err := blob.Writer()
	if err != nil {
		t.Fatalf("Failed to get blob writer: %v", err)
	}
	_, err = writer.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write blob content: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close blob writer: %v", err)
	}
	blobHash := blob.Hash()

	// 3. Create a tree entry for the file
	treeEntry := object.TreeEntry{
		Name: "README.md",
		Mode: 0100644, // Regular file mode
		Hash: blobHash,
	}

	// 4. Create and store the tree object
	tree := &object.Tree{
		Entries: []object.TreeEntry{treeEntry},
	}
	treeObject := storer.NewEncodedObject()
	err = tree.Encode(treeObject)
	if err != nil {
		t.Fatalf("Failed to encode tree object: %v", err)
	}
	treeHash := treeObject.Hash()

	// 5. Create and store the commit object
	commit := &object.Commit{
		Author: object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Committer: object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Message:  "Initial commit",
		TreeHash: treeHash,
	}
	commitObject := storer.NewEncodedObject()
	err = commit.Encode(commitObject)
	if err != nil {
		t.Fatalf("Failed to encode commit object: %v", err)
	}
	commitHash := commitObject.Hash()

	// 6. Initialize the repository using the storer.
	repo, err := git.Init(storer, nil)
	if err != nil {
		t.Fatalf("Failed to initialize in-memory repository: %v", err)
	}

	// Store the blob object (file content)
	_, err = repo.Storer.SetEncodedObject(blob)
	if err != nil {
		t.Fatalf("Failed to store blob object: %v", err)
	}

	// Store the tree object
	_, err = repo.Storer.SetEncodedObject(treeObject)
	if err != nil {
		t.Fatalf("Failed to store tree object: %v", err)
	}

	// Store the commit object itself
	_, err = repo.Storer.SetEncodedObject(commitObject)
	if err != nil {
		t.Fatalf("Failed to store commit object: %v", err)
	}

	// 7. Update HEAD to point to the created commit
	headRef := plumbing.NewHashReference(plumbing.HEAD, commitHash)
	err = repo.Storer.SetReference(headRef)
	if err != nil {
		t.Fatalf("Failed to set HEAD reference: %v", err)
	}

	diff, err := GetDiffWithOpts(repo, &CommitOpts{
		Since: time.Now().AddDate(-1, 0, 0),
	})
	if err != nil {
		t.Errorf("Failed to get diff: %v\n", err)
	}
	fmt.Printf("Diff: %v\n", diff)
	if diff == "" {
		t.Errorf("Expected diff to not be empty, but it was")
	}
	if !strings.Contains(diff, "README.md") {
		t.Errorf("Expected diff to contain 'README.md', but got: %s", diff)
	}
}
