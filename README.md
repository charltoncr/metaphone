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

Ron Charlton
