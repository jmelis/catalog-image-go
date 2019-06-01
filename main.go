package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var repo = "https://github.com/jmelis/test-catalog-image"
var username = "app"
var token = os.Getenv("GITHUB_TOKEN")
var directory = "tmp"
var bundleDir = "b3"
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"

// CheckIfError bla
func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: in-memory
// TODO: create branch if not exists
// TODO: remove dir
// TODO: add dir
// TODO: change file

func main() {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
		URL:      repo,
		Progress: os.Stdout,
	})
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// create bundle dir
	os.MkdirAll(filepath.Join(directory, bundleDir), os.ModePerm)

	// create file
	filename := filepath.Join(directory, bundleDir, "example-git-file")
	err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)
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
