package main

import (
	"math"
	"strings"
	"testing"
)

func init() {
	// init dict
	dictMap["test1"] = []string{"a1", "a2"}
	dictMap["other"] = []string{"a3", "a4"}

	// init languages
	LangFrom = "NL"
	LangTo = "EN"
}

func TestFindByMinDist(t *testing.T) {
	golds := []struct {
		word    string
		expKey  string
		expDist int
	}{
		0: {"test1", "test1", 0},
		1: {"testN", "test1", 1},
		2: {"missed", "", math.MaxInt64},
		3: {"onther", "other", 1},
		4: {"oer", "other", 2},
	}

	for i, gold := range golds {
		key, dist := findByMinDist(gold.word)
		if gold.expKey != key || gold.expDist != dist {
			t.Fatalf("Failed case %d. Expected key and dist are:%s-%d, was:%s-%d", i, gold.expKey, gold.expDist, key, dist)
		}
	}
}

func TestTranslateText(t *testing.T) {
	golds := []struct {
		text           string
		alts           int
		expTranslation string
	}{
		0: {"test1 other", 0, "a1 a3"},
		1: {"test, not", 1, "a1[alt:a2] not"},
		2: {"other testX", 1, "a3[alt:a4] a1[alt:a2]"},
	}

	for i, gold := range golds {
		translated := strings.Trim(TranslateText("", "", gold.text, gold.alts), " ")
		if gold.expTranslation != translated {
			t.Fatalf("Failed case %d. Expected translation is:%s, Translated:%s", i, gold.expTranslation, translated)
		}
	}
}

func TestTranslationWords(t *testing.T) {
	golds := []struct {
		alts     []string
		maxAlt   int
		expected string
	}{
		0: {[]string{"a1", "a2", "a3"}, 1, "a1[alt:a2]"},
		1: {[]string{"a1", "a2", "a3", "a4"}, 5, "a1[alt:a2][alt:a3][alt:a4]"},
		2: {[]string{"a1", "a2"}, 0, "a1"},
	}

	for i, gold := range golds {
		receivedAlts := translationWords(gold.alts, gold.maxAlt)
		if gold.expected != receivedAlts {
			t.Fatalf("Failed case %d. Expected alts:%v, was:%v", i, gold.expected, receivedAlts)
		}
	}
}

func TestIsSupported(t *testing.T) {
	golds := []struct {
		from      string
		to        string
		supported bool
	}{
		0: {"nl", "en", true},
		1: {"nl", "jp", false},
		2: {"en", "nl", false},
	}

	for i, gold := range golds {
		supported := IsSupported(gold.from, gold.to)
		if gold.supported != IsSupported(gold.from, gold.to) {
			t.Fatalf("Failed case %d. Expected support:%v, was:%v", i, gold.supported, supported)
		}
	}
}
