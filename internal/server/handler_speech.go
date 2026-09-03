package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"edgego-voice/internal/audiocache"
	"edgego-voice/internal/audioorchestrator"
	"edgego-voice/internal/config"
	"edgego-voice/internal/edgetts"
	"edgego-voice/internal/personas"
	"edgego-voice/internal/textcleaner"
)

type SpeechRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed"`
	Pitch          string  `json:"pitch"`
	BreakComma     string  `json:"break_comma"`
	BreakPeriod    string  `json:"break_period"`
	Telephony      bool    `json:"telephony"`
	AutoBreath     bool    `json:"auto_breath"`
}

type SpeechHandler struct {
	cfg             *config.Config
	ttsClient       *edgetts.Client
	voiceManager    *edgetts.VoiceManager
	personasManager *personas.Manager
	cache           *audiocache.Cache
	singleFlight    *audiocache.Group[[]byte]
}

func NewSpeechHandler(
	cfg *config.Config,
	ttsClient *edgetts.Client,
	vm *edgetts.VoiceManager,
	pm *personas.Manager,
	cache *audiocache.Cache,
) *SpeechHandler {
	return &SpeechHandler{
		cfg:             cfg,
		ttsClient:       ttsClient,
		voiceManager:    vm,
		personasManager: pm,
		cache:           cache,
		singleFlight:    audiocache.NewGroup[[]byte](),
	}
}

// flushWriter envolve http.ResponseWriter para enviar blocos de áudio imediatamente ao cliente (TTFB baixo).
type flushWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (w *flushWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	if w.flusher != nil {
		w.flusher.Flush()
	}
	return n, err
}

func truncateSnippet(s string) string {
	runes := []rune(s)
	if len(runes) > 60 {
		return string(runes[:57]) + "..."
	}
	return s
}

// GenerateSpeechHandler processa a rota OpenAI compatível POST /v1/audio/speech.
func (h *SpeechHandler) GenerateSpeechHandler(w http.ResponseWriter, r *http.Request) {
	var req SpeechRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error(), "invalid_request_error")
		return
	}

	if req.Input == "" {
		writeJSONError(w, http.StatusBadRequest, "O campo 'input' é obrigatório.", "invalid_request_error")
		return
	}

	// Verificar se foi solicitado via header X-Persona-ID
	personaID := r.Header.Get("X-Persona-ID")
	if personaID != "" {
		p, err := h.personasManager.Get(personaID)
		if err == nil {
			h.synthesizeForPersona(w, r, p, req)
			return
		}
	}

	// 1. Sanitizar texto
	text := req.Input
	if !h.cfg.RemoveFilter {
		text = textcleaner.CleanText(text)
	}

	if text == "" {
		writeJSONError(w, http.StatusBadRequest, "O texto fornecido está vazio após a limpeza.", "invalid_request_error")
		return
	}

	// 2. Preencher valores padrão
	voice := req.Voice
	if voice == "" {
		voice = h.cfg.DefaultVoice
	}
	realVoice := h.voiceManager.ResolveVoice(voice)

	format := req.ResponseFormat
	if format == "" {
		format = h.cfg.DefaultFormat
	}

	speed := req.Speed
	if speed == 0 {
		speed = h.cfg.DefaultSpeed
	}

	rate, err := edgetts.SpeedToRate(speed)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error(), "invalid_request_error")
		return
	}

	pitch := req.Pitch
	if pitch == "" {
		pitch = "+0Hz"
	}

	formatInfo := edgetts.GetFormatInfo(format)

	// 3. Checar Cache SLRU em Memória (~0.1ms)
	cacheKey := audiocache.GenerateKeyWithOptions(
		text,
		realVoice,
		format,
		req.Speed,
		pitch,
		h.cfg.RemoveFilter,
		req.Telephony,
		req.AutoBreath,
	)
	if cachedAudio, hit := h.cache.Get(cacheKey); hit {
		w.Header().Set("Content-Type", formatInfo.MimeType)
		w.Header().Set("Content-Length", strconv.Itoa(len(cachedAudio)))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"speech.%s\"", format))
		w.Header().Set("X-Cache", "HIT")
		_, _ = w.Write(cachedAudio)
		return
	}

	// 4. Cache MISS -> Streaming em Tempo Real com http.Flusher
	w.Header().Set("Content-Type", formatInfo.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"speech.%s\"", format))
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)

	var flusher http.Flusher
	if f, ok := w.(http.Flusher); ok {
		flusher = f
		flusher.Flush() // Notificar cliente que a conexão está aberta e pronta imediatamente
	}
	fw := &flushWriter{w: w, flusher: flusher}

	opts := edgetts.SynthesizeOptions{
		Text:            text,
		Voice:           realVoice,
		Rate:            rate,
		Pitch:           pitch,
		Language:        h.cfg.DefaultLanguage,
		Format:          format,
		BreakComma:      req.BreakComma,
		BreakPeriod:     req.BreakPeriod,
		TelephonyFilter: req.Telephony,
		AutoBreath:      req.AutoBreath,
	}

	// Síntese com deduplicação concorrente (singleflight) para evitar picos de CPU e conexões repetidas
	audioData, err, _ := h.singleFlight.Do(cacheKey, func() ([]byte, error) {
		// Dupla verificação no cache caso outra chamada simultânea já tenha concluído
		if cached, hit := h.cache.Get(cacheKey); hit {
			return cached, nil
		}

		data, _, sErr := audioorchestrator.SynthesizeOrchestrated(
			r.Context(),
			h.ttsClient,
			text,
			opts,
			h.voiceManager.GetAllVoices(),
		)
		if sErr == nil && len(data) > 0 {
			h.cache.Set(cacheKey, data)
		}
		return data, sErr
	})

	if err != nil {
		if r.Context().Err() != nil {
			slog.Warn("Cliente cancelou a requisição antes do término da síntese")
			return
		}
		slog.Error("Erro ao sintetizar áudio", "error", err)
		return
	}

	if _, err := fw.Write(audioData); err != nil {
		slog.Error("Erro ao enviar dados de áudio para o cliente", "error", err)
		return
	}
}

