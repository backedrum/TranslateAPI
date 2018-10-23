package main

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
)

type ServerResponse struct {
	Original        string
	Translated      string
	OriginalLang    string
	TranslateToLang string
}

func translate(res http.ResponseWriter, req *http.Request) {
	text := req.URL.Query().Get("text")
	from := req.URL.Query().Get("from")
	to := req.URL.Query().Get("to")

	if !IsSupported(from, to) {
		http.Error(res,
			"Sorry, server currently doesn't support translation from "+from+" to "+to,
			http.StatusNotImplemented)
		return
	}

	maxAlt := 0
	maxAltStr := req.URL.Query().Get("maxAlt")
	if maxAltStr != "" {
		maxAlt, _ = strconv.Atoi(maxAltStr)
	}

	response := ServerResponse{html.UnescapeString(text), html.UnescapeString(translateFunc(from, to, html.UnescapeString(text), maxAlt)), from, to}

	xml, error := xml.MarshalIndent(response, "", "  ")
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set(
		"Content-Type",
		"application/xml",
	)
	res.Write(xml)
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: server <language from> <language to> <mode> <path to file>")
		fmt.Println("Example: server NL EN prose my_dictionary.tld")
		os.Exit(1)
	}

	LangFrom = os.Args[1]
	LangTo = os.Args[2]
	Mode = os.Args[3]

	if "prose" == Mode {
		translateFunc = TranslateTextWithParse
	} else if "default" != Mode {
		fmt.Println("Translation mode should be set either to prose or to default")
		os.Exit(1)
	}

	InitDictionary(LangFrom, LangTo, os.Args[4])

	http.HandleFunc("/translate", translate)

	http.ListenAndServe(":9000", nil)
}
