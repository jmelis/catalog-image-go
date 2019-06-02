package gitbundlestore

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

// Options TODO
type Options struct {
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
	options Options
}

// NewGitBundleStore TODO
func NewGitBundleStore(options Options) (*GitBundleStore, error) {
	storer := memory.NewStorage()
	fs := memfs.New()

	r, err := git.Init(storer, fs)
	if err != nil {
		return nil, err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{options.Repo},
	})
	if err != nil {
		return nil, err
	}

	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: options.Username,
			Password: options.Token,
		},
	})

	h := plumbing.Hash{}

	refs, _ := r.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			remoteRefName := fmt.Sprintf("refs/remotes/origin/%s", options.GitBranch)
			if ref.Name() == plumbing.ReferenceName(remoteRefName) {
				h = ref.Hash()
			}
		}
		return nil
	})

	if h == plumbing.ZeroHash {
		// TODO: implement
		panic("NOPE")
	}

	refName := fmt.Sprintf("refs/heads/%s", options.GitBranch)
	ref := plumbing.NewHashReference(plumbing.ReferenceName(refName), h)

	// The created reference is saved in the storage.
	err = r.Storer.SetReference(ref)

	// Checkout
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	err = w.Reset(&git.ResetOptions{
		Mode:   git.MergeReset,
		Commit: head.Hash(),
	})
	if err != nil {
		return nil, err
	}

	return &GitBundleStore{r: r, options: options}, nil
}

// AddFile TODO
func (g *GitBundleStore) AddFile(path string, content []byte) error {
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
