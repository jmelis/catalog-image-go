package main

import (
	"os"

	"github.com/jmelis/catalog-image-go/pkg/catalog"
)

var operator = "hive"
var channel = "staging"
var repo = "https://github.com/jmelis/test-catalog-image"
var username = "app"
var token = os.Getenv("GITHUB_TOKEN")
var gitName = "Jaime Melis"
var gitEmail = "j.melis@gmail.com"
var gitBranch = "master"
var gitDir = "/home/jmelis/borrar/tmpgit"
var bundlePath = "/home/jmelis/work/git/catalog-image/test/fixtures/bundles/0.1.700-0000000"

func checkIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	gitStoreOptions := catalog.GitStoreOptions{
		Operator:  operator,
		Channel:   channel,
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

	c, err := catalog.LoadCatalog(store)
	checkIfError(err)

	err = c.PruneAfterCSV("hive-operator.v0.1.598-1af4d6f")
	checkIfError(err)

	// err = c.Bundles.PruneAfterCSV("hive-operator.v0.1.598-1af4d6f")

	// err = c.AddBundle(bundlePath)
	// checkIfError(err)

	err = c.Save()
	checkIfError(err)
}
