# 📜 Histórico de Versões (Changelog)

Todas as mudanças notáveis no projeto **EdgeGo Voice** são documentadas neste arquivo.

O formato é baseado no [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/) e adere ao [Versionamento Semântico](https://semver.org/lang/pt-BR/).

---

## [Não Lançado] (Unreleased)

### ✨ Adicionado
- **Voice Studio Workstation (Redesign Completo do Playground):**
  - **Layout de Estúdio Unificado:** Reformulação da interface em formato de estação de trabalho profissional (Canvas Central de Criação + Inspetor Lateral de Voz e Parâmetros), eliminando caixas vazias e múltiplos cards desconexos.
  - **Componente Oficial AI Elements Persona (`@rive-app/react-webgl2`):** Orbe animado em WebGL2 integrado organicamente no topo do Inspetor, reagindo em tempo real aos estados `idle` (respiração fluida), `thinking` (rotação ativa durante a geração) e `speaking` (ondulação sincronizada com o player de áudio).
  - **Player de Áudio Acoplado:** O Waveform Player surge suavemente diretamente abaixo do botão de ação, proporcionando feedback sonoro imediato no fluxo de criação.
  - **Atalho de Teclado:** Suporte a `Ctrl + Enter` (ou `Cmd + Enter`) para disparar a síntese de áudio instantaneamente.
  - **Design Minimalista com Ionicons:** Substituição de 100% dos emojis da barra de ferramentas e menus por Ionicons monocromáticos de traço limpo, alinhados à estética do **shadcn/ui**, Linear e ElevenLabs.
- **Motor de Efeitos Sonoros Acústicos e Sons Ambiente (`internal/sfx`):**
  - **Trilhas de Ambiente Contínuas:** Mixagem de fundo transparente com atenuação calibrada (ganho 0.45) via tags de envelope (ex: `[callcenter]...[/callcenter]`, `[ruido]`, `[ambiente]`).
  - **Simulação Procedural de Digitação Humana:** Algoritmo orgânico que alterna entre múltiplos arquivos de teclado (`teclado-to-spech.wav`, `teclado-to-spech2.wav`, `efeito-teclado.wav`), intercalando rajadas de teclas (0.35s a 0.70s) com micro-pausas naturais de reflexão (250ms a 500ms).
  - **Camadas Simultâneas:** Continuidade do áudio ambiente por baixo de pausas reais e durante a execução de efeitos sonoros.
  - **Alinhamento Rigoroso de Amostragem (24kHz Mono 48kbps):** Padronização exata de clock e taxa de bits com o Edge TTS para evitar saltos ou distorções de duração na reprodução dos navegadores.
- **Controles de Pausas Naturais (SSML) no Inspetor:**
  - Controles deslizantes dedicados para pausas em vírgulas (`break_comma`: 0 a 800ms) e pontos finais (`break_period`: 0 a 1200ms) com badge indicativo de custo zero ($0).
- **Orquestrador de Silêncio Real e Concatenação em Go (`internal/audioorchestrator`):**
  - **Geração de Silêncio Puro:** Criação de frames válidos de silêncio absoluto em Go para MP3, WAV e PCM com precisão em milissegundos (ex: `[pausa: 2s]`, `[pausa: 500ms]`).
  - **Suporte a Diálogos e Multi-Vozes:** Troca dinâmica de personagens no meio do texto com a tag `[voz: Nome]`.
  - **Síntese Concorrente Paralela:** Processamento concorrente de múltiplos blocos de áudio via goroutines (`golang-concurrency`) e concatenação em buffer pré-alocado (`golang-performance`).
- **Orquestrador SSML Visual & Injetor de Prosódia:**
  - **Ajuste Fino de Tom (Pitch):** Controle deslizante de afinação de frequência (`-20Hz` a `+20Hz`) integrado ao formulário de Personas e ao Voice Studio.
  - **Sanitizador e Parser SSML Inteligente:** Conversão de tags de pausas e formatação em cadências suportadas pelo motor nativo.
- **Suporte a Multi-idiomas (i18n):** Seletor de idioma integrado no painel com suporte a Português (Brasil), Inglês (EUA) e Espanhol (Espanha).

### 🗑️ Removido
- **Integração com Microsoft Azure Speech API:**
  - Excluído o pacote `internal/azuretts/` (`client.go` e `client_test.go`).
  - Removidas as variáveis de ambiente `AZURE_SPEECH_KEY` e `AZURE_SPEECH_REGION` de `config.go`, `docker-compose.yml` e `.env.example`.
  - Eliminado rastreamento de custos em dólar, estimativas financeiras e tokens pagos, operando exclusivamente em modo **100% Go Nativo Gratuito ($0.00)**.
  - Removido cabeçalho flutuante redundante `PersonaHeroCard.tsx`.
- **Página e Menu de Métricas & Consumo:**
  - Removida a aba de navegação "Métricas" na barra superior (`Navbar.tsx`).
  - Excluída a tela de visualização de métricas (`AnalyticsView.tsx`) e limpos os tipos associados (`RequestLogItem`, `TimeBucket`, `UsageStats`), deixando a navegação focada em Personas, Voice Studio e Documentação da API.

---

## [1.0.0] - 2026-09-01

### 🚀 Lançamento Oficial: EdgeGo Voice (Go Native + Next.js)

Reescrita completa da arquitetura do projeto para uma plataforma nativa em **Go (Golang 1.22)** com painel administrativo integrado em **Next.js 14**, substituindo o protótipo legado em Python por um motor de altíssima performance e baixo consumo de recursos.

### ✨ Funcionalidades (Features)
- **Motor TTS em Golang 1.22:**
  - Streaming em tempo real de áudio binário via `http.Flusher` conectado diretamente aos WebSockets do Microsoft Edge TTS.
  - Consumo de memória ultrabaixo (~20MB de RAM) e alto throughput para concorrência massiva sem bloqueios de GIL.
- **100% Compatível com a API OpenAI (`/v1/audio/speech`):**
  - *Drop-in replacement* para SDKs oficiais da OpenAI (Python, Node.js, PHP, cURL) e plataformas de automação (n8n, Typebot, Dify, Flowise, Evolution API).
- **Sistema de Personas de Áudio:**
  - Gerenciamento de perfis pré-configurados de áudio (voz, formato, velocidade e pitch).
  - Rota dedicada de síntese direta: `POST /v1/persona/{id}/speech`.
  - Persistência atômica em arquivo JSON com suporte a concorrência (`sync.RWMutex`).
- **Cache LRU em Memória Thread-Safe:**
  - Respostas imediatas em **< 1ms** para requisições repetidas.
  - Chaveamento criptográfico SHA-256 (texto + voz + velocidade + formato + filtros).
  - Limite configurável de memória máxima (MB) e tempo de expiração TTL.
- **Higienizador Inteligente de Texto (Text Cleaner):**
  - Remoção automática de marcações Markdown (`**negrito**`, `# títulos`, links), URLs e caracteres de controle para fala limpa e natural.
- **Painel Administrativo SPA (Next.js 14):**
  - Dashboard dark mode com estética refinada e componentes **shadcn/ui**.
  - **Playground Interativo:** Síntese em tempo real com estatísticas de latência (ms), tamanho do arquivo, taxa de transferência e player de áudio com visualização de ondas (Waveform).
  - **Gerador de Snippets de Integração:** Exemplos de código prontos para cURL, Python (OpenAI SDK), JavaScript/TypeScript, Go e automações via Webhook (n8n/Evolution API).
  - **Catálogo de Vozes Neurais Multilíngues:** Suporte a vozes de alta fidelidade em português (Brasil e Portugal), inglês, espanhol, francês, alemão, italiano e japonês com identificação visual por bandeiras de países.
- **Distribuição em Container Único:**
  - Imagem Docker Alpine minimalista (~18MB) compilando o binário Go e servindo o bundle estático do Next.js na mesma porta (`5050`).

---

## [0.2.0] - 2025-10-02 (Legado - Protótipo Python)

- Suporte a vozes neurais multilinguais e mapeamento via `voices.json`.
- Conversão de formato de áudio via `pydub` e `ffmpeg`.
