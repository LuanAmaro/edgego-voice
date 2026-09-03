package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"edgego-voice/internal/audiocache"
	"edgego-voice/internal/config"
	"edgego-voice/internal/edgetts"
	"edgego-voice/internal/personas"
	"edgego-voice/internal/server"
)

func main() {
	// 1. Configurar Logger Estruturado
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logHandler))

	// 2. Carregar Configurações
	cfg := config.Load()
	if cfg.DetailedErrorLogging {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	// 3. Inicializar Componentes de Negócio, Clientes TTS & Cache LRU
	voiceManager := edgetts.NewVoiceManager("voices.json")
	personasManager := personas.NewManager("personas.json")
	ttsClient := edgetts.NewClient(cfg.Proxy)
	audioCache := audiocache.NewCache(cfg.CacheEnabled, cfg.CacheMaxMB, cfg.CacheTTL)

	// 4. Montar Roteador e Servidor HTTP
	router := server.SetupRouter(cfg, ttsClient, voiceManager, personasManager, audioCache)

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 5. Exibir Banner de Inicialização
	fmt.Println("=====================================================================")
	fmt.Println("      🔊 Iniciando o EdgeGo Voice Engine 🇧🇷 (100% Go Nativo - $0)")
	fmt.Println("=====================================================================")
	fmt.Println("🐹 Linguagem: Go (Golang) - Alta Performance e Baixa Latência")
	fmt.Println("⚙️  Servidor: Go net/http + Chi Router (Streaming)")
	fmt.Println("🎭 Módulo de Personas: Ativo (CRUD & API dedicadas)")
	fmt.Printf("🧠 Cache de Áudio LRU: %v (Max: %d MB)\n", cfg.CacheEnabled, cfg.CacheMaxMB)
	fmt.Println("🟢 Motor de Áudio: Edge TTS + SFX Mixer (100% Gratuito / Zero Tokens)")
	fmt.Printf("🌐 Servidor / Painel Web: http://%s\n", addr)
	fmt.Printf("🔑 Exigir chave de API: %v\n", cfg.RequireAPIKey)
	fmt.Printf("🎤 Voz Padrão: %s\n", cfg.DefaultVoice)
	fmt.Println("=====================================================================")

	// 6. Iniciar Servidor em Goroutine com Graceful Shutdown
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Falha fatal ao iniciar o servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	// Aguardar sinais de término do SO
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Encerrando servidor gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("Erro durante o encerramento do servidor", "error", err)
	}

	slog.Info("Servidor finalizado com sucesso.")
}
