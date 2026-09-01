package textcleaner

import (
	"testing"
)

func TestSplitIntoSentences(t *testing.T) {
	text := "Olá! Este é um texto de teste para avaliar a qualidade de uma voz sintetizada por inteligência artificial. Vamos observar a clareza das palavras, a naturalidade da pronúncia, o ritmo das frases e a variação da entonação. Uma boa voz deve soar agradável, fluida e próxima de uma conversa real."

	chunks := SplitIntoSentences(text, 60)

	if len(chunks) < 2 {
		t.Fatalf("Esperava que o texto fosse dividido em pelo menos 2 frases, obtido: %d", len(chunks))
	}

	// Verificar se os chunks não estão vazios
	for i, c := range chunks {
		if len(c) == 0 {
			t.Errorf("Chunk %d está vazio", i)
		}
	}
}

func TestSplitShortText(t *testing.T) {
	text := "Texto curto."
	chunks := SplitIntoSentences(text, 80)
	if len(chunks) != 1 {
		t.Fatalf("Esperava 1 chunk para texto curto, obtido: %d", len(chunks))
	}
}
