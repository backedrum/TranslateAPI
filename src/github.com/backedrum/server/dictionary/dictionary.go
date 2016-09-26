package dictionary

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Entry struct {
	Original   string `xml:"form>orth"`
	Translated []string `xml:"sense>cit>quote"`
}

var dictMap = make(map[string][]string)
var reg, error = regexp.Compile("[^a-zA-Z\\d\\s:]+")

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

	for i := range words {

		if val, ok := dictMap[strings.ToLower(words[i])]; ok {
			buf.WriteString(translationWords(val, maxAlt) + " ")
		} else {
			buf.WriteString(words[i] + " ")
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