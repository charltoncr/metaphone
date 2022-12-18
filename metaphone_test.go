// metaphone_Test.go - test util/metaphone.go by comparing lines from
// testInputData.txt.gz as input.  Comparison is with testInputData.txt.gz
// that contains output from a test program using dmetaph.cpp code.
// testInputData.txt.gz contains 171,110 words from an Aspell wordlist.
// testWantData.txt.gz contains the same number of lines output by a
// test program and the original C++ source code that metaphone.go was ported
// from.
// Author: Ron Charlton
// Date:   2022-12-12
// This file is public domain.  Public domain is per CC0 1.0; see
// https://creativecommons.org/publicdomain/zero/1.0/ for information.

package metaphone

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

var rcs_id_mt string = "$Id: metaphone_test.go,v 1.11 2022-12-17 23:29:19-05 ron Exp $"

func TestMetaphone(t *testing.T) {
	var words, want []string
	var err error

	if words, err = readFileLines("testInputData.txt.gz"); err == nil {
		want, err = readFileLines("testWantData.txt.gz")
	}
	if err != nil {
		t.Fatalf("%v", err)
	}

	idx := 0
	for _, word := range words {
		if len(word) > 0 {
			m, m2 := DoubleMetaphone(word, 6)
			got := fmt.Sprintf("'%s' '%s' %s", m, m2, word)
			if got != want[idx] {
				t.Errorf("At line %d got: <%s>;  want: <%s>", idx, got, want[idx])
			}
		}
		idx++
	}
}

// readFileLines reads a (gzipped) text file and returns its lines.
// Parameter name is the file's name.
func readFileLines(name string) (lines []string, err error) {
	var b []byte
	var r io.Reader
	var fp *os.File

	if fp, err = os.Open(name); err != nil {
		err = fmt.Errorf("trying to open file %s: %v", name, err)
		return
	}
	defer fp.Close()
	r = fp
	if strings.HasSuffix(name, ".gz") {
		if r, err = gzip.NewReader(r); err != nil {
			err = fmt.Errorf(
				"trying to make a gzip reader for file %s: %v", name, err)
			return
		}
	}
	if b, err = io.ReadAll(r); err != nil {
		err = fmt.Errorf("trying to read word list file %s: %v", name, err)
		return
	}
	lines = strings.Split(string(b), "\n")
	return
}

func TestConvenience(t *testing.T) {
	metaph, err := NewMetaphMapFromFile("testInputData.txt.gz", 6)
	if err != nil {
		t.Fatalf("%v", err)
	}
	words := metaph.MatchWord("knewmoania")
	if len(words) != 11 {
		t.Errorf("got: %d;  want: 11", len(words))
	}
}
