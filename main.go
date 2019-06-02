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
var directory = "tmp"
var bundleDir = "b4"
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"
var gitBranch = "master"

// CheckIfError bla
func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: in-memory
// TODO: create branch if not exists (https://github.com/src-d/go-git/blob/master/_examples/branch/main.go)
// TODO: add dir function
// TODO: remove dir
// TODO: change file

func main() {
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

	w, err := r.Worktree()
	CheckIfError(err)

	// create bundle dir
	fs.MkdirAll(bundleDir, os.ModePerm)

	// create file
	file, err := fs.Create(filepath.Join(bundleDir, "example-git-file"))
	CheckIfError(err)

	_, err = file.Write([]byte("hello from memfs!"))
	CheckIfError(err)

	// Adds the bundleDir to the staging area.
	_, err = w.Add(bundleDir)
	CheckIfError(err)

	// Commit
	commitMsg := fmt.Sprintf("add bundledir %s", bundleDir)
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	CheckIfError(err)

	// Push
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})
	CheckIfError(err)
}

// // GitBundleStore TODO
// type GitBundleStore struct {
// 	r *git.Repository
// }

// // NewGitBundleStore TODO
// func NewGitBundleStore() GitBundleStore {
// 	storer := memory.NewStorage()
// 	fs := memfs.New()

// 	r, err := git.Clone(storer, fs, &git.CloneOptions{
// 		Auth: &http.BasicAuth{
// 			Username: username,
// 			Password: token,
// 		},
// 		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", gitBranch)),
// 		URL:           repo,
// 	})
// 	CheckIfError(err)

// 	return GitBundleStore{r}
// }
