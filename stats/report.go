package stats

import (
	"fmt"
	"strings"
)

type report struct {
	// is the duplicate counter tracking the number of unique lines that have
	// been duplicated, or the number of times all given lines are duplicated?
	// I don't know. So, we're going to say it's the latter
	dCount int64
	// floating point arithmetic is obviously not decimal arithmetic.
	// this is the slight inaccuracy inherent in the system right now.
	// working with Go's decimal type would be far more of a pain.
	medLength float64
	devLength float64
	medTokens float64
	devTokens float64
	keywords  map[string]int64
}

func (r *report) kwMerge(kw map[string]int) {
	if r.keywords == nil {
		r.keywords = make(map[string]int64)
	}
	for key, count := range kw {
		if current, exists := r.keywords[key]; exists {
			r.keywords[key] = current + int64(count)
		} else {
			r.keywords[key] = int64(count)
		}
	}
}

func (r report) String() string {
	lines := make([]string, 5)
	lines[0] = fmt.Sprintf("num dupes\t%d", r.dCount)
	lines[1] = fmt.Sprintf("med length\t%f", r.medLength)
	lines[2] = fmt.Sprintf("std length\t%f", r.devLength)
	lines[3] = fmt.Sprintf("med tokens\t%f", r.medTokens)
	lines[4] = fmt.Sprintf("std tokens\t%f", r.devTokens)
	for word, count := range r.keywords {
		lines = append(lines,
			fmt.Sprintf("%s\t%d", word, count),
		)
	}
	return strings.Join(lines, "\n")
}
