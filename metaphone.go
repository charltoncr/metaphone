// metaphone.go is an open source implementation of Double Metaphone.
//
// Ported by Ron Charlton from http://aspell.net/metaphone/dmetaph.cpp on
// 2022-12-04. dmetaph.cpp is open source, as is this Go port.
//
// See the original Dr.Dobb's article by Lawrence Philips at
// https://drdobbs.com/the-double-metaphone-search-algorithm/184401251?pgno=2
//
// $Id: metaphone.go,v 3.22 2022-12-15 09:34:52-05 ron Exp $

///
// From metaph.cpp:
//
// Double Metaphone (c) 1998, 1999 by Lawrence Philips
//
// Slightly modified by Kevin Atkinson to fix several bugs and
// to allow it to give back more than 4 characters.
///

// Package metaphone is an open source implementation of Double Metaphone.
package metaphone

import "strings"

// DoubleMetaphone returns primary and secondary codes for word.
// Metaph and metaph2 are each limited to maxlength characters.
// The original Double Metaphone code set maxlength to 4.
// Non-alphabetic characters in word are ignored.  Upper/Lower case distinctions
// in word are also ignored.
// Metaph will be the same for similar sounding words/misspellings, as will
// metaph2.  If either metaph or metaph2 from one word matches either
// metaph or metaph2 from another word, and len(metaph2) > 0,
// then the two words sound similar.  Example code:
//
//	import "github.com/charltoncr/metaphone"
//	// ...
//	m, m2 := metaphone.DoubleMetaphone("knewmoanya", 6)
//	n, n2 := metaphone.DoubleMetaphone("pneumonia", 6)
//	// m is "NMN". n also is "NMN".  The maximum allowed length for each
//	// of m, m2, n and n2 is 6.
//	if m == n || m == n2 || m2 == n || len(m2) > 0 && m2 == n2 {
//	  // match
//	}
//	// ...
func DoubleMetaphone(word string, maxlength int) (metaph, metaph2 string) {
	const pad = "     " // 5 spaces

	length := len(word)
	if length < 1 {
		return
	}
	if maxlength < 1 {
		maxlength = 4
	}

	var current = 0
	last := length - 1                     //zero based index
	var primary, secondary strings.Builder // becomes metaph, metaph2

	word = strings.ToUpper(word)

	// pad with spaces at end
	word += pad
	rword := []rune(word)
	rwordLen := len(rword)
	alternate := false

	// Found determines whether s is in word.
	Found := func(s string) bool {
		return strings.Contains(word, s)
	}

	SlavoGermanic := func() bool {
		return Found("W") || Found("K") || Found("CZ") // never reached: || Found("WITZ")
	}

	// MetaphAdd adds a string to primary and secondary.  Call it with 1 or 2
	// arguments.  The first argument is appended to primary (and to
	// secondary if a second argument is empty or not provided).  Any second,
	// non-empty argument is appended to secondary.
	MetaphAdd := func(s ...string) {
		if len(s) < 1 || len(s) > 2 {
			panic("MetaphAdd requires one or two arguments")
		}
		main := s[0]
		primary.WriteString(main)
		if len(s) == 1 {
			secondary.WriteString(main)
		} else {
			alt := s[1]
			if len(alt) > 0 {
				alternate = true
				if alt[0] != ' ' {
					secondary.WriteString(alt)
				}
			} else if len(main) > 0 && main[0] != ' ' {
				secondary.WriteString(main)
			}
		}
	}

	// GetAt returns the rune at index 'at' in rword, or rune(0) if 'at' is
	// out of range.
	GetAt := func(at int) rune {
		if at < 0 || at >= rwordLen {
			return 0
		}
		return rword[at]
	}

	// IsVowel returns true if the rune at index 'at' in rword is a vowel.
	IsVowel := func(at int) bool {
		if at < 0 || at >= rwordLen {
			return false
		}
		return strings.ContainsRune("AEIOUY", GetAt(at))
	}

	// StringAt determines if any of a list of string arguments appear
	// in rword at start and length long.
	StringAt := func(start, length int, s ...string) bool {
		if start < 0 || (start+length) >= rwordLen {
			return false
		}
		target := string(rword[start : start+length])
		for i := 0; i < len(s); i++ {
			if s[i] == target {
				return true
			}
		}
		return false
	}

	//skip these when at start of word
	if StringAt(0, 2, "GN", "KN", "PN", "WR", "PS") {
		current += 1
	}

	//Initial 'X' is pronounced 'Z' e.g. 'Xavier'
	if GetAt(0) == 'X' {
		MetaphAdd("S") //'Z' maps to 'S'
		current += 1
	}

	///////////main loop//////////////////////////
	for current < length &&
		(primary.Len() < maxlength || secondary.Len() < maxlength) {
		switch GetAt(current) {
		case 'A', 'E', 'I', 'O', 'U', 'Y':
			if current == 0 {
				//all init vowels now map to 'A'
				MetaphAdd("A")
			}
			current += 1
		case 'B':
			//"-mb", e.g", "dumb", already skipped over...
			MetaphAdd("P")

			if GetAt(current+1) == 'B' {
				current += 2
			} else {
				current += 1
			}
		case 'Ç':
			MetaphAdd("S")
			current += 1
		case 'C':
			//various germanic
			if current > 1 &&
				!IsVowel(current-2) &&
				StringAt((current-1), 3, "ACH") &&
				((GetAt(current+2) != 'I') && ((GetAt(current+2) != 'E') ||
					StringAt((current-2), 6, "BACHER", "MACHER"))) {
				MetaphAdd("K")
				current += 2
				break
			}

			//special case 'caesar'
			if current == 0 && StringAt(current, 6, "CAESAR") {
				MetaphAdd("S")
				current += 2
				break
			}

			//italian 'chianti'
			if StringAt(current, 4, "CHIA") {
				MetaphAdd("K")
				current += 2
				break
			}

			if StringAt(current, 2, "CH") {
				//find 'michael'
				if current > 0 && StringAt(current, 4, "CHAE") {
					MetaphAdd("K", "X")
					current += 2
					break
				}

				//greek roots e.g. 'chemistry', 'chorus'
				if (current == 0) &&
					(StringAt((current+1), 5, "HARAC", "HARIS") ||
						StringAt((current+1), 3, "HOR", "HYM", "HIA", "HEM")) &&
					!StringAt(0, 5, "CHORE") {
					MetaphAdd("K")
					current += 2
					break
				}

				//germanic, greek, or otherwise 'ch' for 'kh' sound
				if (StringAt(0, 4, "VAN ", "VON ") || StringAt(0, 3, "SCH")) ||
					// 'architect but not 'arch', 'orchestra', 'orchid'
					StringAt((current-2), 6, "ORCHES", "ARCHIT", "ORCHID") ||
					StringAt((current+2), 1, "T", "S") ||
					((StringAt((current-1), 1, "A", "O", "U", "E") || (current == 0)) &&
						//e.g., 'wachtler', 'wechsler', but not 'tichner'
						StringAt((current+2), 1, "L", "R", "N", "M", "B", "H", "F", "V", "W", " ")) {
					MetaphAdd("K")
				} else {
					if current > 0 {
						if StringAt(0, 2, "MC") {
							//e.g., "McHugh"
							MetaphAdd("K")
						} else {
							MetaphAdd("X", "K")
						}
					} else {
						MetaphAdd("X")
					}
				}
				current += 2
				break
			}
			//e.g, 'czerny'
			if StringAt(current, 2, "CZ") && !StringAt((current-2), 4, "WICZ") {
				MetaphAdd("S", "X")
				current += 2
				break
			}

			//e.g., 'focaccia'
			if StringAt((current + 1), 3, "CIA") {
				MetaphAdd("X")
				current += 3
				break
			}

			//double 'C', but not if e.g. 'McClellan'
			if StringAt(current, 2, "CC") && !((current == 1) && (GetAt(0) == 'M')) {
				//'bellocchio' but not 'bacchus'
				if StringAt((current+2), 1, "I", "E", "H") && !StringAt((current+2), 2, "HU") {
					//'accident', 'accede' 'succeed'
					if (current == 1 && (GetAt(current-1) == 'A')) ||
						StringAt((current-1), 5, "UCCEE", "UCCES") {
						MetaphAdd("KS")
						//'bacci', 'bertucci', other italian
					} else {
						MetaphAdd("X")
					}
					current += 3
					break
				} else { //Pierce's rule
					MetaphAdd("K")
					current += 2
					break
				}
			}
			if StringAt(current, 2, "CK", "CG", "CQ") {
				MetaphAdd("K")
				current += 2
				break
			}

			if StringAt(current, 2, "CI", "CE", "CY") {
				//italian vs. english
				if StringAt(current, 3, "CIO", "CIE", "CIA") {
					MetaphAdd("S", "X")
				} else {
					MetaphAdd("S")
				}
				current += 2
				break
			}

			//else
			MetaphAdd("K")

			//name sent in 'mac caffrey', 'mac gregor
			if StringAt((current + 1), 2, " C", " Q", " G") {
				current += 3
			} else {
				if StringAt((current+1), 1, "C", "K", "Q") &&
					!StringAt((current+1), 2, "CE", "CI") {
					current += 2
				} else {
					current += 1
				}
			}
		case 'D':
			if StringAt(current, 2, "DG") {
				if StringAt((current + 2), 1, "I", "E", "Y") {
					//e.g. 'edge'
					MetaphAdd("J")
					current += 3
				} else {
					//e.g. 'edgar'
					MetaphAdd("TK")
					current += 2
				}
				break
			}
			if StringAt(current, 2, "DT", "DD") {
				MetaphAdd("T")
				current += 2
				break
			}

			//else
			MetaphAdd("T")
			current += 1
		case 'F':
			if GetAt(current+1) == 'F' {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("F")
		case 'G':
			if GetAt(current+1) == 'H' {
				if current > 0 && !IsVowel(current-1) {
					MetaphAdd("K")
					current += 2
					break
				}

				if current < 3 {
					//'ghislane', ghiradelli
					if current == 0 {
						if GetAt(current+2) == 'I' {
							MetaphAdd("J")
						} else {
							MetaphAdd("K")
						}
						current += 2
						break
					}
				}
				//Parker's rule (with some further refinements) - e.g., 'hugh'
				if !((current > 1 && StringAt((current-2), 1, "B", "H", "D")) ||
					//e.g., 'bough'
					(current > 2 && StringAt((current-3), 1, "B", "H", "D")) ||
					//e.g., 'broughton'
					(current > 3 && StringAt((current-4), 1, "B", "H"))) {
					//e.g., 'laugh', 'McLaughlin', 'cough', 'gough', 'rough', 'tough'
					if current > 2 && GetAt(current-1) == 'U' &&
						StringAt((current-3), 1, "C", "G", "L", "R", "T") {
						MetaphAdd("F")
					} else {
						if current > 0 && GetAt(current-1) != 'I' {
							MetaphAdd("K")
						}
					}
				}
				current += 2
				break
			}

			if GetAt(current+1) == 'N' {
				if current == 1 && IsVowel(0) && !SlavoGermanic() {
					MetaphAdd("KN", "N")
				} else {
					//not e.g. 'cagney'
					if !StringAt((current+2), 2, "EY") &&
						(GetAt(current+1) != 'Y') && !SlavoGermanic() {
						MetaphAdd("N", "KN")
					} else {
						MetaphAdd("KN")
					}
				}
				current += 2
				break
			}

			//'tagliaro'
			if StringAt((current+1), 2, "LI") && !SlavoGermanic() {
				MetaphAdd("KL", "L")
				current += 2
				break
			}

			//-ges-,-gep-,-gel-, -gie- at beginning
			if current == 0 &&
				((GetAt(current+1) == 'Y') ||
					StringAt((current+1), 2, "ES", "EP", "EB", "EL", "EY", "IB", "IL", "IN", "IE", "EI", "ER")) {
				MetaphAdd("K", "J")
				current += 2
				break
			}

			// -ger-,  -gy-
			if (StringAt(current+1, 2, "ER") || (GetAt(current+1) == 'Y')) &&
				!StringAt(0, 6, "DANGER", "RANGER", "MANGER") &&
				!StringAt(current-1, 1, "E", "I") &&
				!StringAt(current-1, 3, "RGY", "OGY") {
				MetaphAdd("K", "J")
				current += 2
				break
			}

			// italian e.g, 'biaggi'
			if StringAt((current+1), 1, "E", "I", "Y") ||
				StringAt((current-1), 4, "AGGI", "OGGI") {
				//obvious germanic
				if StringAt(0, 4, "VAN ", "VON ") || StringAt(0, 3, "SCH") ||
					StringAt(current+1, 2, "ET") {
					MetaphAdd("K")
				} else {
					//always soft if french ending
					if StringAt((current + 1), 4, "IER ") {
						MetaphAdd("J")
					} else {
						MetaphAdd("J", "K")
					}
				}
				current += 2
				break
			}

			if GetAt(current+1) == 'G' {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("K")
		case 'H':
			//only keep if first & before vowel or btw. 2 vowels
			if ((current == 0) || IsVowel(current-1)) && IsVowel(current+1) {
				MetaphAdd("H")
				current += 2
			} else { //also takes care of 'HH'
				current += 1
			}
		case 'J':
			//obvious spanish, 'jose', 'san jacinto'
			if StringAt(current, 4, "JOSE") || StringAt(0, 4, "SAN ") {
				if ((current == 0) && (GetAt(current+4) == ' ')) || StringAt(0, 4, "SAN ") {
					MetaphAdd("H")
				} else {
					MetaphAdd("J", "H")
				}
				current += 1
				break
			}

			if current == 0 && !StringAt(current, 4, "JOSE") {
				MetaphAdd("J", "A") //Yankelovich/Jankelowicz
			} else {
				//spanish pron. of e.g. 'bajador'
				if IsVowel(current-1) &&
					!SlavoGermanic() &&
					((GetAt(current+1) == 'A') || (GetAt(current+1) == 'O')) {
					MetaphAdd("J", "H")
				} else {
					if current == last {
						MetaphAdd("J", " ")
					} else {
						if !StringAt((current+1), 1, "L", "T", "K", "S", "N", "M", "B", "Z") &&
							!StringAt((current-1), 1, "S", "K", "L") {
							MetaphAdd("J")
						}
					}
				}
			}
			if GetAt(current+1) == 'J' { //it could happen!
				current += 2
			} else {
				current += 1
			}
		case 'K':
			if GetAt(current+1) == 'K' {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("K")
		case 'L':
			if GetAt(current+1) == 'L' {
				//spanish e.g. 'cabrillo', 'gallegos'
				if ((current == (length - 3)) &&
					StringAt((current-1), 4, "ILLO", "ILLA", "ALLE")) ||
					((StringAt((last-1), 2, "AS", "OS") || StringAt(last, 1, "A", "O")) &&
						StringAt((current-1), 4, "ALLE")) {
					MetaphAdd("L", " ")
					current += 2
					break
				}
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("L")
		case 'M':
			if (StringAt((current-1), 3, "UMB") &&
				(((current + 1) == last) || StringAt((current+2), 2, "ER"))) ||
				//'dumb','thumb'
				(GetAt(current+1) == 'M') {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("M")
		case 'N':
			if GetAt(current+1) == 'N' {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("N")
		case 'Ñ':
			current += 1
			MetaphAdd("N")
		case 'P':
			if GetAt(current+1) == 'H' {
				MetaphAdd("F")
				current += 2
				break
			}
			//also account for "campbell", "raspberry"
			if StringAt((current + 1), 1, "P", "B") {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("P")
		case 'Q':
			if GetAt(current+1) == 'Q' {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("K")
		case 'R':
			//french e.g. 'rogier', but exclude 'hochmeier'
			if current == last && !SlavoGermanic() &&
				StringAt((current-2), 2, "IE") &&
				!StringAt((current-4), 2, "ME", "MA") {
				MetaphAdd("", "R")
			} else {
				MetaphAdd("R")
			}

			if GetAt(current+1) == 'R' {
				current += 2
			} else {
				current += 1
			}
		case 'S':
			//special cases 'island', 'isle', 'carlisle', 'carlysle'
			if StringAt((current - 1), 3, "ISL", "YSL") {
				current += 1
				break
			}

			//special case 'sugar-'
			if (current == 0) && StringAt(current, 5, "SUGAR") {
				MetaphAdd("X", "S")
				current += 1
				break
			}

			if StringAt(current, 2, "SH") {
				//germanic
				if StringAt((current + 1), 4, "HEIM", "HOEK", "HOLM", "HOLZ") {
					MetaphAdd("S")
				} else {
					MetaphAdd("X")
				}
				current += 2
				break
			}

			//italian & armenian
			if StringAt(current, 3, "SIO", "SIA") || StringAt(current, 4, "SIAN") {
				if !SlavoGermanic() {
					MetaphAdd("S", "X")
				} else {
					MetaphAdd("S")
				}
				current += 3
				break
			}

			//german & anglicisations, e.g. 'smith' match 'schmidt', 'snider' match 'schneider'
			//also, -sz- in slavic language altho in hungarian it is pronounced 's'
			if current == 0 &&
				StringAt((current+1), 1, "M", "N", "L", "W") ||
				StringAt((current+1), 1, "Z") {
				MetaphAdd("S", "X")
				if StringAt((current + 1), 1, "Z") {
					current += 2
				} else {
					current += 1
				}
				break
			}

			if StringAt(current, 2, "SC") {
				//Schlesinger's rule
				if GetAt(current+2) == 'H' {
					//dutch origin, e.g. 'school', 'schooner'
					if StringAt((current + 3), 2, "OO", "ER", "EN", "UY", "ED", "EM") {
						//'schermerhorn', 'schenker'
						if StringAt((current + 3), 2, "ER", "EN") {
							MetaphAdd("X", "SK")
						} else {
							MetaphAdd("SK")
						}
						current += 3
					} else {
						if current == 0 && !IsVowel(3) && (GetAt(3) != 'W') {
							MetaphAdd("X", "S")
						} else {
							MetaphAdd("X")
						}
						current += 3
					}
					break
				}

				if StringAt((current + 2), 1, "I", "E", "Y") {
					MetaphAdd("S")
					current += 3
					break
				}
				//else
				MetaphAdd("SK")
				current += 3
				break
			}

			//french e.g. 'resnais', 'artois'
			if current == last && StringAt((current-2), 2, "AI", "OI") {
				MetaphAdd("", "S")
			} else {
				MetaphAdd("S")
			}

			if StringAt((current + 1), 1, "S", "Z") {
				current += 2
			} else {
				current += 1
			}
		case 'T':
			if StringAt(current, 4, "TION") {
				MetaphAdd("X")
				current += 3
				break
			}

			if StringAt(current, 3, "TIA", "TCH") {
				MetaphAdd("X")
				current += 3
				break
			}

			if StringAt(current, 2, "TH") ||
				StringAt(current, 3, "TTH") {
				//special case 'thomas', 'thames' or germanic
				if StringAt((current+2), 2, "OM", "AM") ||
					StringAt(0, 4, "VAN ", "VON ") ||
					StringAt(0, 3, "SCH") {
					MetaphAdd("T")
				} else {
					MetaphAdd("0", "T")
				}
				current += 2
				break
			}

			if StringAt((current + 1), 1, "T", "D") {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("T")
		case 'V':
			if GetAt(current+1) == 'V' {
				current += 2
			} else {
				current += 1
			}
			MetaphAdd("F")
		case 'W':
			//can also be in middle of word
			if StringAt(current, 2, "WR") {
				MetaphAdd("R")
				current += 2
				break
			}

			if current == 0 &&
				(IsVowel(current+1) || StringAt(current, 2, "WH")) {
				//Wasserman should match Vasserman
				if IsVowel(current + 1) {
					MetaphAdd("A", "F")
				} else {
					//need Uomo to match Womo
					MetaphAdd("A")
				}
			}

			//Arnow should match Arnoff
			if (current == last && IsVowel(current-1)) ||
				StringAt((current-1), 5, "EWSKI", "EWSKY", "OWSKI", "OWSKY") ||
				StringAt(0, 3, "SCH") {
				MetaphAdd("", "F")
				current += 1
				break
			}

			//polish e.g. 'filipowicz'
			if StringAt(current, 4, "WICZ", "WITZ") {
				MetaphAdd("TS", "FX")
				current += 4
				break
			}

			//else skip it
			current += 1
		case 'X':
			//french e.g. breaux
			if !((current == last) &&
				(StringAt((current-3), 3, "IAU", "EAU") ||
					StringAt((current-2), 2, "AU", "OU"))) {
				MetaphAdd("KS")
			}

			if StringAt((current + 1), 1, "C", "X") {
				current += 2
			} else {
				current += 1
			}
		case 'Z':
			//chinese pinyin e.g. 'zhao'
			if GetAt(current+1) == 'H' {
				MetaphAdd("J")
				current += 2
				break
			} else {
				if StringAt((current+1), 2, "ZO", "ZI", "ZA") ||
					(SlavoGermanic() && ((current > 0) && GetAt(current-1) != 'T')) {
					MetaphAdd("S", "TS")
				} else {
					MetaphAdd("S")
				}
			}

			if GetAt(current+1) == 'Z' {
				current += 2
			} else {
				current += 1
			}
		default:
			current += 1
		}
	}

	metaph = primary.String()
	if len(metaph) > maxlength {
		metaph = metaph[:maxlength]
	}

	if alternate {
		metaph2 = secondary.String()
		if len(metaph2) > maxlength {
			metaph2 = metaph2[:maxlength]
		}
	}

	return
}
