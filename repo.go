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
	// Search(query string) []Pid
	Store(id Pid, code string) error
	Load(pid Pid) (string, error)
}

type dirStore struct {
	dir string
}

func (i *dirStore) path(id Pid) string {
	return filepath.Join(i.dir, string(id)+".js")
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
