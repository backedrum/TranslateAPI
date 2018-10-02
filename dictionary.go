package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/jdkato/prose.v2"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TEI file mappings
type Entry struct {
	Original   string   `xml:"form>orth"`
	Translated []string `xml:"sense>cit>quote"`
}

var LangFrom string
var LangTo string
var Mode string

var dictMap = make(map[string][]string)

var translateFunc = TranslateText

func InitDictionary(langFrom, langTo, path string) {
	LangFrom = langFrom
	LangTo = langTo

	filePath, error := filepath.Abs(path)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	xmlFile, error := os.Open(filePath)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "entry" {
				var e Entry
				decoder.DecodeElement(&e, &startElement)
				// put pair into the map
				dictMap[strings.ToLower(e.Original)] = e.Translated
			}
		}
	}
}

func IsSupported(from, to string) bool {
	return LangFrom == from && LangTo == to
}

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

	resetSeq := func() {
		seqSoFar = ""
		distSoFar = math.MaxInt64
	}

	setSeq := func(newDictSeq string, newDistSoFar int) {
		seqSoFar = newDictSeq
		distSoFar = newDistSoFar
	}

	for _, tok := range tokens {
		nonPunctuation := regexp.MustCompile("^[A-Z]+$")
		if !nonPunctuation.MatchString(tok.Tag) {
			// non empty seqSoFar?
			if seqSoFar != "" {
				result = append(result, seqSoFar)
				resetSeq()
			}

			result = append(result, tok.Text+" ")
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
				seqSoFar = dictSeq
				distSoFar = dist
				continue
			}

			// flush previous sequence
			result = append(result, seqSoFar)

			// start with building of a new sequence
			setSeq(findByMinDist(strings.ToLower(tok.Text)))
		}

	}

	// final flush
	if seqSoFar != "" {
		result = append(result, seqSoFar)
	}

	return result
}

// TODO take punctuation into account
func TranslateText(langFrom, langTo, text string, maxAlt int) string {
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

func translationWords(val []string, maxAlt int) string {
	var buf bytes.Buffer

	for i := range val {
		if i > maxAlt {
			break
		}

		if i == 0 {
			buf.WriteString(val[i])
		} else {
			buf.WriteString("[alt:" + val[i] + "]")
		}
	}

	return buf.String()
}

func findByMinDist(word string) (string, int) {
	minDist := math.MaxInt64
	result := ""

	for k := range dictMap {
		// limit key set in order to limit number of distance calculations
		if k[0] != word[0] {
			continue
		}

		options := levenshtein.Options{
			InsCost: 1,
			DelCost: 1,
			SubCost: 1,
			Matches: func(sourceCharacter rune, targetCharacter rune) bool {
				return sourceCharacter == targetCharacter
			},
		}
		dist := levenshtein.DistanceForStrings([]rune(word), []rune(k), options)

		if dist < minDist {
			minDist = dist
			result = k
		}

	}

	return result, minDist
}
