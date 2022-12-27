// Convenience functions that use DoubleMetaphone.
// Created 2022-12-16 by Ron Charlton and placed in the public domain.
//
// $Id: convenience.go,v 1.20 2022-12-27 09:19:23-05 ron Exp $

package metaphone

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// MetaphMap defines a MetaphMap for a wordlist and maximum metaph/metaph2
// length from DoubleMetaphone.
type MetaphMap struct {
	mapper map[string][]string
	// maximum length of metaph and metaph2 in DoubleMetaphone.
	maxlen int
}

// NewMetaphMap returns a MetaphMap made from wordlist and a maximum
// length for the DoubleMetaphone return values.
// The MetaphMap can be used with MatchWord to find all words in the
// MetaphMap that sound like a given word or misspelling.
// Argument maxLen is 4 in the original Double Metaphone algorithm.
// Case is ignored in the words in wordlist, as are non-alphabetic
// characters.
func NewMetaphMap(wordlist []string, maxLen int) *MetaphMap {
	MMap := make(map[string][]string)
	for _, word := range wordlist {
		m, m2 := DoubleMetaphone(word, maxLen)
		if len(m) > 0 {
			MMap[m] = append(MMap[m], word)
		}
		if len(m2) > 0 {
			MMap[m2] = append(MMap[m2], word)
		}
	}
	return &MetaphMap{
		mapper: MMap,
		maxlen: maxLen,
	}
}

// NewMetaphMapFromFile returns a MetaphMap made from a file containing a
// word list, and using a maximum length for the DoubleMetaphone return values.
// The file can be a gzipped file with its name ending with ".gz".
// The MetaphMap can be used with MatchWord to find all words in the
// MetaphMap that sound like a given word or misspelling.
// Argument maxLen is 4 in the original Double Metaphone algorithm.
// Case and non-alphabetic characters in the file are ignored.
func NewMetaphMapFromFile(fileName string, maxLen int) (
	metaph *MetaphMap, err error) {
	var b []byte
	var r io.Reader
	var fp *os.File

	if fp, err = os.Open(fileName); err != nil {
		err = fmt.Errorf("trying to open file %s: %v", fileName, err)
		return
	}
	defer fp.Close()
	r = fp
	if strings.HasSuffix(fileName, ".gz") {
		if r, err = gzip.NewReader(r); err != nil {
			err = fmt.Errorf(
				"trying to make a gzip reader for file %s: %v", fileName, err)
			return
		}
	}
	if b, err = io.ReadAll(r); err != nil {
		err = fmt.Errorf("trying to read file %s: %v", fileName, err)
		return
	}
	lines := strings.Split(string(b), "\n")
	return NewMetaphMap(lines, maxLen), err
}

// Len returns the number of sound-alike entries in metaph.
func (metaph *MetaphMap) Len() int {
	return len(metaph.mapper)
}

// MatchWord returns all words in metaph that sound like word.
// Case and non-alphabetic characters in word are ignored.  Typical use:
//
//		import "fmt"
//		import "metaphone"
//		// ...
//		// File wordlistFileName should contain a comprehesive word
//	 	// list, one word per line.  Errors are ignored here.
//		metaphMap, _ := metaphone.NewMetaphMapFromFile(wordlistFileName, 4)
//		matches := metaphMap.MatchWord("knewmoanya")
//		for _, word = range matches {
//			fmt.Println(word)
//		}
func (metaph *MetaphMap) MatchWord(word string) (output []string) {
	m, m2 := DoubleMetaphone(word, metaph.maxlen)
	if len(m) > 0 {
		output = metaph.mapper[m]
	}
	if len(m2) > 0 {
		output = append(output, metaph.mapper[m2]...)
	}
	output = removeDups(output)
	return
}

// removeDups removes duplicates within s.
func removeDups(s []string) (out []string) {
	m := make(map[string]struct{})
	for _, w := range s {
		m[w] = struct{}{}
	}
	for o := range m {
		out = append(out, o)
	}
	return
}
