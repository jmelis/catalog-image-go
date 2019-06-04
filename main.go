package main

import (
	"os"

	"github.com/jmelis/catalog-image-go/pkg/catalog"
)

var operator = "hive"
var repo = "https://github.com/jmelis/test-catalog-image"
var username = "app"
var token = os.Getenv("GITHUB_TOKEN")
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"
var gitBranch = "master"
var gitDir = "/home/jmelis/borrar/tmpgit"

func checkIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	gitStoreOptions := catalog.GitStoreOptions{
		Operator:  operator,
		Repo:      repo,
		Username:  username,
		Token:     token,
		GitName:   gitName,
		GitEmail:  gitEmail,
		GitBranch: gitBranch,
		GitDir:    gitDir,
	}

	store, err := catalog.NewGitStore(gitStoreOptions)
	checkIfError(err)

	c := catalog.NewCatalog(operator, store)

	c.Load()
	c.Save()

	// err = c.WriteFile()
	// checkIfError(err)

	// err = c.Save()
	// checkIfError(err)
}
