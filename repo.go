package main

import (
	"fmt"
	"sort"
	"strings"
)

type codeRepo interface {
	Default() Pid
	Search(query string) []Pid
	Store(id Pid, code string) error
	Load(pid Pid) (string, error)
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
