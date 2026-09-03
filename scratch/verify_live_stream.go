package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type StreamMessage struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Voice    string `json:"voice,omitempty"`
	Sentence string `json:"sentence,omitempty"`
}

func main() {
	wsURL := "ws://localhost:5050/v1/audio/stream?key=minha-chave-secreta"
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Falha ao conectar via WebSocket: %v (status=%v)", err, resp)
	}
	defer conn.Close()

	fmt.Println("Conectado com sucesso ao WebSocket de streaming!")

	// 1. Enviar config inicial
	_ = conn.WriteJSON(StreamMessage{
		Type:  "config",
		Voice: "pt-BR-FranciscaNeural",
	})

	// 2. Simular tokens vindos de um LLM com tags intercaladas
	tokens := []string{
		"[callcenter]",
		"Olá! ",
		"Seja bem-vindo ",
		"ao suporte EdgeGo. ",
		"[pausa: 500ms] ",
		"Deixe-me verificar seu sistema ",
		"[teclado:2s]. ",
		"Pronto, tudo certo! ",
		"[/callcenter]",
	}

	go func() {
		for _, t := range tokens {
			time.Sleep(100 * time.Millisecond) // simula latência de geração de tokens da LLM
			_ = conn.WriteJSON(StreamMessage{
				Type: "text",
				Text: t,
			})
		}
		time.Sleep(200 * time.Millisecond)
		_ = conn.WriteJSON(StreamMessage{Type: "flush"})
	}()

	totalBytes := 0
	numAudioChunks := 0

	for {
		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.BinaryMessage {
			numAudioChunks++
			totalBytes += len(data)
			fmt.Printf(" [AUDIO CHUNK #%d] Recebidos %d bytes de áudio binário!\n", numAudioChunks, len(data))
		} else if msgType == websocket.TextMessage {
			var msg StreamMessage
			_ = json.Unmarshal(data, &msg)
			if msg.Type == "sentence_start" {
				fmt.Printf(" >>> Servidor iniciou síntese da sentença: %q\n", msg.Sentence)
			} else if msg.Type == "sentence_end" {
				fmt.Printf(" <<< Servidor concluiu síntese da sentença: %q\n", msg.Sentence)
			} else if msg.Type == "done" {
				fmt.Println(" === STREAM CONCLUÍDO COM SUCESSO! ===")
				break
			}
		}
	}

	fmt.Printf("\nResultado Final: %d pacotes de áudio, total de %d bytes recebidos.\n", numAudioChunks, totalBytes)
}
