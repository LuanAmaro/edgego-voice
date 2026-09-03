package edgetts

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// SSMLOptions define as configurações completas para geração do SSML.
type SSMLOptions struct {
	Voice       string
	Rate        string // ex: "+0%", "+50%", "-10%"
	Pitch       string // ex: "+0Hz", "+5Hz", "-10Hz", "+5%"
	Volume      string // ex: "+0%"
	Lang        string // ex: "pt-BR"
	BreakComma  string // ex: "150ms" (pausa automática após vírgulas)
	BreakPeriod string // ex: "350ms" (pausa automática após pontuações finais)
}

var (
	// reBreakTag detecta tags <break...> para converter em pausas naturais aceitas pelo motor
	reBreakTag = regexp.MustCompile(`(?i)\s*<break\s+[^>]*\/?>\s*`)
	// reSSMLGeneric detecta outras tags SSML para higienização limpa
	reSSMLGeneric = regexp.MustCompile(`(?i)</?(?:say-as|emphasis|sub|phoneme|prosody|voice|speak|s|p|lang|audio|mstts:[a-z-]+)\b[^>]*>`)
	// reComma detecta vírgulas após letras/palavras seguidas de espaço ou fim de linha, sem quebrar números como 10,50
	reComma = regexp.MustCompile(`([a-zA-Z\p{L}]),\s+`)
	// rePeriod detecta pontuações finais (. ! ?) após letras/números seguidas de espaço ou fim de linha, sem quebrar casas decimais como 3.14
	rePeriod = regexp.MustCompile(`([a-zA-Z\p{L}0-9])([.!?]+)(\s+|$)`)
)

// SpeedToRate converte um multiplicador de velocidade (ex: 1.5) para a string de taxa do SSML (ex: "+50%").
func SpeedToRate(speed float64) (string, error) {
	if speed < 0.25 || speed > 2.0 {
		return "", fmt.Errorf("velocidade deve estar entre 0.25 e 2.0, recebido: %.2f", speed)
	}
	percentage := int((speed - 1.0) * 100)
	if percentage >= 0 {
		return fmt.Sprintf("+%d%%", percentage), nil
	}
	return fmt.Sprintf("%d%%", percentage), nil
}

// BuildSSML constrói o XML SSML legado (mantido para compatibilidade).
func BuildSSML(text, voice, rate, pitch, volume, lang string) string {
	return BuildSSMLWithOptions(text, SSMLOptions{
		Voice:  voice,
		Rate:   rate,
		Pitch:  pitch,
		Volume: volume,
		Lang:   lang,
	})
}

// BuildSSMLWithOptions constrói o XML SSML com suporte a pausas inteligentes e afinação de tom.
func BuildSSMLWithOptions(text string, opts SSMLOptions) string {
	pitch := normalizePitch(opts.Pitch)
	volume := opts.Volume
	if volume == "" {
		volume = "+0%"
	}
	rate := opts.Rate
	if rate == "" {
		rate = "+0%"
	}
	lang := opts.Lang
	if lang == "" {
		lang = "pt-BR"
	}
	voice := opts.Voice

	// 1. Converter tags <break> para pausas naturais e higienizar tags XML internas
	cleanInput := reBreakTag.ReplaceAllString(text, " ... ")
	cleanInput = reSSMLGeneric.ReplaceAllString(cleanInput, "")
	cleanInput = strings.TrimSpace(cleanInput)

	// 2. Injetar pausas calculadas em pontuações caso configurado
	processedText := injectNaturalPauses(cleanInput, opts.BreakComma, opts.BreakPeriod)

	var sb strings.Builder
	// Estimar tamanho total para evitar realocações de memória (golang-performance)
	sb.Grow(len(processedText) + len(voice) + 240)

	sb.WriteString("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xmlns:mstts='https://www.w3.org/2001/mstts' xmlns:emo='http://www.w3.org/2009/10/emotionml' xml:lang='")
	sb.WriteString(lang)
	sb.WriteString("'><voice name='")
	sb.WriteString(voice)
	sb.WriteString("'><prosody pitch='")
	sb.WriteString(pitch)
	sb.WriteString("' rate='")
	sb.WriteString(rate)
	sb.WriteString("' volume='")
	sb.WriteString(volume)
	sb.WriteString("'>")

	_ = xml.EscapeText(&sb, []byte(processedText))

	sb.WriteString("</prosody></voice></speak>")

	return sb.String()
}

// normalizeDuration padroniza strings de tempo (ex: "150" -> "150ms").
func normalizeDuration(d string) string {
	d = strings.TrimSpace(d)
	if d == "" || d == "0" || d == "0ms" || d == "0s" {
		return ""
	}
	if !strings.HasSuffix(d, "ms") && !strings.HasSuffix(d, "s") {
		return d + "ms"
	}
	return d
}

// normalizePitch padroniza o valor de tom (ex: "5" -> "+5Hz", "-5" -> "-5Hz", "+10%" -> "+10%").
func normalizePitch(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p == "0" || p == "0Hz" || p == "0%" {
		return "+0Hz"
	}
	if strings.HasSuffix(p, "Hz") || strings.HasSuffix(p, "%") || strings.HasSuffix(p, "st") {
		if !strings.HasPrefix(p, "+") && !strings.HasPrefix(p, "-") {
			return "+" + p
		}
		return p
	}
	if !strings.HasPrefix(p, "+") && !strings.HasPrefix(p, "-") {
		return "+" + p + "Hz"
	}
	return p + "Hz"
}

// injectNaturalPauses injeta pausas naturais em pontuações finais e vírgulas sem corromper números decimais.
func injectNaturalPauses(text, breakComma, breakPeriod string) string {
	breakComma = normalizeDuration(breakComma)
	breakPeriod = normalizeDuration(breakPeriod)

	if breakComma == "" && breakPeriod == "" {
		return text
	}

	result := text
	if breakComma != "" && breakComma != "0ms" {
		// Substituição sutil para respiro em vírgulas
		result = reComma.ReplaceAllString(result, "$1, ")
	}
	if breakPeriod != "" && breakPeriod != "0ms" {
		// Substituição para pausa estendida em pontuação final
		result = rePeriod.ReplaceAllString(result, "$1$2.. $3")
	}
	return result
}


