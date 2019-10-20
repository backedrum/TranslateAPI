package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
)

const FROM_LANGUAGE_IDX = 1
const TO_LANGUAGE_IDX = 2
const TRANSLATION_MODE_IDX = 3
const DICTIONARY_PATH_IDX = 4

const DEFAULT_MODE = "default"
const PROSE_MODE = "prose"

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

	response := ServerResponse{html.UnescapeString(text),
		html.UnescapeString(translateFunc(from, to, html.UnescapeString(text), maxAlt)), from, to}

	json, error := json.MarshalIndent(response, "", "  ")
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(json)
}

func inspectEntry(res http.ResponseWriter, req *http.Request) {
	langFrom := req.URL.Query().Get("lang-from")
	langTo := req.URL.Query().Get("lang-to")
	textToInspect := req.URL.Query().Get("text")

	if !IsSupported(langFrom, langTo) {
		http.Error(res,
			"Sorry, server currently doesn't support translation from "+langFrom+" to "+langTo,
			http.StatusNotImplemented)
	}

	entry := Inspect(langFrom, langTo, textToInspect)

	json, error := json.MarshalIndent(*entry, "", "")
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(json)
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: server <language from> <language to> <mode> <path to file>")
		fmt.Println("Example: server NL EN prose my_dictionary.tld")
		os.Exit(1)
	}

	LangFrom = os.Args[FROM_LANGUAGE_IDX]
	LangTo = os.Args[TO_LANGUAGE_IDX]
	Mode = os.Args[TRANSLATION_MODE_IDX]

	if PROSE_MODE == Mode {
		translateFunc = TranslateTextWithParse
	} else if DEFAULT_MODE != Mode {
		fmt.Println("Translation mode should be set either to prose or to default")
		os.Exit(1)
	}

	InitDictionary(LangFrom, LangTo, os.Args[DICTIONARY_PATH_IDX])

	http.HandleFunc("/translate", translate)

	http.HandleFunc("/inspect", inspectEntry)

	http.ListenAndServe(":9000", nil)
}
