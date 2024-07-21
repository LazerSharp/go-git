package gogit

import (
	"fmt"
	"testing"
)

func TestRepoFile(t *testing.T) {
	repo, err := NewGitRepository("testdata", false)
	if err != nil {
		t.Fatal(err)
	}
	f, err := RepoFile(repo, false, "temps")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("file:", f)

}

func TestCheckoutCommit(t *testing.T) {
	repo, err := NewGitRepository("testdata", false)
	if err != nil {
		t.Fatal(err)
	}
	err = checkoutCommit(
		repo,
		"d8050a3db8f5a4e8ff9d7da7b90bf0948d0bec4b",
		"temp", "chkout",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckout(t *testing.T) {
	repo, err := NewGitRepository("testdata", false)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.Checkout(
		"d8050a3db8f5a4e8ff9d7da7b90bf0948d0bec4b",
	)
	if err != nil {
		t.Fatal(err)
	}
}
