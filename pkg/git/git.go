package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GetRepo opens the git repository in the current directory
func GetRepo() (*git.Repository, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	return git.PlainOpen(dir)
}

// BranchExists checks if a branch exists in the repository
func BranchExists(repo *git.Repository, branch string) bool {
	_, err := repo.Reference(plumbing.NewBranchReferenceName(branch), true)
	return err == nil
}

// CheckoutBranch checks out a branch in the repository
func CheckoutBranch(repo *git.Repository, branch string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}
	return wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)})
}

// CreateBranch creates a new branch from the current HEAD
func CreateBranch(repo *git.Repository, branch string) error {
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("error getting HEAD: %v", err)
	}
	ref := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+branch), head.Hash())
	return repo.Storer.SetReference(ref)
}
