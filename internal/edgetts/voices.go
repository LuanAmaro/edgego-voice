package edgetts

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

// VoiceManager gerencia os apelidos e mapeamento de vozes da OpenAI para Microsoft Edge Neural.
type VoiceManager struct {
	mu     sync.RWMutex
	voices map[string]string
}

var defaultVoices = map[string]string{
	"alloy":   "pt-BR-ThalitaMultilingualNeural",
	"ash":     "en-US-AndrewMultilingualNeural",
	"ballad":  "en-US-BrianMultilingualNeural",
	"cedar":   "en-AU-WilliamMultilingualNeural",
	"coral":   "en-US-AvaMultilingualNeural",
	"echo":    "it-IT-GiuseppeMultilingualNeural",
	"fable":   "en-US-EmmaMultilingualNeural",
	"onyx":    "fr-FR-RemyMultilingualNeural",
	"nova":    "de-DE-SeraphinaMultilingualNeural",
	"sage":    "fr-FR-VivienneMultilingualNeural",
	"shimmer": "de-DE-FlorianMultilingualNeural",
	"verse":   "ko-KR-HyunsuMultilingualNeural",
}

// NewVoiceManager inicializa o gerenciador carregando de voices.json se disponível.
func NewVoiceManager(filePath string) *VoiceManager {
	vm := &VoiceManager{
		voices: make(map[string]string, len(defaultVoices)),
	}

	// Preencher com defaults primeiro
	for k, v := range defaultVoices {
		vm.voices[k] = v
	}

	if filePath == "" {
		filePath = "voices.json"
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Warn("Ficheiro voices.json não encontrado ou inacessível, utilizando mapeamento padrão em memória", "path", filePath, "error", err)
		return vm
	}

	var customVoices map[string]string
	if err := json.Unmarshal(data, &customVoices); err != nil {
		slog.Error("Erro ao decodificar voices.json, mantendo mapeamento padrão", "error", err)
		return vm
	}

	for k, v := range customVoices {
		vm.voices[k] = v
	}

	slog.Info("Mapeamento de vozes carregado com sucesso", "total_vozes", len(vm.voices))
	return vm
}

// ResolveVoice retorna a voz real do Edge correspondente ao apelido ou o próprio nome fornecido.
func (vm *VoiceManager) ResolveVoice(aliasOrName string) string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	if realVoice, ok := vm.voices[aliasOrName]; ok {
		return realVoice
	}
	return aliasOrName
}
