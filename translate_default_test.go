package main

import (
	"strings"
	"testing"
)

func TestTranslateDefault(t *testing.T) {
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
		translated := strings.Trim(TranslateDefault("", "", gold.text, gold.alts), " ")
		if gold.expTranslation != translated {
			t.Fatalf("Failed case %d. Expected translation is:%s, Translated:%s", i, gold.expTranslation, translated)
		}
	}
}
