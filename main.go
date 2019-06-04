package main

import (
	"log"
	"os"

	"github.com/jmelis/catalog-image-go/pkg/catalog"
)

// ../catalog-image/test/fixtures/bundles/0.1.506-14cff03

var operator = "hive"
var repo = "https://github.com/jmelis/test-catalog-image"
var username = "app"
var token = os.Getenv("GITHUB_TOKEN")
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"
var gitBranch = "master2"

// CheckIfError bla
func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	gitStoreOptions := catalog.GitStoreOptions{
		Repo:      repo,
		Username:  username,
		Token:     token,
		GitName:   gitName,
		GitEmail:  gitEmail,
		GitBranch: gitBranch,
	}

	store, err := catalog.NewGitStore(gitStoreOptions)
	CheckIfError(err)

	c := catalog.NewCatalog(operator, store)

	c.Load()
	// err = bundleStore.DeleteFile("b4/a")
	// CheckIfError(err)

	// err = bundleStore.Save()
	// CheckIfError(err)
}
