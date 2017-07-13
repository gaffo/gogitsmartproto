package gogitsmartproto

import (
	"io"
	"bufio"
	"strings"
)

const (
	pr_initial = 0
	pr_header = 2
	pr_refs = 3
	pr_finished = 4
)

func ParseRefs(rdr io.Reader) (refs Refs, err error) {
	refs.Refs = make([]Ref, 0, 1024)
	reader := bufio.NewScanner(rdr)
	state := pr_initial
	for reader.Scan() {
		line := reader.Text()
		switch state {
		case pr_initial:
			state = pr_header
			continue
		case pr_header:
			state = pr_refs
			continue
		case pr_refs:
			if line == "0000" {
				state = pr_finished
				break
			}
			trimmed := line[4:]
			parts := strings.Split(trimmed, " ")
			refs.Refs = append(refs.Refs, Ref{
				Checksum: parts[0],
				Ref: parts[1],
			})
		}
	}
	return refs, nil
}