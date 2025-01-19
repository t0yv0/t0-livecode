package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type codeRepo interface {
	Default() Pid
	List() ([]Pid, error)
	// Search(query string) []Pid
	Store(id Pid, code string) error
	Load(pid Pid) (string, error)
}

type dirStore struct {
	dir string
}

var _ codeRepo = (*dirStore)(nil)

func (i *dirStore) suffix() string {
	return ".js"
}

func (i *dirStore) path(id Pid) string {
	return filepath.Join(i.dir, string(id)+i.suffix())
}

func (i *dirStore) List() ([]Pid, error) {
	var pids []Pid
	entries, err := os.ReadDir(i.dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), i.suffix()) {
			continue
		}
		pids = append(pids, Pid(strings.TrimSuffix(e.Name(), i.suffix())))
	}
	sort.Slice(pids, func(i, j int) bool {
		return pids[i] < pids[j]
	})
	return pids, nil
}

func (i *dirStore) Load(id Pid) (string, error) {
	b, err := os.ReadFile(i.path(id))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (i *dirStore) Default() Pid {
	return "starter"
}

func (i *dirStore) Store(id Pid, code string) error {
	return os.WriteFile(i.path(id), []byte(code), 0655)
}

func newDirStore(dir string) (*dirStore, error) {
	s := &dirStore{dir}
	if _, err := s.Load(s.Default()); err != nil {
		if err := s.Store(s.Default(), starterCode); err != nil {
			return nil, err
		}
	}
	return s, nil
}

type inMemCodeStore struct {
	codes map[Pid]string
}

func (i *inMemCodeStore) List() ([]Pid, error) {
	var r []Pid
	for c := range i.codes {
		r = append(r, c)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r, nil
}

func (i *inMemCodeStore) Search(query string) (res []Pid) {
	defer sort.Slice(res, func(i, j int) bool {
		return string(res[i]) < string(res[j])
	})
	for k := range i.codes {
		if strings.Contains(string(k), query) {
			res = append(res, k)
		}
		if len(res) >= 20 {
			return
		}
	}
	return
}

func (i *inMemCodeStore) Load(id Pid) (string, error) {
	if c, ok := i.codes[id]; ok {
		return c, nil
	}
	return "", fmt.Errorf("NOT FOUND")
}

func (i *inMemCodeStore) Default() Pid {
	return "starter"
}

func (i *inMemCodeStore) Store(id Pid, code string) error {
	i.codes[id] = code
	return nil
}

const starterCode = `function setup() {
  createCanvas(400, 400);
}

function draw() {
  background(220);
}
`

func newInMemCodeStore() codeRepo {
	return &inMemCodeStore{
		codes: map[Pid]string{
			"starter": starterCode,
		},
	}
}
