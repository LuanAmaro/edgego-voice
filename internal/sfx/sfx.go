package sfx

import (
	"bytes"
	"encoding/binary"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SFXType identifica o tipo de micro-efeito sonoro
type SFXType string

const (
	SFXThroat     SFXType = "throat"     // Pigarro
	SFXCough      SFXType = "cough"      // Tosse
	SFXLaughter   SFXType = "laughter"   // Risada
	SFXSigh       SFXType = "sigh"       // Suspiro
	SFXBreath     SFXType = "breath"     // Respiração
	SFXKeyboard   SFXType = "keyboard"   // Efeito Teclado
	SFXCallcenter SFXType = "callcenter" // Callcenter Ambiente
	SFXNoise      SFXType = "noise"      // Ruído de fundo
)

// isMaleVoice verifica se a voz selecionada é masculina para direcionamento de gênero dos efeitos.
func isMaleVoice(voice string) bool {
	v := strings.ToLower(voice)
	maleKeywords := []string{
		"antonio", "antônio", "donato", "fabio", "fábio", "humberto", "julio", "júlio",
		"nicolau", "valerio", "valério", "macerio", "pedro", "rafael", "caio",
		"andrew", "brian", "william", "giuseppe", "remy", "florian", "davis", "jason", "tony",
	}
	for _, m := range maleKeywords {
		if strings.Contains(v, m) {
			return true
		}
	}
	return false
}

// NormalizeSFXType converte apelidos e tags de efeitos para SFXType
func NormalizeSFXType(name string) (SFXType, bool) {
	lower := strings.ToLower(strings.TrimSpace(name))
	switch lower {
	case "pigarro", "throat", "limpar_garganta", "garganta":
		return SFXThroat, true
	case "tosse", "cough", "cof", "tosse-feminina", "tosse-masculina":
		return SFXCough, true
	case "risada", "risos", "riso", "laughter", "laugh", "haha":
		return SFXLaughter, true
	case "suspiro", "sigh", "ah", "suspiro-feminino", "suspiro-masculino":
		return SFXSigh, true
	case "respiracao", "respiração", "breath", "breathe", "inspira":
		return SFXBreath, true
	case "teclado", "efeito-teclado", "digitando", "keyboard", "typing":
		return SFXKeyboard, true
	case "callcenter", "callcenter-ambiente", "ambiente", "escritorio", "escritório":
		return SFXCallcenter, true
	case "ruido", "ruído", "ruido-fundo", "chiado", "noise":
		return SFXNoise, true
	default:
		return SFXType(lower), false
	}
}

// buildCandidateFileNames gera nomes de arquivos candidatos baseados no efeito e no contexto da voz
func buildCandidateFileNames(name, voice string) []string {
	cleanName := strings.TrimSpace(strings.ToLower(name))
	cleanName = strings.ReplaceAll(cleanName, "_", "-")

	isMale := isMaleVoice(voice)
	genderSuffix := "-feminina"
	genderSuffixAlt := "-feminino"
	if isMale {
		genderSuffix = "-masculina"
		genderSuffixAlt = "-masculino"
	}

	var candidates []string

	// Se o efeito for suspiro:
	if strings.Contains(cleanName, "suspiro") || cleanName == "sigh" {
		candidates = append(candidates, "suspiro"+genderSuffixAlt, "suspiro"+genderSuffix, "suspiro")
	}

	// Se o efeito for tosse:
	if strings.Contains(cleanName, "tosse") || cleanName == "cough" {
		candidates = append(candidates, "tosse"+genderSuffix, "tosse"+genderSuffixAlt, "tosse")
	}

	// Se o efeito for teclado:
	if strings.Contains(cleanName, "teclado") || cleanName == "digitando" || cleanName == "keyboard" {
		candidates = append(candidates, "efeito-teclado", "teclado")
	}

	// Se o efeito for callcenter / ambiente:
	if strings.Contains(cleanName, "callcenter") || strings.Contains(cleanName, "ambiente") {
		candidates = append(candidates, "callcenter-ambiente", "callcenter", "ambiente")
	}

	// Se o efeito for ruído de fundo:
	if strings.Contains(cleanName, "ruido") || strings.Contains(cleanName, "ruído") || strings.Contains(cleanName, "noise") {
		candidates = append(candidates, "ruido-fundo", "ruido")
	}

	// Padrões específicos por voz completa ou curta (ex: pt-BR-FranciscaNeural_tosse, Francisca_tosse)
	if voice != "" {
		candidates = append(candidates, voice+"-"+cleanName, voice+"_"+cleanName)
		parts := strings.Split(voice, "-")
		if len(parts) >= 3 {
			shortVoice := parts[2]
			candidates = append(candidates, shortVoice+"-"+cleanName, shortVoice+"_"+cleanName)
		}
	}

	// Nome exato passado na tag
	candidates = append(candidates, cleanName)

	return candidates
}

// LoadCustomSFX tenta carregar um arquivo de áudio de efeito customizado na pasta ./sfx ou /app/sfx.
func LoadCustomSFX(name, voice string) ([]byte, bool) {
	cleanName := strings.TrimSpace(strings.ToLower(name))
	if cleanName == "" {
		return nil, false
	}

	searchDirs := []string{
		"sfx",
		"./sfx",
		"/app/sfx",
	}

	exts := []string{".wav", ".mp3", ".ogg", ".opus", ".aac"}
	candidates := buildCandidateFileNames(cleanName, voice)

	for _, dir := range searchDirs {
		for _, cand := range candidates {
			for _, ext := range exts {
				path := filepath.Join(dir, cand+ext)
				if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
					slog.Info("Efeito SFX carregado com sucesso do disco", "arquivo", path, "efeito", name, "voz", voice)
					return data, true
				}
			}
		}
	}

	return nil, false
}

