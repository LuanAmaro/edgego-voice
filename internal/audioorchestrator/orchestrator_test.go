package audioorchestrator

import (
	"strings"
	"testing"
	"time"
)

func TestParseSegments(t *testing.T) {
	text := "Olá! [pausa: 2s] Este é o primeiro teste. [voz: Antonio] [velocidade: 1.2x] Agora falando mais rápido com outra voz! [pause: 500ms] <break time=\"1.5s\"/> Finalizando."

	voiceMap := map[string]string{
		"pt-BR-AntonioNeural": "Antonio",
		"pt-BR-ThalitaMultilingualNeural": "Thalita",
	}

	segments := ParseSegments(text, "pt-BR-ThalitaMultilingualNeural", "+0%", "+0Hz", voiceMap)

	if len(segments) != 7 {
		t.Fatalf("Esperava 7 segmentos, obteve %d: %+v", len(segments), segments)
	}

	// 1. Fala inicial
	if segments[0].Type != SegmentSpeech || segments[0].Text != "Olá!" || segments[0].Voice != "pt-BR-ThalitaMultilingualNeural" {
		t.Errorf("Segmento 0 incorreto: %+v", segments[0])
	}

	// 2. Silêncio de 2 segundos
	if segments[1].Type != SegmentSilence || segments[1].Duration != 2*time.Second {
		t.Errorf("Segmento 1 (silêncio) incorreto: %+v", segments[1])
	}

	// 3. Fala intermédia
	if segments[2].Type != SegmentSpeech || !strings.Contains(segments[2].Text, "Este é o primeiro teste.") {
		t.Errorf("Segmento 2 incorreto: %+v", segments[2])
	}

	// 4. Fala com voz trocada (Antonio) e velocidade +20%
	if segments[3].Type != SegmentSpeech || segments[3].Voice != "pt-BR-AntonioNeural" || segments[3].Rate != "+20%" {
		t.Errorf("Segmento 3 (troca de voz) incorreto: %+v", segments[3])
	}

	// 5. Pausa de 500ms
	if segments[4].Type != SegmentSilence || segments[4].Duration != 500*time.Millisecond {
		t.Errorf("Segmento 4 (500ms) incorreto: %+v", segments[4])
	}

	// 6. Pausa SSML de 1.5s
	if segments[5].Type != SegmentSilence || segments[5].Duration != 1500*time.Millisecond {
		t.Errorf("Segmento 5 (1.5s SSML) incorreto: %+v", segments[5])
	}

	// 7. Fala final
	if segments[6].Type != SegmentSpeech || segments[6].Text != "Finalizando." {
		t.Errorf("Segmento 6 incorreto: %+v", segments[6])
	}
}

func TestVoiceResolution(t *testing.T) {
	voiceMap := map[string]string{
		"alloy": "pt-BR-ThalitaMultilingualNeural",
		"echo":  "it-IT-GiuseppeMultilingualNeural",
	}

	testCases := []struct {
		input    string
		expected string
	}{
		{"pt-BR-AntonioNeural", "pt-BR-AntonioNeural"},
		{"Antonio", "pt-BR-AntonioNeural"},
		{"Antônio", "pt-BR-AntonioNeural"},
		{"antônio", "pt-BR-AntonioNeural"},
		{"Francisca", "pt-BR-FranciscaNeural"},
		{"Thalita", "pt-BR-ThalitaMultilingualNeural"},
		{"Andrew", "en-US-AndrewMultilingualNeural"},
		{"Álvaro", "es-ES-AlvaroNeural"},
		{"alvaro", "es-ES-AlvaroNeural"},
		{"alloy", "pt-BR-ThalitaMultilingualNeural"},
		{"echo", "it-IT-GiuseppeMultilingualNeural"},
		{"pt-PT-RaquelNeural", "pt-PT-RaquelNeural"},
	}

	fallback := "pt-BR-ThalitaMultilingualNeural"

	for _, tc := range testCases {
		res := resolveVoiceName(tc.input, voiceMap, fallback)
		if res != tc.expected {
			t.Errorf("Para input '%s', esperava '%s', obteve '%s'", tc.input, tc.expected, res)
		}
	}
}

