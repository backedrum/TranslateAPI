package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"math"
	"os"
	"path/filepath"
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

var translateFunc = TranslateDefault

func InitDictionary(langFrom, langTo, path string) {
	LangFrom = strings.ToUpper(langFrom)
	LangTo = strings.ToUpper(langTo)

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
	return LangFrom == strings.ToUpper(from) && LangTo == strings.ToUpper(to)
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

// Find a similar text and distance to it for a given text string.
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
