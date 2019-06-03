package bundlestore

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GitBundleStoreOptions TODO
type GitBundleStoreOptions struct {
	Repo      string
	Username  string
	Token     string
	GitName   string
	GitEmail  string
	GitBranch string
}

// GitBundleStore TODO
type GitBundleStore struct {
	r       *git.Repository
	options GitBundleStoreOptions
}

// NewGitBundleStore TODO
func NewGitBundleStore(options GitBundleStoreOptions) (*GitBundleStore, error) {
	storer := memory.NewStorage()
	fs := memfs.New()
	origin := "origin"

	r, err := git.Init(storer, fs)
	if err != nil {
		return nil, err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: origin,
		URLs: []string{options.Repo},
	})
	if err != nil {
		return nil, err
	}

	err = r.Fetch(&git.FetchOptions{
		RemoteName: origin,
		Auth: &http.BasicAuth{
			Username: options.Username,
			Password: options.Token,
		},
	})

	refHeadHash := plumbing.Hash{}

	refs, _ := r.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			remoteRefName := fmt.Sprintf("refs/remotes/origin/%s", options.GitBranch)
			if ref.Name() == plumbing.ReferenceName(remoteRefName) {
				refHeadHash = ref.Hash()
			}
		}
		return nil
	})

	refName := fmt.Sprintf("refs/heads/%s", options.GitBranch)
	if refHeadHash == plumbing.ZeroHash {
		refHead := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName(refName))
		err = r.Storer.SetReference(refHead)
		if err != nil {
			return nil, err
		}
	} else {
		refBranch := plumbing.NewHashReference(plumbing.ReferenceName(refName), refHeadHash)
		err = r.Storer.SetReference(refBranch)
		if err != nil {
			return nil, err
		}

		refHead := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName(refName))
		err = r.Storer.SetReference(refHead)
		if err != nil {
			return nil, err
		}

		// Checkout
		w, err := r.Worktree()
		if err != nil {
			return nil, err
		}

		err = w.Reset(&git.ResetOptions{
			Mode:   git.HardReset,
			Commit: refHeadHash,
		})
		if err != nil {
			return nil, err
		}
	}

	return &GitBundleStore{r: r, options: options}, nil
}

// DeleteFile TODO
func (g *GitBundleStore) DeleteFile(path string) error {
	w, err := g.r.Worktree()
	if err != nil {
		return err
	}

	// Adds the path to the staging area.
	_, err = w.Remove(path)
	return nil
}

// WriteFile TODO
func (g *GitBundleStore) WriteFile(path string, content []byte) error {
	w, err := g.r.Worktree()
	if err != nil {
		return err
	}

	// create bundle dir
	dirPath := filepath.Dir(path)
	baseName := filepath.Base(path)
	w.Filesystem.MkdirAll(dirPath, os.ModePerm)

	// create file
	file, err := w.Filesystem.Create(filepath.Join(dirPath, baseName))
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	// Adds the path to the staging area.
	_, err = w.Add(path)
	return err
}

// Save TODO
func (g *GitBundleStore) Save() error {
	w, err := g.r.Worktree()
	if err != nil {
		return err
	}

	// Commit
	commitMsg := fmt.Sprintf("commit")
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.options.GitName,
			Email: g.options.GitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// Push
	err = g.r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: g.options.Username,
			Password: g.options.Token,
		},
	})
	return err
}
