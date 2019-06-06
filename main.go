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
			Name:  "git-repo",
			Usage: "git-repo",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "git-username",
			Usage: "git-username",
			Value: "app",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "git-token",
			Usage:  "git-token",
			EnvVar: "GIT_TOKEN",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "git-author-name",
			Usage: "git-author-name",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "git-author-email",
			Usage: "git-author-email",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "git-branch",
			Usage: "git-branch",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "git-work-dir",
			Usage: "git-work-dir",
		}),
		cli.BoolFlag{
			Name:  "git-no-push",
			Usage: "git-no-push",
		},
		cli.StringFlag{
			Name:  "prune",
			Usage: "prune descendants of `CSV`",
		},
		cli.StringFlag{
			Name:  "prune-hash",
			Usage: "prune descendants of `HASH`",
		},
		cli.StringFlag{
			Name:  "hash",
			Usage: "sets the hash for the bundle",
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
					Repo:        c.String("git-repo"),
					Username:    c.String("git-username"),
					Token:       c.String("git-token"),
					AuthorName:  c.String("git-author-name"),
					AuthorEmail: c.String("git-author-email"),
					Branch:      c.String("git-branch"),
					WorkDir:     c.String("git-work-dir"),
					Push:        !c.Bool("git-no-push"),
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
					if err := cl.PruneAfterCSV(pruneCSV); err != nil {
						return err
					}
				}

				if pruneHash := c.String("prune-hash"); pruneHash != "" {
					if err := cl.PruneAfterHash(pruneHash); err != nil {
						return err
					}
				}

				if err := cl.AddBundle(bundleDir, c.String("hash")); err != nil {
					return err
				}

				return cl.Save()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
