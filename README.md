# Double Metaphone

DoubleMetaphone returns two codes for a word.  Two words sounds similar when
either non-empty code of one word matches either of the other word's codes.
Producing codes for the words in a
word list and matching them with the codes for an attempted spelling
allows a spell checker program to suggest correct spellings.

Example use:

```go
import "github.com/charltoncr/metaphone"
// ...
m, m2 := metaphone.DoubleMetaphone("knewmoanya", 4)
n, n2 := metaphone.DoubleMetaphone("pneumonia", 4)
if m == n || m == n2 || m2 == n || len(m2) > 0 && m2 == n2 {
    // match
}
// m is "NMN", as is n, so the two spellings match.
// The maximum allowed length for each of m, m2,
// n and n2 is 4 in this case.
```

# Double Metaphone Convenience Functions

The Double Metaphone convenience functions ease the use of DoubleMetaphone.
Two function calls are sufficient to read all words in a file, create a
map of words that have the same metaphone return values, and find all words
in the map that match a given word/misspelling.  See the example below.

- func NewMetaphMap(wordlist []string, maxLen int) *MetaphMap
- func NewMetaphMapFromFile(fileName string, maxLen int) (*MetaphMap, error)
- func (metaph *MetaphMap) MatchWord(word string) (output []string)
- func (metaph *MetaphMap) Len() int

**NewMetaphMap** returns a MetaphMap made from a wordlist and a maximum length
for the DoubleMetaphone return values.

**NewMetaphMapFromFile** returns a MetaphMap made from a word list file and
a maximum length for the DoubleMetaphone return values.

**MatchWord** returns all words in metaph that sound like word. Case in word
is ignored.

**Len** returns the number of sounds-alike keys in the metaph map.

Example use:

```go
package main

import (
    "fmt"
    "github.com/charltoncr/metaphone"
)
func main() {
    // The file specified by fileName should contain a comprehesive word
    // list with one word per line.  (Error check is omitted for brevity.)
    fileName := "spellCheckerWords.txt" // (could be a *.txt.gz file)
    metaphMap, _ := metaphone.NewMetaphMapFromFile(fileName, 4)
    matches := metaphMap.MatchWord("knewmoanya")
    for _, word := range matches {
        fmt.Println(word)
    }
}
```

Ron Charlton
