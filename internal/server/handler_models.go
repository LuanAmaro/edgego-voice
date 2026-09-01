package server

import (
	"encoding/json"
	"net/http"
)

type ModelItem struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type ModelsResponse struct {
	Data []ModelItem `json:"data"`
}

var availableModels = []ModelItem{
	{ID: "tts-1", Object: "model"},
	{ID: "tts-1-hd", Object: "model"},
}

// ListModelsHandler retorna os modelos disponíveis compatíveis com a OpenAI.
func ListModelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ModelsResponse{
		Data: availableModels,
	})
}