func TestGenerateSilenceMP3(t *testing.T) {
	dur := 1 * time.Second
	silence, err := GenerateSilence("mp3", dur)
	if err != nil {
		t.Fatalf("Erro ao gerar silêncio MP3: %v", err)
	}

	if len(silence) == 0 {
		t.Fatalf("Buffer de silêncio MP3 está vazio")
	}

	// 1s com frames de 24ms (288 bytes) deve ter ~41 frames = ~11808 bytes
	if len(silence) < 10000 {
		t.Errorf("Tamanho do buffer de silêncio menor que o esperado: %d bytes", len(silence))
	}
}

func TestGenerateSilenceWAV(t *testing.T) {
	dur := 500 * time.Millisecond
	wavSilence, err := GenerateSilence("wav", dur)
	if err != nil {
		t.Fatalf("Erro ao gerar silêncio WAV: %v", err)
	}

	if len(wavSilence) < 44 {
		t.Fatalf("Header WAV inválido ou curto demais: %d bytes", len(wavSilence))
	}

	if string(wavSilence[0:4]) != "RIFF" || string(wavSilence[8:12]) != "WAVE" {
		t.Errorf("Assinatura WAV inválida")
	}
}

func TestParseEmotionAndSFXSegments(t *testing.T) {
	text := "Olá a todos! [suspiro] [alegre]Temos excelentes notícias![/alegre] [teclado] Até mais."

	segments := ParseSegments(text, "pt-BR-FranciscaNeural", "+0%", "+0Hz", nil)

	if len(segments) != 5 {
		t.Fatalf("Esperava 5 segmentos, obteve %d: %+v", len(segments), segments)
	}

	// Seg 0: Fala neutra (Edge)
	if segments[0].Provider != ProviderEdge || segments[0].Text != "Olá a todos!" || segments[0].Style != "" {
		t.Errorf("Segmento 0 inválido: %+v", segments[0])
	}

	// Seg 1: SFX Suspiro (Direto)
	if segments[1].Type != SegmentSFX || segments[1].Provider != ProviderSFX || segments[1].SFXType != "suspiro" {
		t.Errorf("Segmento 1 (SFX suspiro) inválido: %+v", segments[1])
	}

	// Seg 2: Fala (Edge)
	if segments[2].Provider != ProviderEdge || segments[2].Text != "Temos excelentes notícias!" {
		t.Errorf("Segmento 2 inválido: %+v", segments[2])
	}

	// Seg 3: SFX Teclado (Direto)
	if segments[3].Type != SegmentSFX || segments[3].Provider != ProviderSFX || segments[3].SFXType != "teclado" {
		t.Errorf("Segmento 3 (SFX teclado) inválido: %+v", segments[3])
	}

	// Seg 4: Fala neutra (Edge)
	if segments[4].Provider != ProviderEdge || segments[4].Text != "Até mais." || segments[4].Style != "" {
		t.Errorf("Segmento 4 inválido: %+v", segments[4])
	}
}

func TestParseAmbientBlockAndDurationSFX(t *testing.T) {
	// 1. Bloco de ambiente envolvente
	text1 := "[callcenter] Olá tudo bem, como posso ajudar você hoje? [/callcenter]"
	segs1 := ParseSegments(text1, "pt-BR-FranciscaNeural", "+0%", "+0Hz", nil)
	if len(segs1) != 1 {
		t.Fatalf("Esperava 1 segmento com trilha ambiente, obteve %d", len(segs1))
	}
	if segs1[0].AmbientTrack != "callcenter" {
		t.Errorf("Esperava AmbientTrack=callcenter, obteve: %s", segs1[0].AmbientTrack)
	}

	// 2. SFX com duração programada
	text2 := "Só um momento vou consultar no sistema [teclado:5s] Pronto!"
	segs2 := ParseSegments(text2, "pt-BR-FranciscaNeural", "+0%", "+0Hz", nil)
	if len(segs2) != 3 {
		t.Fatalf("Esperava 3 segmentos, obteve %d", len(segs2))
	}
	if segs2[1].Type != SegmentSFX || segs2[1].SFXType != "teclado" || segs2[1].Duration != 5*time.Second {
		t.Errorf("Segmento 1 de SFX com duração inválido: %+v", segs2[1])
	}
}

