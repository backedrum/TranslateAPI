package dictionary

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Entry struct {
	Original   string   `xml:"form>orth"`
	Translated []string `xml:"sense>cit>quote"`
}

var dictMap = make(map[string][]string)
var reg, _ = regexp.Compile("[^a-zA-Z\\d\\s:]+")

func InitDictionary(path string) {
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

func TranslateTextAsWordsList(text string, maxAlt int) string {
	words := strings.Split(reg.ReplaceAllString(text, ""), " ")

	var buf bytes.Buffer

	for i := 0; i < len(words); i++ {

		if val, ok := dictMap[strings.ToLower(words[i])]; ok {
			buf.WriteString(translationWords(val, maxAlt) + " ")
		} else {
			// try to find two words match
			if i+1 < len(words) {
				if found, ok := dictMap[strings.ToLower(words[i]+" "+words[i+1])]; ok {
					buf.WriteString(translationWords(found, maxAlt) + " ")
					i++
					continue
				}
			}

			// find by min distance
			found := findByMinDist(strings.ToLower(words[i]))
			if found == "" {
				buf.WriteString(words[i] + " ")
			} else {
				buf.WriteString(translationWords(dictMap[found], maxAlt))
			}

		}
	}

	return buf.String()
}

func translationWords(val []string, maxAlt int) string {
	var buf bytes.Buffer

	for i, _ := range val {
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

func findByMinDist(word string) string {
	minDist := -1
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

		if minDist == -1 || dist < minDist {
			minDist = dist
			result = k
		}

	}

	return result
}
