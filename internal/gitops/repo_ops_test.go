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
	storer := memory.NewStorage()
	blobHash := createBlob(t, storer, "Hello, World!")
	treeEntry := object.TreeEntry{Name: "README.md", Mode: 0100644, Hash: blobHash}
	treeHash := createTree(t, storer, treeEntry)
	commitHash := createCommit(t, storer, treeHash)

	repo, err := git.Init(storer, nil)
	if err != nil {
		t.Fatalf("Failed to initialize in-memory repository: %v", err)
	}
	headRef := plumbing.NewHashReference(plumbing.HEAD, commitHash)
	if err := repo.Storer.SetReference(headRef); err != nil {
		t.Fatalf("Failed to set HEAD reference: %v", err)
	}

	diff, err := GetDiffWithOpts(repo, &CommitOpts{Since: time.Now().AddDate(-1, 0, 0)})
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

// Helper functions for creating in-memory objects
func createBlob(t *testing.T, storer *memory.Storage, content string) plumbing.Hash {
	t.Helper()
	blob := storer.NewEncodedObject()
	blob.SetType(plumbing.BlobObject)
	writer, err := blob.Writer()
	if err != nil {
		t.Fatalf("Failed to get blob writer: %v", err)
	}
	_, err = writer.Write([]byte(content))
	if err != nil {
		t.Fatalf("Failed to write blob content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close blob writer: %v", err)
	}
	_, err = storer.SetEncodedObject(blob)
	if err != nil {
		t.Fatalf("Failed to store blob object: %v", err)
	}
	return blob.Hash()
}

func createTree(t *testing.T, storer *memory.Storage, entry object.TreeEntry) plumbing.Hash {
	t.Helper()
	tree := &object.Tree{Entries: []object.TreeEntry{entry}}
	treeObj := storer.NewEncodedObject()
	if err := tree.Encode(treeObj); err != nil {
		t.Fatalf("Failed to encode tree object: %v", err)
	}
	_, err := storer.SetEncodedObject(treeObj)
	if err != nil {
		t.Fatalf("Failed to store tree object: %v", err)
	}
	return treeObj.Hash()
}

func createCommit(t *testing.T, storer *memory.Storage, treeHash plumbing.Hash) plumbing.Hash {
	t.Helper()
	now := time.Now()
	commit := &object.Commit{
		Author:    object.Signature{Name: "Test User", Email: "test@example.com", When: now},
		Committer: object.Signature{Name: "Test User", Email: "test@example.com", When: now},
		Message:   "Initial commit",
		TreeHash:  treeHash,
	}
	commitObj := storer.NewEncodedObject()
	if err := commit.Encode(commitObj); err != nil {
		t.Fatalf("Failed to encode commit object: %v", err)
	}
	_, err := storer.SetEncodedObject(commitObj)
	if err != nil {
		t.Fatalf("Failed to store commit object: %v", err)
	}
	return commitObj.Hash()
}
