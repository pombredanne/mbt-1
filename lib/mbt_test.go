package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"text/template"
	"time"

	git "github.com/libgit2/git2go"
)

type TestRepository struct {
	Dir           string
	Repo          *git.Repository
	LastCommit    *git.Oid
	CurrentBranch string
}

type TestApplication struct {
	Name           string
	Build          string
	BuildPlatforms []string
	Properties     map[string]string
}

func (r *TestRepository) InitApplication(p string) error {
	return r.InitApplicationWithOptions(p, &TestApplication{
		Name:           path.Base(p),
		Build:          "./build.sh",
		BuildPlatforms: []string{"darwin", "linux"},
		Properties:     map[string]string{"foo": "bar", "jar": "car"},
	})
}

func (r *TestRepository) InitApplicationWithOptions(p string, app *TestApplication) error {
	appDir := path.Join(r.Dir, p)
	err := os.MkdirAll(appDir, 0755)
	if err != nil {
		return err
	}

	t, err := template.New("appspec").Parse(`name: {{ .Name }}
buildPlatforms: 
  {{ range $p := .BuildPlatforms }}- {{ $p }}
  {{ end }}
build: {{ .Build }}
{{ if .Properties }}
properties: 
  {{ range $k, $v := .Properties }}{{ $k }}: {{ $v }}
  {{ end }}
{{ end }}
`)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, app)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(appDir, "appspec.yaml"), buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *TestRepository) WriteContent(file, content string) error {
	fpath := path.Join(r.Dir, file)
	dir := path.Dir(fpath)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(fpath, []byte(content), 0744)
}

func (r *TestRepository) Commit(message string) error {
	idx, err := r.Repo.Index()
	if err != nil {
		return err
	}

	err = idx.AddAll([]string{"."}, git.IndexAddCheckPathspec, func(p string, f string) int {
		return 0
	})
	if err != nil {
		return err
	}

	oid, err := idx.WriteTree()
	if err != nil {
		return err
	}

	tree, err := r.Repo.LookupTree(oid)
	if err != nil {
		return err
	}

	sig := &git.Signature{
		Email: "alice@wonderland.com",
		Name:  "alice",
		When:  time.Now(),
	}

	parents := []*git.Commit{}
	isEmpty, err := r.Repo.IsEmpty()
	if err != nil {
		return nil
	}

	if !isEmpty {
		currentBranch, err := r.Repo.Head()
		if err != nil {
			return err
		}

		bc, err := r.Repo.LookupCommit(currentBranch.Target())
		if err != nil {
			return err
		}

		parents = append(parents, bc)
	}

	r.LastCommit, err = r.Repo.CreateCommit("HEAD", sig, sig, message, tree, parents...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TestRepository) SwitchToBranch(name string) error {
	_, err := r.Repo.LookupBranch(name, git.BranchAll)
	if err != nil {
		head, err := r.Repo.Head()
		if err != nil {
			return err
		}

		hc, err := r.Repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		_, err = r.Repo.CreateBranch(name, hc, false)
		if err != nil {
			return err
		}
	}

	err = r.Repo.SetHead(fmt.Sprintf("refs/heads/%s", name))
	if err != nil {
		return err
	}

	return r.Repo.CheckoutHead(&git.CheckoutOpts{
		Strategy: git.CheckoutForce | git.CheckoutRemoveUntracked | git.CheckoutDontWriteIndex,
	})
}

func (r *TestRepository) Remove(p string) error {
	return os.RemoveAll(path.Join(r.Dir, p))
}

func (r *TestRepository) Rename(old, new string) error {
	return os.Rename(path.Join(r.Dir, old), path.Join(r.Dir, new))
}

func createTestRepository(dir string) (*TestRepository, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	repo, err := git.InitRepository(dir, false)
	if err != nil {
		return nil, err
	}

	return &TestRepository{dir, repo, nil, "master"}, nil
}

func clean() {
	os.RemoveAll(".tmp")
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
