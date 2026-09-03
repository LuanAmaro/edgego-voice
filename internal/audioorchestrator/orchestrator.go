package audioorchestrator

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"

	"edgego-voice/internal/edgetts"
	"edgego-voice/internal/sfx"
)

// SynthesisMetrics sumariza o consumo de caracteres e efeitos nesta execução.
type SynthesisMetrics struct {
	EdgeChars  int
	SFXCount   int
	DurationMs int64
}

// SynthesizeOrchestrated executa a síntese de múltiplos blocos e intercalação de silêncio real/SFX.
func SynthesizeOrchestrated(
	ctx context.Context,
	edgeClient *edgetts.Client,
	text string,
	defaultOpts edgetts.SynthesizeOptions,
	voiceMap map[string]string,
) ([]byte, SynthesisMetrics, error) {
	rawSegments := ParseSegments(text, defaultOpts.Voice, defaultOpts.Rate, defaultOpts.Pitch, voiceMap)
	if len(rawSegments) == 0 {
		return nil, SynthesisMetrics{}, fmt.Errorf("nenhum segmento de áudio para sintetizar")
	}

	segments := optimizeAndMergeSegments(rawSegments)
	metrics := SynthesisMetrics{}

	// Caso otimizado: texto simples sem pausas extras, sem SFX e sem som ambiente
	if len(segments) == 1 && segments[0].Type == SegmentSpeech && segments[0].AmbientTrack == "" {
		opts := defaultOpts
		opts.Text = segments[0].Text
		opts.Voice = segments[0].Voice
		opts.Rate = segments[0].Rate
		opts.Pitch = segments[0].Pitch
		audio, err := edgeClient.Synthesize(ctx, opts)
		if err == nil {
			metrics.EdgeChars = utf8.RuneCountInString(segments[0].Text)
		}
		return audio, metrics, err
	}

	type segmentResult struct {
		index     int
		data      []byte
		edgeChars int
		err       error
	}

	results := make([][]byte, len(segments))
	var wg sync.WaitGroup
	resultChan := make(chan segmentResult, len(segments))

	for i, seg := range segments {
		// 1. Silêncio real local (ou continuidade do som ambiente durante a pausa)
		if seg.Type == SegmentSilence {
			var silenceBytes []byte
			var err error
			if seg.AmbientTrack != "" {
				silenceBytes, err = sfx.GenerateDurationSFXWithVolume(seg.AmbientTrack, defaultOpts.Voice, defaultOpts.Format, seg.Duration, 0.45)
			} else {
				silenceBytes, err = GenerateSilence(defaultOpts.Format, seg.Duration)
			}
			if err != nil {
				return nil, metrics, fmt.Errorf("falha ao gerar silêncio para segmento %d: %w", i, err)
			}
			results[i] = silenceBytes
			continue
		}

		// 2. Efeito Sonoro (SFX) - Carrega arquivo de áudio customizado de ./sfx ou sintetiza acústico puro
		if seg.Type == SegmentSFX {
			var sfxBytes []byte
			var err error
			if seg.Duration > 0 {
				sfxBytes, err = sfx.GenerateDurationSFX(seg.SFXType, defaultOpts.Voice, defaultOpts.Format, seg.Duration)
			} else {
				sfxBytes, err = sfx.GetAudioForSFX(seg.SFXType, defaultOpts.Voice, defaultOpts.Format)
			}
			if err != nil {
				return nil, metrics, fmt.Errorf("falha ao carregar sfx %s: %w", seg.SFXType, err)
			}

			// Se houver uma trilha de ambiente ativa e diferente do SFX atual, mixa o som de ambiente sob o SFX!
			if seg.AmbientTrack != "" && !strings.EqualFold(seg.AmbientTrack, seg.SFXType) {
				if mixed, mixErr := sfx.MixAmbientTrack(sfxBytes, seg.AmbientTrack, defaultOpts.Voice, defaultOpts.Format, 0.45); mixErr == nil && len(mixed) > 0 {
					sfxBytes = mixed
				}
			}

			results[i] = sfxBytes
			metrics.SFXCount++
			continue
		}

		// 3. Síntese de Fala Paralela via Edge TTS
		wg.Add(1)
		go func(idx int, s Segment) {
			defer wg.Done()

			charCount := utf8.RuneCountInString(s.Text)

			segOpts := defaultOpts
			segOpts.Text = s.Text
			segOpts.Voice = s.Voice
			segOpts.Rate = s.Rate
			segOpts.Pitch = s.Pitch
			segOpts.Volume = s.Volume
			segOpts.BreakComma = defaultOpts.BreakComma
			segOpts.BreakPeriod = defaultOpts.BreakPeriod

			audioData, err := edgeClient.Synthesize(ctx, segOpts)
			if err == nil && len(audioData) > 0 && s.AmbientTrack != "" {
				// Se houver trilha de fundo, mixar som ambiente
				if mixed, mixErr := sfx.MixAmbientTrack(audioData, s.AmbientTrack, s.Voice, defaultOpts.Format, 0.45); mixErr == nil && len(mixed) > 0 {
					audioData = mixed
				}
			}

			resultChan <- segmentResult{
				index:     idx,
				data:      audioData,
				edgeChars: charCount,
				err:       err,
			}
		}(i, seg)
	}

	wg.Wait()
	close(resultChan)

	// Coletar resultados preservando a ordem cronológica
	for res := range resultChan {
		if res.err != nil {
			return nil, metrics, fmt.Errorf("falha na síntese do segmento %d: %w", res.index, res.err)
		}
		results[res.index] = res.data
		metrics.EdgeChars += res.edgeChars
	}

	// Estimar tamanho total para pré-alocação de memória (golang-performance)
	totalLen := 0
	for _, chunk := range results {
		totalLen += len(chunk)
	}

	finalAudio := make([]byte, 0, totalLen)
	for _, chunk := range results {
		finalAudio = append(finalAudio, chunk...)
	}

	return finalAudio, metrics, nil
}

// optimizeAndMergeSegments mescla segmentos de fala adjacentes consecutivos de mesma voz para prosódia fluida.
func optimizeAndMergeSegments(segments []Segment) []Segment {
	if len(segments) == 0 {
		return segments
	}

	var merged []Segment
	for _, seg := range segments {
		if len(merged) == 0 {
			merged = append(merged, seg)
			continue
		}

		last := &merged[len(merged)-1]
		if last.Type == SegmentSpeech && seg.Type == SegmentSpeech &&
			last.Voice == seg.Voice &&
			last.Rate == seg.Rate &&
			last.Pitch == seg.Pitch &&
			last.AmbientTrack == seg.AmbientTrack {
			last.Text = strings.TrimSpace(last.Text) + " " + strings.TrimSpace(seg.Text)
		} else {
			merged = append(merged, seg)
		}
	}

	return merged
}