// GeneratePersonaSpeechHandler processa a rota dedicada POST /v1/persona/{id}/speech.
func (h *SpeechHandler) GeneratePersonaSpeechHandler(w http.ResponseWriter, r *http.Request) {
	personaID := chi.URLParam(r, "id")
	p, err := h.personasManager.Get(personaID)
	if err != nil {
		if errors.Is(err, personas.ErrPersonaNotFound) {
			writeJSONError(w, http.StatusNotFound, "Persona não encontrada", "not_found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	// Validação de autenticação se a persona tiver chave exclusiva
	if h.cfg.RequireAPIKey {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" || (token != h.cfg.APIKey && (p.APIKey == "" || token != p.APIKey)) {
			writeJSONError(w, http.StatusUnauthorized, "Chave de API inválida para esta Persona.", "authentication_error")
			return
		}
	}

	var req SpeechRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error(), "invalid_request_error")
		return
	}

	if req.Input == "" {
		writeJSONError(w, http.StatusBadRequest, "O campo 'input' é obrigatório.", "invalid_request_error")
		return
	}

	h.synthesizeForPersona(w, r, p, req)
}

func (h *SpeechHandler) synthesizeForPersona(w http.ResponseWriter, r *http.Request, p personas.Persona, req SpeechRequest) {
	text := req.Input
	removeFilter := p.RemoveFilter || h.cfg.RemoveFilter
	if !removeFilter {
		text = textcleaner.CleanText(text)
	}

	if text == "" {
		writeJSONError(w, http.StatusBadRequest, "O texto fornecido está vazio após a limpeza.", "invalid_request_error")
		return
	}

	voice := p.Voice
	if req.Voice != "" {
		voice = req.Voice
	}
	realVoice := h.voiceManager.ResolveVoice(voice)

	format := p.Format
	if req.ResponseFormat != "" {
		format = req.ResponseFormat
	}
	if format == "" {
		format = "mp3"
	}

	speed := p.Speed
	if req.Speed > 0 {
		speed = req.Speed
	}
	if speed <= 0 {
		speed = 1.0
	}

	rate, err := edgetts.SpeedToRate(speed)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error(), "invalid_request_error")
		return
	}

	pitch := p.Pitch
	if req.Pitch != "" {
		pitch = req.Pitch
	}
	if pitch == "" {
		pitch = "+0Hz"
	}

	formatInfo := edgetts.GetFormatInfo(format)

	// Checar Cache SLRU em Memória (~0.1ms)
	cacheKey := audiocache.GenerateKeyWithOptions(text, realVoice, format, speed, pitch, removeFilter, req.Telephony, req.AutoBreath)
	if cachedAudio, hit := h.cache.Get(cacheKey); hit {
		w.Header().Set("Content-Type", formatInfo.MimeType)
		w.Header().Set("Content-Length", strconv.Itoa(len(cachedAudio)))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.%s\"", p.ID, format))
		w.Header().Set("X-Persona-ID", p.ID)
		w.Header().Set("X-Cache", "HIT")
		_, _ = w.Write(cachedAudio)
		return
	}

	// Cache MISS -> Streaming imediato com http.Flusher
	w.Header().Set("Content-Type", formatInfo.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.%s\"", p.ID, format))
	w.Header().Set("X-Persona-ID", p.ID)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Cache", "MISS")

	var flusher http.Flusher
	if f, ok := w.(http.Flusher); ok {
		flusher = f
	}
	fw := &flushWriter{w: w, flusher: flusher}

	breakComma := p.BreakComma
	if req.BreakComma != "" {
		breakComma = req.BreakComma
	}
	breakPeriod := p.BreakPeriod
	if req.BreakPeriod != "" {
		breakPeriod = req.BreakPeriod
	}

	opts := edgetts.SynthesizeOptions{
		Text:            text,
		Voice:           realVoice,
		Rate:            rate,
		Pitch:           pitch,
		Language:        h.cfg.DefaultLanguage,
		Format:          format,
		BreakComma:      breakComma,
		BreakPeriod:     breakPeriod,
		TelephonyFilter: req.Telephony,
		AutoBreath:      req.AutoBreath,
	}

	// Síntese de persona com deduplicação concorrente (singleflight)
	audioData, err, _ := h.singleFlight.Do(cacheKey, func() ([]byte, error) {
		if cached, hit := h.cache.Get(cacheKey); hit {
			return cached, nil
		}

		data, _, sErr := audioorchestrator.SynthesizeOrchestrated(
			r.Context(),
			h.ttsClient,
			text,
			opts,
			h.voiceManager.GetAllVoices(),
		)
		if sErr == nil && len(data) > 0 {
			h.cache.Set(cacheKey, data)
		}
		return data, sErr
	})

	if err != nil {
		if r.Context().Err() != nil {
			slog.Warn("Cliente cancelou a requisição antes do término da síntese")
			return
		}
		slog.Error("Erro ao sintetizar áudio para persona", "persona", p.ID, "error", err)
		return
	}

	if _, err := fw.Write(audioData); err != nil {
		slog.Error("Erro ao enviar áudio de persona para o cliente", "error", err)
		return
	}
}
