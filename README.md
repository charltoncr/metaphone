# Double Metaphone

DoubleMetaphone returns two codes for a word.  Words with
matching codes sound similar.  Producing codes for the words in a word list
and matching them with the codes for an attempted spelling allows a spell
checker program to suggest possible correct spellings.

Example use:

```go
import "github.com/charltoncr/metaphone"
// ...
m, m2 := metaphone.DoubleMetaphone("knewmoanya", 4)
n, n2 := metaphone.DoubleMetaphone("pneumonia", 4)
if m == n || m == n2 || m2 == n || len(m2) > 0 && m2 == n2 {
    // match
}
// m is "NMN", as is n, so the two spelling match.
// The maximum allowed length for each of m, m2, n and n2
// is 4 in this case.
```

# Double Metaphone Convenience Functions

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

**Len** returns the number of sound-alike entries in the metaph sounds map.

Example use:

```go
import "fmt"
import "github.com/charltoncr/metaphone"
// ...
// File wordlistFileName should contain a comprehesive word list
// with one word per line.
metaphMap := metaphone.NewMetaphMapFromFile(wordlistFileName, 4)
matches := metaphMap.MatchWord("knewmoanya")
for _, word = range matches {
    fmt.Println(word)
}
```

Ron Charlton
