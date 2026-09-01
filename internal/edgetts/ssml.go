package edgetts

import (
	"encoding/xml"
	"fmt"
	"strings"
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

// BuildSSML constrói o XML SSML com alocação otimizada via strings.Builder.
func BuildSSML(text, voice, rate, pitch, volume, lang string) string {
	if pitch == "" {
		pitch = "+0Hz"
	}
	if volume == "" {
		volume = "+0%"
	}
	if lang == "" {
		lang = "pt-BR"
	}

	var sb strings.Builder
	// Estimar tamanho total para evitar realocações de memória
	sb.Grow(len(text) + len(voice) + 200)

	sb.WriteString("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='")
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

	_ = xml.EscapeText(&sb, []byte(text))

	sb.WriteString("</prosody></voice></speak>")

	return sb.String()
}
