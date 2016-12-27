package main

import (
	"encoding/xml"
	"fmt"
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
	maxAltStr := req.URL.Query().Get("max-alt")
	if maxAltStr != "" {
		maxAlt, _ = strconv.Atoi(maxAltStr)
	}

	response := ServerResponse{text, TranslateText(from, to, text, maxAlt), from, to}

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
	if len(os.Args) != 4 {
		fmt.Println("Usage: server <language from> <language to> <path to file>")
		os.Exit(1)
	}

	LangFrom = os.Args[1]
	LangTo = os.Args[2]

	InitDictionary(LangFrom, LangTo, os.Args[3])

	http.HandleFunc("/translate", translate)

	http.ListenAndServe(":9000", nil)
}
