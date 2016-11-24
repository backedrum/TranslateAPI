package main

import (
	"encoding/xml"
	"github.com/backedrum/server/dictionary"
	"net/http"
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

	maxAlt := 0
	maxAltStr := req.URL.Query().Get("max-alt")
	if maxAltStr != "" {
		maxAlt, _ = strconv.Atoi(maxAltStr)
	}

	response := ServerResponse{text, dictionary.TranslateText(text, maxAlt), from, to}

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
	// TODO get rid of hardcoded path
	dictionary.InitDictionary("./src/github.com/backedrum/server/dictionary-files/nld-eng.tei")
	http.HandleFunc("/translate", translate)

	http.ListenAndServe(":9000", nil)
}
