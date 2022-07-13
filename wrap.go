package gg

import (
	"strings"
	"unicode"
)

type measureStringer interface {
	MeasureString(s string) (w, h float64)
}

type stringsWithSpace []stringWithSpace

func newStringsWithSpace() *stringsWithSpace {
	return &stringsWithSpace{}
}

func (ss *stringsWithSpace) add(s string, isSpace bool) {
	if !isSpace {
		*ss = append(*ss, stringWithSpace{s: s})
	} else {
		if len(*ss) == 0 {
			*ss = append(*ss, stringWithSpace{spaces: s})
		} else {
			(*ss)[len(*ss)-1].spaces = s
		}
	}
}

type stringWithSpace struct {
	s      string
	spaces string
}

func splitOnSpace(x string) stringsWithSpace {
	ss := newStringsWithSpace()
	pi := 0
	ps := false
	for i, c := range x {
		s := unicode.IsSpace(c)
		if s != ps && i > 0 {
			ss.add(x[pi:i], ps)
			pi = i
		}
		ps = s
	}
	ss.add(x[pi:], ps)
	return *ss
}

func wordWrap(m measureStringer, s string, width float64) []string {
	// Loop through lines
	var result []string
	for _, line := range strings.Split(s, "\n") {
		// We want to keep empty lines
		if len(line) == 0 {
			result = append(result, "")
		}

		// Split on space
		ss := splitOnSpace(line)

		x := ""
		for _, s := range ss {
			w, _ := m.MeasureString(x + s.s)
			if w > width {
				if x == "" {
					result = append(result, s.s)
					x = ""
					continue
				} else {
					result = append(result, x)
					x = ""
				}
			}
			x += s.s + s.spaces
		}
		if x != "" {
			result = append(result, x)
		}
	}
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}
	return result
}
