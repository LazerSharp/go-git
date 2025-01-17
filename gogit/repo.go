package gogit

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/LazerSharp/go-git/config"
)

type GitRepository struct {
	WorkTree string
	GitDir   string
}

func NewGitRepository(pth string, force bool) (*GitRepository, error) {

	if !force && !isDir(pth) {
		return nil, fmt.Errorf("%s is not a Directory", pth)
	}

	gitDir := path.Join(pth, ".git")

	return &GitRepository{
		WorkTree: pth,
		GitDir:   gitDir,
	}, nil
}

func RepoPath(r *GitRepository, pth ...string) string {
	args := append([]string{r.GitDir}, pth...)
	return path.Join(args...)
}

func RepoFile(r *GitRepository, mkDir bool, pth ...string) (string, error) {
	slog.Debug("RepoFile paths --->", "path", pth)
	_, err := RepoDir(r, mkDir, pth[:len(pth)-1]...)
	if err != nil {
		return "", err
	}
	return RepoPath(r, pth...), nil
}

func RepoDir(r *GitRepository, mkDir bool, pth ...string) (string, error) {
	slog.Debug("RepoDir paths --->", "path", pth)
	p := RepoPath(r, pth...)
	f, err := os.Stat(p) // get File Info
	if err == nil {      // file / dir exists
		if f.IsDir() {
			return p, nil
		} else {
			// it is not a dir. file maybe
			return "", fmt.Errorf("%s is not a Directory", p)
		}
	}

	if mkDir {
		slog.Debug("making dir", "path", p)
		err = os.MkdirAll(p, 0700)
		if err != nil {
			return "", err
		}
	}
	return p, nil
}

func RepoCreate(pth string) error {
	repo, err := NewGitRepository(pth, true)
	if err != nil {
		return err
	}
	if ewt, wt := exists(repo.WorkTree); ewt {
		if !wt.IsDir() {
			return fmt.Errorf("%s is not a Directory", repo.WorkTree)
		}
		if egd, gd := exists(repo.GitDir); egd {
			if !gd.IsDir() {
				return fmt.Errorf("%s is not a Directory", repo.GitDir)
			}
			if empty, _ := isDirEmpty(repo.GitDir); !empty {
				return fmt.Errorf("%s Directory is not empty", repo.GitDir)
			}
		}
	} else {
		if err = os.Mkdir(repo.WorkTree, 0644); err != nil {
			return err
		}
	}

	Must(RepoDir(repo, true, "branches"))
	Must(RepoDir(repo, true, "objects"))
	Must(RepoDir(repo, true, "refs", "tags"))
	Must(RepoDir(repo, true, "refs", "heads"))

	Check(os.WriteFile(Must(RepoFile(repo, false, "description")),
		[]byte("Unnamed repository; edit this file 'description' to name the repository.\n"),
		0644))
	Check(os.WriteFile(Must(RepoFile(repo, false, "HEAD")),
		[]byte("ref: refs/heads/master\n"),
		0644))

	// write git config file
	cfgPath := Must(RepoFile(repo, false, "config"))
	f := Must(os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY, 0644))
	defer f.Close()
	cfg := config.DeaultConfig()
	cfg.Write(f)
	return nil
}

func CatFile(typ, obj string) {
	repo := Must(NewGitRepository(".", false))
	o := Must(ReadObject(repo, obj))
	if o == nil {
		log.Fatal("Unable to read!")
	}
	o.Serialize(os.Stdout)
}

func (r *GitRepository) Checkout(cHash string) (err error) {

	tmpDir := "chkout"
	tpath := []string{"temp", tmpDir}
	td, err := RepoFile(r, false, tpath...)
	if err != nil {
		return err
	}
	defer func() {
		err = os.RemoveAll(td)
	}()
	err = checkoutCommit(r, cHash, tpath...)
	if err != nil {
		return err
	}

	// remove workdir content
	fpaths, err := filepath.Glob(filepath.Join(r.WorkTree, "*"))
	if err != nil {
		return err
	}

	for _, p := range fpaths {
		if p == r.GitDir {
			continue
		}
		err = os.RemoveAll(p)
		if err != nil {
			return err
		}
	}
	// copy files / folders from td

	tpaths, err := filepath.Glob(filepath.Join(td, "*"))
	if err != nil {
		return err
	}

	for _, tp := range tpaths {
		_, f := filepath.Split(tp)

		fmt.Printf("Moving %s to %s", tp, filepath.Join(r.WorkTree, f))
		os.Rename(tp, filepath.Join(r.WorkTree, f))
	}

	return err
}

func checkoutCommit(repo *GitRepository, cHash string, pth ...string) error {

	obj, err := ReadObject(repo, cHash)
	if err != nil {
		return err
	}
	if obj.Type() != "commit" {
		return fmt.Errorf("invalid commit #: %s", cHash)
	}

	commit := obj.(*Commit)

	return checkoutTree(repo, *commit.Tree, pth...)
}

func checkoutTree(repo *GitRepository, tHash string, pth ...string) error {

	obj, err := ReadObject(repo, tHash)
	if err != nil {
		return err
	}
	if obj.Type() != "tree" {
		return fmt.Errorf("invalid tree #: %s", tHash)
	}

	tree := obj.(*Tree)

	for _, entr := range tree.Entries {
		switch entr.Type {
		case BlobEntry:
			p := append(pth, entr.Path)
			err = checkoutBlob(repo, entr.Hash, p...)
		case TreeEntry:
			p := append(pth, entr.Path)
			err = checkoutTree(repo, entr.Hash, p...)
		default:
			err = fmt.Errorf("invalid type %s", entr.Type)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func checkoutBlob(repo *GitRepository, bHash string, pth ...string) error {

	// open work tree (code) file
	fpth, err := RepoFile(repo, true, pth...)
	if err != nil {
		return err
	}
	f, err := os.Create(fpth)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		Check(f.Close())
	}(f)

	// read objevt from git
	o, err := ReadObject(repo, bHash)
	if err != nil {
		return err
	}
	if o.Type() != "blob" {
		return fmt.Errorf("%s is not a file", fpth)
	}

	// stream file comtent from Blob object to work space file
	return o.Serialize(f)
}
