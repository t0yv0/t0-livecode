package main

import (
	"fmt"
	"regexp"
)

type Pid string

func ParsePid(raw string) (Pid, error) {
	pattern := `^[a-zA-Z][-_a-zA-Z0-9]*$`
	var re = regexp.MustCompile(pattern)
	if !re.MatchString(raw) {
		return "", fmt.Errorf("Invalid program ID: %q, must match %q", raw, pattern)
	}
	return Pid(raw), nil
}
