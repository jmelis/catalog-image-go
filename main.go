package main

import (
	"fmt"
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

// var gitDir = "/home/jmelis/borrar/tmpgit"
var bundlePath = "/home/jmelis/work/git/catalog-image/test/fixtures/bundles/0.1.700-0000000"

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
		// GitDir:    gitDir,
	}

	store, err := catalog.NewGitStore(gitStoreOptions)
	checkIfError(err)

	c := catalog.NewCatalog(operator, store)

	c.Load()
	for _, b := range c.Bundles {
		fmt.Println(b.CSV.Version())
	}
	fmt.Println("hi")
	c.AddBundle(bundlePath)
	for _, b := range c.Bundles {
		fmt.Println(b.CSV.Version())
	}

	// c.Save()

	// err = c.WriteFile()
	// checkIfError(err)

	// err = c.Save()
	// checkIfError(err)
}
