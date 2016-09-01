// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"fmt"
	"math"
	"regexp"
)

type Task struct {
	Level       uint
	Match       []string
	Dmatch      []string
}

func (t Task) matches(sol string) bool {
	re, err := regexp.Compile(sol)
	if err != nil {
		return false
	}

	for _, w := range t.Match {
		if re.FindString(w) != w {
			return false
		}
	}

	for _, w := range t.Dmatch {
		if re.FindString(w) == w {
			return false
		}
	}

	return true
}

// calculates a level for the task, based on the
// the words it has and shouldn't match
func (t *Task) calcLevel() {
	var (
		c float64               // count
		r = make(map[rune]uint) // character occurrence map
		a float64               // average rune count
	)

	for _, w := range t.Match {
		c++
		for _, b := range w {
			r[b]++
		}
	}

	for _, w := range t.Dmatch {
		c++
		for _, b := range w {
			r[b]++
		}
	}

	for _, n := range r {
		a += float64(n)
	}
	a /= float64(len(r))

	t.Level = uint(math.Ceil(math.Abs(math.Log10(math.Pow(a, c)))))
}

// Check if task is "acceptable"
// currently only checks if a word in M is in D too
func (t Task) isAcceptable() error {
	if (len(t.Match) == 0 || len(t.Dmatch) == 0) {
		return fmt.Errorf("Either words to match (%d) or not to match (%d) are empty", len(t.Match), len(t.Dmatch))
	}

	for _, w := range t.Match {
		for _, c := range t.Dmatch {
			if w == c && w != "" {
				return fmt.Errorf("Words to match and words not to match have duplicates (%q and %q)", w, c)
			}
		}
	}

	return nil
}

// writes a solution for a task into the database
func (_ *Task) submit(sol string) {
	fmt.Printf("Suggested \"%s\"\n", sol)
}
