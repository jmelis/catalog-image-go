package catalog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GitStoreOptions TODO
type GitStoreOptions struct {
	Operator  string
	Repo      string
	Username  string
	Token     string
	GitName   string
	GitEmail  string
	GitBranch string
	// GitDir cloned repo path. If empty it will clone in memory.
	GitDir string
}

// GitStore TODO
type GitStore struct {
	r       *git.Repository
	options GitStoreOptions
}

// NewGitStore TODO
func NewGitStore(options GitStoreOptions) (*GitStore, error) {
	var fs billy.Filesystem
	var storer storage.Storer

	if options.GitDir != "" {
		fs = osfs.New(options.GitDir)
		dot, _ := fs.Chroot(".git")
		storer = filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	} else {
		storer = memory.NewStorage()
		fs = memfs.New()
	}

	origin := "origin"

	r, err := git.Init(storer, fs)
	if err != nil {
		return nil, err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: origin,
		URLs: []string{options.Repo},
	})
	if err != nil {
		return nil, err
	}

	err = r.Fetch(&git.FetchOptions{
		RemoteName: origin,
		Auth: &http.BasicAuth{
			Username: options.Username,
			Password: options.Token,
		},
	})

	refHeadHash := plumbing.Hash{}

	refs, _ := r.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			remoteRefName := fmt.Sprintf("refs/remotes/origin/%s", options.GitBranch)
			if ref.Name() == plumbing.ReferenceName(remoteRefName) {
				refHeadHash = ref.Hash()
			}
		}
		return nil
	})

	refName := fmt.Sprintf("refs/heads/%s", options.GitBranch)
	if refHeadHash == plumbing.ZeroHash {
		refHead := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName(refName))
		err = r.Storer.SetReference(refHead)
		if err != nil {
			return nil, err
		}
	} else {
		refBranch := plumbing.NewHashReference(plumbing.ReferenceName(refName), refHeadHash)
		err = r.Storer.SetReference(refBranch)
		if err != nil {
			return nil, err
		}

		refHead := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName(refName))
		err = r.Storer.SetReference(refHead)
		if err != nil {
			return nil, err
		}

		// Checkout
		w, err := r.Worktree()
		if err != nil {
			return nil, err
		}

		err = w.Reset(&git.ResetOptions{
			Mode:   git.HardReset,
			Commit: refHeadHash,
		})
		if err != nil {
			return nil, err
		}
	}

	return &GitStore{r: r, options: options}, nil
}

func (g *GitStore) deleteFile(path string) error {
	w, err := g.r.Worktree()
	if err != nil {
		return err
	}

	// Adds the path to the staging area.
	_, err = w.Remove(path)
	return nil
}

func (g *GitStore) writeFile(path string, content []byte) error {
	w, err := g.r.Worktree()
	if err != nil {
		return err
	}

	// create leading dirs
	dirPath := filepath.Dir(path)
	baseName := filepath.Base(path)
	w.Filesystem.MkdirAll(dirPath, os.ModePerm)

	// create file
	file, err := w.Filesystem.Create(filepath.Join(dirPath, baseName))
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	// Adds the path to the staging area.
	_, err = w.Add(path)
	return err
}

func (g *GitStore) readFile(path string) ([]byte, error) {
	w, err := g.r.Worktree()
	if err != nil {
		return nil, err
	}

	fs := w.Filesystem

	fd, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	// read file
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (g *GitStore) load() (Bundles, error) {
	operator := g.options.Operator

	w, err := g.r.Worktree()
	if err != nil {
		return nil, err
	}

	fs := w.Filesystem

	files, err := fs.ReadDir(operator)
	if err != nil {
		return nil, err
	}

	var bundles Bundles
	for _, bundleDir := range files {
		if bundleDir.IsDir() {
			dirPath := filepath.Join(operator, bundleDir.Name())

			// read csv
			csvFilePath := filepath.Join(dirPath, CSVName(operator, bundleDir.Name()))

			content, err := g.readFile(csvFilePath)
			if err != nil {
				return nil, err
			}

			csv := CSV{
				version: bundleDir.Name(),
				content: content,
			}

			// read rest of files
			sideFiles, err := fs.ReadDir(dirPath)
			if err != nil {
				return nil, err
			}

			var sidefiles []SideFile
			for _, sideFile := range sideFiles {
				sideFilePath := filepath.Join(dirPath, sideFile.Name())

				if strings.HasSuffix(sideFile.Name(), CSVSuffix) {
					continue
				}

				if !strings.HasSuffix(sideFile.Name(), ".yaml") {
					return nil, fmt.Errorf("only '.yaml' is supported")
				}

				content, err := g.readFile(sideFilePath)
				if err != nil {
					return nil, err
				}

				sidefile := SideFile{
					name:    sideFile.Name(),
					content: content,
				}
				sidefiles = append(sidefiles, sidefile)
			}

			bundles = append(bundles, Bundle{csv: csv, sidefiles: sidefiles})
		}
	}

	return bundles, nil
}

func (g *GitStore) save(bundles Bundles) error {
	for _, bundle := range bundles {
		bundleDir := filepath.Join(g.options.Operator, bundle.csv.version)

		// write sidefiles
		for _, sf := range bundle.sidefiles {
			sfPath := filepath.Join(bundleDir, sf.name)
			g.writeFile(sfPath, sf.content)
		}

		// write csv
		csvPath := filepath.Join(bundleDir, CSVName(g.options.Operator, bundle.csv.version))
		g.writeFile(csvPath, bundle.csv.content)
	}

	// TODO: create packagefile
	return g.commit()
}

func (g *GitStore) commit() error {
	w, err := g.r.Worktree()
	if err != nil {
		return err
	}

	// Commit
	commitMsg := fmt.Sprintf("commit1")
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.options.GitName,
			Email: g.options.GitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// // Push
	// err = g.r.Push(&git.PushOptions{
	// 	Auth: &http.BasicAuth{
	// 		Username: g.options.Username,
	// 		Password: g.options.Token,
	// 	},
	// })
	return err
}
