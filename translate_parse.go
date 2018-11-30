package main

import (
	"bytes"
	"gopkg.in/jdkato/prose.v2"
	"log"
	"math"
	"regexp"
	"strings"
)

func TranslateTextWithParse(langFrom, langTo, text string, maxAlt int) string {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatal(err)
		return "Sorry, cannot build a new document from text " + text
	}

	var res bytes.Buffer

	for _, seq := range splitToSequences(doc.Tokens()) {
		if dictMap[seq] != nil {
			res.WriteString(translationWords(dictMap[seq], maxAlt) + " ")
			continue
		}

		res.WriteString(seq)
	}

	return res.String()
}

func splitToSequences(tokens []prose.Token) []string {
	var result []string

	seqSoFar := ""
	distSoFar := math.MaxInt64

	// start building a next sequence
	resetSeq := func() {
		seqSoFar = ""
		distSoFar = math.MaxInt64
	}

	// set current sequence
	setSeq := func(newDictSeq string, newDistSoFar int) {
		seqSoFar = newDictSeq
		distSoFar = newDistSoFar
	}

	for _, tok := range tokens {
		// we use this regexp in order to determine whether
		// a certain token should be translated
		shouldBeTranslated := regexp.MustCompile("^[A-Z]+$")
		if !shouldBeTranslated.MatchString(tok.Tag) {
			// non empty seqSoFar?
			if seqSoFar != "" {
				result = append(result, seqSoFar)
				resetSeq()
			}

			result = append(result, tok.Text)
			continue
		}

		// sequence extension and flush related logic
		newSeq := strings.ToLower(tok.Text)
		if seqSoFar != "" {
			newSeq = seqSoFar + " " + newSeq
		}

		dictSeq, dist := findByMinDist(newSeq)

		if dictSeq != "" {

			// better sequence? increase
			if dist < distSoFar {
				setSeq(dictSeq, dist)
				continue
			}
		}

		// does not make sense to continue, so flushing a previous sequence and start with a new one
		if seqSoFar != "" {
			result = append(result, seqSoFar)
		}

		// deal with the current token
		dictSeq, _ = findByMinDist(strings.ToLower(tok.Text))
		if dictSeq == "" {
			result = append(result, tok.Text + " ")
			resetSeq()
		} else {
			setSeq(findByMinDist(strings.ToLower(tok.Text)))
		}
	}

	// final flush
	if seqSoFar != "" {
		result = append(result, seqSoFar)
	}

	return result
}
