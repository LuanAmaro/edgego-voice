package sfx

import (
	"testing"
)

func TestNormalizeSFXType(t *testing.T) {
	testCases := []struct {
		input    string
		expected SFXType
		valid    bool
	}{
		{"respiracao", SFXBreath, true},
		{"respiração", SFXBreath, true},
		{"respiro", SFXBreath, true},
		{"breath", SFXBreath, true},
		{"hum", SFXType("hum"), true},
		{"hmm", SFXType("hum"), true},
		{"filler:hmm", SFXType("hum"), true},
		{"entendi", SFXType("entendi"), true},
		{"certo", SFXType("certo"), true},
		{"deixa-ver", SFXType("deixa-ver"), true},
		{"callcenter", SFXCallcenter, true},
		{"teclado", SFXKeyboard, true},
		{"tosse", SFXCough, true},
		{"suspiro", SFXSigh, true},
		{"ruido", SFXNoise, true},
		{"desconhecido", SFXType("desconhecido"), false},
	}

	for _, tc := range testCases {
		res, ok := NormalizeSFXType(tc.input)
		if ok != tc.valid {
			t.Errorf("Para input '%s', esperava valid=%v, obteve %v", tc.input, tc.valid, ok)
		}
		if res != tc.expected {
			t.Errorf("Para input '%s', esperava '%s', obteve '%s'", tc.input, tc.expected, res)
		}
	}
}

func TestBuildCandidateFileNames(t *testing.T) {
	// 1. Respiração com voz feminina
	femaleCandidates := buildCandidateFileNames("respiracao", "pt-BR-FranciscaNeural")
	foundFem := false
	for _, c := range femaleCandidates {
		if c == "suspiro-feminino" || c == "suspiro-feminina" {
			foundFem = true
			break
		}
	}
	if !foundFem {
		t.Errorf("Esperava encontrar candidato feminino para respiração, obteve: %+v", femaleCandidates)
	}

	// 2. Respiração com voz masculina
	maleCandidates := buildCandidateFileNames("respiracao", "pt-BR-AntonioNeural")
	foundMasc := false
	for _, c := range maleCandidates {
		if c == "suspiro-masculino" || c == "suspiro-masculina" {
			foundMasc = true
			break
		}
	}
	if !foundMasc {
		t.Errorf("Esperava encontrar candidato masculino para respiração, obteve: %+v", maleCandidates)
	}
}

func TestApplyTelephonyFilter_EmptyAudio(t *testing.T) {
	out, err := ApplyTelephonyFilter(nil, "mp3")
	if err != nil {
		t.Fatalf("Erro inesperado com áudio vazio: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("Esperava áudio vazio, obteve %d bytes", len(out))
	}
}
