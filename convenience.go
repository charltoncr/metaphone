// Convenience functions that use DoubleMetaphone.
// Created 2022-12-16 by Ron Charlton and placed in the public domain.
//
// $Id: convenience.go,v 1.13 2022-12-17 13:43:38-05 ron Exp $

package metaphone

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
// Case is ignored in wordlist, as are non-alphabetic characters.
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

// Len returns the number of sound-alike entries in metaph.
func (metaph *MetaphMap) Len() int {
	return len(metaph.mapper)
}

// MatchWord returns all words in metaph that sound like word.
// Case and non-alphabetic characters in word are ignored.  Typical use:
//
//	import "fmt"
//	import "metaphone"
//	// ...
//	// wordlist should contain all words in a comprehesive word list.
//	metaphMap := metaphone.NewMetaphMap(wordlist, 4)
//	matches := metaphMap.MatchWord("knewmoanya")
//	for _, word = range matches {
//		fmt.Println(word)
//	}
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

// removeDups removes duplicates within in.
func removeDups(in []string) (out []string) {
	m := make(map[string]struct{})
	for _, w := range in {
		m[w] = struct{}{}
	}
	for o := range m {
		out = append(out, o)
	}
	return
}
