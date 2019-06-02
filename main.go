package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var repo = "https://github.com/jmelis/test-catalog-image"
var username = "app"
var token = os.Getenv("GITHUB_TOKEN")
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"
var gitBranch = "master"

// CheckIfError bla
func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// Qs:
// - New? Pointer to struct?
// - Interface? Register with init() ?
// - Create branch?
// - Chmod file?
// - GitBundleStoreOptions. Pointer to struct?

// TODO: create branch if not exists (https://github.com/src-d/go-git/blob/master/_examples/branch/main.go)
// TODO: add dir function
// TODO: remove dir
// TODO: change file

func main() {
	bundleStore := NewGitBundleStore()

	bundleStore.AddFile("b5/a", []byte("hello AddFile1"))
	bundleStore.AddFile("ttt", []byte("hello AddFile2"))

	bundleStore.Save()
}

// GitBundleStore TODO
type GitBundleStore struct {
	r *git.Repository
}

// NewGitBundleStore TODO
func NewGitBundleStore() GitBundleStore {
	storer := memory.NewStorage()
	fs := memfs.New()

	r, err := git.Clone(storer, fs, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", gitBranch)),
		URL:           repo,
	})
	CheckIfError(err)

	return GitBundleStore{r}
}

// AddFile TOOD
func (g GitBundleStore) AddFile(path string, content []byte) {
	w, err := g.r.Worktree()
	CheckIfError(err)

	// get Filesystem
	fs := w.Filesystem

	// create bundle dir
	dirPath := filepath.Dir(path)
	baseName := filepath.Base(path)

	fs.MkdirAll(dirPath, os.ModePerm)

	// create file
	file, err := fs.Create(filepath.Join(dirPath, baseName))
	CheckIfError(err)

	_, err = file.Write(content)
	CheckIfError(err)

	// Adds the path to the staging area.
	_, err = w.Add(path)
	CheckIfError(err)
}

// Save TODO
func (g GitBundleStore) Save() error {
	w, err := g.r.Worktree()
	CheckIfError(err)

	// Commit
	commitMsg := fmt.Sprintf("commit")
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	CheckIfError(err)

	// Push
	err = g.r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})
	CheckIfError(err)

	return nil
}
