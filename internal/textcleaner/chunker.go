package textcleaner

import (
	"strings"
	"unicode"
)

// SplitIntoSentences divide um texto em micro-frases/cláusulas naturais (40 a 90 chars) para TTFB instantâneo (~200ms).
func SplitIntoSentences(text string, minChunkLen int) []string {
	text = strings.TrimSpace(text)
	if minChunkLen <= 0 {
		minChunkLen = 45
	}
	if len(text) <= minChunkLen {
		return []string{text}
	}

	var sentences []string
	var current strings.Builder

	runes := []rune(text)
	n := len(runes)

	for i := 0; i < n; i++ {
		r := runes[i]
		current.WriteRune(r)

		// Pontuações fortes (fim de frase) e médias (dois pontos, vírgula após tamanho razoável)
		isStrongPunct := (r == '.' || r == '!' || r == '?' || r == '\n' || r == ';')
		isSoftPunct := (r == ':' || r == ',')

		shouldBreak := false

		if isStrongPunct && current.Len() >= minChunkLen {
			shouldBreak = true
		} else if isSoftPunct && current.Len() >= (minChunkLen + 20) {
			shouldBreak = true
		}

		if shouldBreak {
			// Não quebrar se for número decimal (ex: "1.5" ou "10.000")
			if i+1 < n {
				next := runes[i+1]
				if unicode.IsDigit(next) {
					continue
				}
			}

			chunk := strings.TrimSpace(current.String())
			if chunk != "" {
				sentences = append(sentences, chunk)
				current.Reset()
			}
		}
	}

	if current.Len() > 0 {
		chunk := strings.TrimSpace(current.String())
		if chunk != "" {
			if len(sentences) > 0 && len(chunk) < 25 {
				// Anexar pedaço muito pequeno à última frase
				sentences[len(sentences)-1] += " " + chunk
			} else {
				sentences = append(sentences, chunk)
			}
		}
	}

	if len(sentences) == 0 {
		return []string{text}
	}
	return sentences
}
