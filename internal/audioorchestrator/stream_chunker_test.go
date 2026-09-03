package audioorchestrator

import (
	"strings"
	"testing"
)

func TestChunker_FragmentedTags(t *testing.T) {
	chunker := NewTagAwareChunker()

	// Simulação de tokens do LLM chegando fragmentados no meio de tags
	tokens := []string{
		"[pa", "usa: ", "500", "ms] ",
		"Olá ", "seja ", "bem-vindo! ",
		"Só um ", "instante ", "enquanto ", "consulto ", "o sistema ",
		"[tec", "lado:2s]. ",
		"Pronto ", "encontrei ", "seus dados.",
	}

	var emitted []string
	for _, tok := range tokens {
		chunks := chunker.Feed(tok)
		emitted = append(emitted, chunks...)
	}
	emitted = append(emitted, chunker.Flush()...)

	if len(emitted) < 2 {
		t.Fatalf("Esperava pelo menos 2 sentenças emitidas, obteve %d: %+v", len(emitted), emitted)
	}

	// Verifica se a primeira sentença preservou a tag de pausa inteira
	if !strings.Contains(emitted[0], "[pausa: 500ms]") {
		t.Errorf("Tag de pausa foi corrompida no streaming: %s", emitted[0])
	}

	// Verifica se a segunda sentença preservou a tag de teclado inteira
	combined := strings.Join(emitted, " ")
	if !strings.Contains(combined, "[teclado:2s]") {
		t.Errorf("Tag de teclado com duração foi corrompida: %s", combined)
	}
}

func TestChunker_EnvelopePreservation(t *testing.T) {
	chunker := NewTagAwareChunker()

	tokens := []string{
		"[callcenter]", "Olá! ", "Tudo bem? ", "Como posso ajudar? ", "[/callcenter]",
	}

	var emitted []string
	for _, tok := range tokens {
		chunks := chunker.Feed(tok)
		emitted = append(emitted, chunks...)
	}
	emitted = append(emitted, chunker.Flush()...)

	if len(emitted) < 2 {
		t.Fatalf("Esperava sentenças separadas, obteve %d: %+v", len(emitted), emitted)
	}

	// Todas as sentenças geradas dentro do envelope callcenter devem manter a mixagem de ambiente
	for i, s := range emitted {
		if !strings.Contains(s, "[callcenter]") {
			t.Errorf("Sentença %d perdeu o envelope de ambiente: %s", i, s)
		}
	}
}

func TestChunker_FlushRemaining(t *testing.T) {
	chunker := NewTagAwareChunker()

	// Texto sem pontuação terminal no final
	chunks := chunker.Feed("Olá este é um teste sem ponto final")
	if len(chunks) != 0 {
		t.Fatalf("Não deveria ter emitido antes do fim da frase")
	}

	final := chunker.Flush()
	if len(final) != 1 || !strings.Contains(final[0], "sem ponto final") {
		t.Fatalf("Flush deveria ter emitido o texto restante: %+v", final)
	}
}
