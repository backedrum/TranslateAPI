package main

import (
	"bytes"
	"strings"
)

// TODO take punctuation into account
func TranslateDefault(langFrom, langTo, text string, maxAlt int) string {
	var res bytes.Buffer

	// split text into words
	words := strings.Fields(text)

	for i := 0; i < len(words); i++ {
		word := strings.ToLower(words[i])
		foundWord, distWord := findByMinDist(word)

		if i+1 < len(words) {
			nextWord := strings.ToLower(words[i+1])
			foundPair, distPair := findByMinDist(word + " " + nextWord)

			// pair has better result than a single word?
			if distPair < distWord {
				res.WriteString(translationWords(dictMap[foundPair], maxAlt) + " ")
				i += 2
				continue
			}
		}

		if foundWord == "" {
			res.WriteString(word + " ")
		} else {
			res.WriteString(translationWords(dictMap[foundWord], maxAlt) + " ")
		}

	}

	return res.String()
}

