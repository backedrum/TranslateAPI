package main

import (
	"gopkg.in/jdkato/prose.v2"
	"testing"
)

func TestTranslateTextWithParse(t *testing.T) {
	golds := []struct {
		text           string
		translatedText string
	}{
		{"Homer Simpson.", "Homer Simpson."},
		{"test1 27", "a1[alt:a2] 27"},
	}

	for i, gold := range golds {
		translatedText := TranslateTextWithParse("NL", "EN", gold.text, 3)
		if gold.translatedText != translatedText {
			t.Fatalf("Failed case %d. Expected translation '%s' but actual was '%s'", i, gold.translatedText, translatedText)
		}
	}

}

func TestSplitToSequences(t *testing.T) {
	golds := []struct {
		tokens            []prose.Token
		expectedSequences []string
	}{
		0: {[]prose.Token{{Tag: "(", Text: "("}, {Tag: "JJ", Text: "Something"}, {Tag: ")", Text: ")"}},
			[]string{"(", "Something", ")"}},
		1: {[]prose.Token{{Tag: "PERSON", Text: "Homer Simpson"}, {Tag: ".", Text: "."}},
			[]string{"Homer Simpson", "."}},
		2: {[]prose.Token{{Tag: "NN", Text: "http://google.com"}, {Tag: "SYM", Text: ":D"}},
			[]string{"http://google.com", ":D"}},
		3: {[]prose.Token{{Tag: "JJ", Text: "test1"}, {Tag: "SYM", Text: ":)"}},
			[]string{"test1", ":)"}},
	}

	for i, gold := range golds {
		for j, seq := range splitToSequences(gold.tokens) {
			if seq != gold.expectedSequences[j] {
				t.Fatalf("Failed case %d. Expected sequence:'%s' but splitted was:'%s'", i, gold.expectedSequences[j], seq)
			}
		}
	}
}
