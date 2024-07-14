package gogit

import (
	"fmt"
	"log/slog"
	"os"
	"path"

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
