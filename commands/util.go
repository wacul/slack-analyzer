package commands

import (
	"fmt"
	"io"
	"time"
)

type StringFilter map[string]struct{}

func NewStringFilter(targets []string) StringFilter {
	m := map[string]struct{}{}
	for _, c := range targets {
		m[c] = struct{}{}
	}
	return StringFilter(m)
}

func (f StringFilter) Match(s string) bool {
	if len(f) > 0 {
		_, ok := f[s]
		return ok
	} else {
		return true
	}
}

type StringCounter map[string]int

func (s *StringCounter) Add(key string, count int) {
	if s == nil || *s == nil {
		*s = StringCounter{}
	}
	map[string]int(*s)[key] += count
}

func (s *StringCounter) Fprint(w io.Writer) {
	if s == nil {
		return
	}
	var index int
	for key, count := range *s {
		fmt.Fprintf(w, "%s : %d", key, count)
		if (index-3)%4 == 0 {
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, "  ")
		}
		index++
	}
}

type DateTime struct {
	raw    string
	parsed time.Time
}

const (
	DateFormat = "2006-01-02T15:04:05"
)

func (d *DateTime) Set(v string) error {
	p, err := time.Parse(DateFormat, v)
	if err != nil {
		return err
	}
	d.raw = v
	d.parsed = p
	return nil
}

func (d DateTime) String() string {
	return fmt.Sprintf("%s (%s)", d.raw, d.parsed.Format(time.RFC3339))
}

func (d DateTime) Time() time.Time {
	return d.parsed
}