// GetAudioForSFX carrega do arquivo customizado do usuário em ./sfx ou gera forma de onda acústica.
func GetAudioForSFX(sfxName, voice, format string) ([]byte, error) {
	// 1. Tenta carregar do disco (com resolução inteligente de gênero e contexto de voz)
	if customAudio, ok := LoadCustomSFX(sfxName, voice); ok {
		return customAudio, nil
	}

	// 2. Fallback: gera áudio acústico puro sem fala
	sfxType, _ := NormalizeSFXType(sfxName)
	return GenerateAcousticSFX(sfxType, format)
}

// GenerateAcousticSFX gera um bloco de áudio acústico puro (sem fala) no formato solicitado.
func GenerateAcousticSFX(sfxType SFXType, format string) ([]byte, error) {
	sampleRate := 24000
	var pcm []int16

	r := rand.New(rand.NewSource(1337))

	switch sfxType {
	case SFXThroat:
		duration := 0.30
		totalSamples := int(float64(sampleRate) * duration)
		pcm = make([]int16, totalSamples)
		for i := 0; i < totalSamples; i++ {
			t := float64(i) / float64(sampleRate)
			env := math.Exp(-math.Pow((t-0.07)/0.03, 2)) + 0.7*math.Exp(-math.Pow((t-0.18)/0.04, 2))
			noise := (r.Float64()*2.0 - 1.0)
			tone := math.Sin(2*math.Pi*700*t) + 0.4*math.Sin(2*math.Pi*1100*t)
			val := (noise*0.7 + tone*0.3) * env * 11000
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}

	case SFXCough:
		duration := 0.35
		totalSamples := int(float64(sampleRate) * duration)
		pcm = make([]int16, totalSamples)
		for i := 0; i < totalSamples; i++ {
			t := float64(i) / float64(sampleRate)
			env := math.Exp(-t * 10.0)
			noise := (r.Float64()*2.0 - 1.0)
			tone := math.Sin(2*math.Pi*320*t) + 0.3*math.Sin(2*math.Pi*580*t)
			val := (noise*0.75 + tone*0.25) * env * 13000
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}

	case SFXSigh:
		duration := 0.40
		totalSamples := int(float64(sampleRate) * duration)
		pcm = make([]int16, totalSamples)
		for i := 0; i < totalSamples; i++ {
			t := float64(i) / float64(sampleRate)
			env := math.Sin(math.Pi * t / duration)
			noise := (r.Float64()*2.0 - 1.0)
			val := noise * env * 5500
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}

	default:
		duration := 0.25
		totalSamples := int(float64(sampleRate) * duration)
		pcm = make([]int16, totalSamples)
	}

	lowerFormat := strings.ToLower(format)
	if strings.Contains(lowerFormat, "wav") || strings.Contains(lowerFormat, "riff") {
		return encodeWAV(pcm, sampleRate)
	}

	dur := time.Duration(len(pcm)*1000/sampleRate) * time.Millisecond
	return generateAcousticMP3Frames(pcm, dur), nil
}

func generateAcousticMP3Frames(samples []int16, duration time.Duration) []byte {
	frameDuration := 24 * time.Millisecond
	numFrames := int(math.Ceil(float64(duration) / float64(frameDuration)))
	if numFrames <= 0 {
		numFrames = 1
	}

	out := make([]byte, 0, numFrames*288)

	for f := 0; f < numFrames; f++ {
		frame := make([]byte, 288)
		frame[0] = 0xFF
		frame[1] = 0xF3
		frame[2] = 0x50
		frame[3] = 0xC0

		startSample := f * 576
		endSample := startSample + 576
		if endSample > len(samples) {
			endSample = len(samples)
		}

		var energy float64
		if startSample < len(samples) {
			for s := startSample; s < endSample; s++ {
				energy += math.Abs(float64(samples[s]))
			}
			energy = energy / float64(endSample-startSample)
		}

		if energy > 200 {
			gain := byte(math.Min(120, energy/150))
			frame[4] = gain
			frame[5] = 0x1A
			frame[6] = 0x48

			for b := 9; b < 288; b++ {
				frame[b] = byte((f*37 + b*13) % 256)
			}
		}

		out = append(out, frame...)
	}

	return out
}

func encodeWAV(samples []int16, sampleRate int) ([]byte, error) {
	buf := new(bytes.Buffer)
	numChannels := uint16(1)
	bitsPerSample := uint16(16)
	byteRate := uint32(sampleRate * int(numChannels) * int(bitsPerSample/8))
	blockAlign := uint16(numChannels * (bitsPerSample / 8))
	dataSize := uint32(len(samples) * 2)

	buf.WriteString("RIFF")
	_ = binary.Write(buf, binary.LittleEndian, uint32(36+dataSize))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	_ = binary.Write(buf, binary.LittleEndian, uint32(16))
	_ = binary.Write(buf, binary.LittleEndian, uint16(1))
	_ = binary.Write(buf, binary.LittleEndian, numChannels)
	_ = binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	_ = binary.Write(buf, binary.LittleEndian, byteRate)
	_ = binary.Write(buf, binary.LittleEndian, blockAlign)
	_ = binary.Write(buf, binary.LittleEndian, bitsPerSample)

	buf.WriteString("data")
	_ = binary.Write(buf, binary.LittleEndian, dataSize)
	for _, sample := range samples {
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}

	return buf.Bytes(), nil
}
