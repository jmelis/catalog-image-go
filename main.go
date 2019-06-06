package main

import (
	"log"
	"os"

	"github.com/jmelis/catalog-image-go/pkg/catalog"
	"github.com/urfave/cli/altsrc"
	"gopkg.in/urfave/cli.v1"
)

func main() {

	addFlags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "operator",
			Usage: "operator",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "channel",
			Usage: "channel",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "gitRepo",
			Usage: "gitRepo",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "gitUsername",
			Usage: "gitUsername",
			Value: "app",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "gitToken",
			Usage:  "gitToken",
			EnvVar: "GIT_TOKEN",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "gitAuthorName",
			Usage: "gitAuthorName",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "gitAuthorEmail",
			Usage: "gitAuthorEmail",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "gitBranch",
			Usage: "gitBranch",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "gitWorkDir",
			Usage: "gitWorkDir",
		}),
		cli.StringFlag{
			Name:  "prune",
			Usage: "prune descendants of `CSV` (example: hive-operator.v0.1.506-14cff03)",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
	}

	app := cli.NewApp()
	app.Name = "catalog-image"
	app.Usage = "manage grpc catalogs for OLM"
	app.Commands = []cli.Command{
		{
			Name:      "add",
			Usage:     "add bundle to the repo",
			UsageText: "add BUNDLE_DIR",
			Flags:     addFlags,
			Before:    altsrc.InitInputSourceWithContext(addFlags, altsrc.NewYamlSourceFromFlagFunc("config")),
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return cli.NewExitError("missing argument: BUNDLE_DIR", 1)
				}

				bundleDir := c.Args().Get(0)

				stat, err := os.Stat(bundleDir)
				if err != nil {
					return err
				}
				if !stat.IsDir() {
					return cli.NewExitError("invalid argument is not a directory: BUNDLE_DIR", 1)
				}

				gitStoreOptions := catalog.GitStoreOptions{
					Operator:    c.String("operator"),
					Channel:     c.String("channel"),
					Repo:        c.String("gitRepo"),
					Username:    c.String("gitUsername"),
					Token:       c.String("gitToken"),
					AuthorName:  c.String("gitAuthorName"),
					AuthorEmail: c.String("gitAuthorEmail"),
					Branch:      c.String("gitBranch"),
					WorkDir:     c.String("gitWorkDir"),
				}

				store, err := catalog.NewGitStore(gitStoreOptions)
				if err != nil {
					return err
				}

				cl, err := catalog.LoadCatalog(store)
				if err != nil {
					return err
				}

				if pruneCSV := c.String("prune"); pruneCSV != "" {
					err = cl.PruneAfterCSV(pruneCSV)
					if err != nil {
						return err
					}
				}

				err = cl.AddBundle(bundleDir)
				if err != nil {
					return err
				}

				err = cl.Save()
				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
