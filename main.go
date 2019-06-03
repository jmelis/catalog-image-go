package main

import (
	"os"

	"github.com/jmelis/catalog-image-go/pkg/bundlestore"
)

var repo = "https://github.com/jmelis/test-catalog-image"
var username = "app"
var token = os.Getenv("GITHUB_TOKEN")
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"
var gitBranch = "master2"

// CheckIfError bla
func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	gitBundleStoreOptions := bundlestore.GitBundleStoreOptions{
		Repo:      repo,
		Username:  username,
		Token:     token,
		GitName:   gitName,
		GitEmail:  gitEmail,
		GitBranch: gitBranch,
	}

	bundleStore, err := bundlestore.NewGitBundleStore(gitBundleStoreOptions)
	CheckIfError(err)

	err = bundleStore.AddFile("b4/a", []byte("hello AddFile1"))
	CheckIfError(err)

	err = bundleStore.Save()
	CheckIfError(err)
}
