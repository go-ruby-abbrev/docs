// SPDX-License-Identifier: BSD-3-Clause
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-ruby-abbrev/abbrev"
)

// words is the fixed representative word list, identical to the Ruby driver:
// a few dozen words in several families that share prefixes, so the
// unambiguous-abbreviation table is non-trivial (ambiguous shared prefixes are
// dropped, unique ones kept).
var words = []string{
	"ruby", "rules", "ruler", "rubygems", "rubocop",
	"cat", "car", "card", "care", "cargo", "carpet",
	"go", "golang", "gopher", "gone", "good", "goose",
	"test", "testing", "tester", "tested", "testify",
	"program", "programmer", "programming", "progress",
	"data", "database", "datacenter", "date", "dating",
}

// canon serializes an abbrev table to a single deterministic line so the Go and
// Ruby drivers can be checked byte-identical before timing: keys sorted, each
// "k=v", joined by ";".
func canon(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = k + "=" + m[k]
	}
	return strings.Join(parts, ";")
}

func main() {
	// --dump: print the canonical abbrev table for the cross-runtime output
	// check, then exit without timing anything.
	if len(os.Args) > 1 && os.Args[1] == "--dump" {
		fmt.Println(canon(abbrev.Abbrev(words)))
		return
	}

	bench("abbrev", 2000, func() { sink = abbrev.Abbrev(words) })
}
