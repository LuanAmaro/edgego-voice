# 📜 Histórico de Versões (Changelog)

Todas as mudanças notáveis no projeto **EdgeGo Voice** são documentadas neste arquivo.

O formato é baseado no [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/) e adere ao [Versionamento Semântico](https://semver.org/lang/pt-BR/).

---

## [Não Lançado] (Unreleased)

### ✨ Adicionado
- **Suporte a Multi-idiomas (i18n):** Adicionado seletor de idioma integrado no painel com suporte a **Português (Brasil)** [padrão], **Inglês (Estados Unidos)** e **Espanhol (Espanha)**, com persistência automática no `localStorage`.
- **Modal de Seleção de Idioma (shadcn/ui):** Interface construída com os componentes e padrões do **shadcn/ui** exibindo as bandeiras dos países (`CountryFlag`), nomes nativos e feedback visual de seleção.
- **Seletor no Topbar / Navbar:** O status estático foi substituído por um botão interativo que exibe a bandeira do país atual, código da linguagem e o indicador pulsante verde de status online.
- **Modal de Exclusão de Personas Seguro e Minimalista:** Confirmação de deleção exigindo a digitação do slug exato da persona para evitar exclusões acidentais, seguindo a composição e padrões do **shadcn/ui**.
- **Configuração de Versionamento (`.gitignore`):** Regras completas de ignore para binários Go, artefatos de build do Next.js (`frontend/.next`, `frontend/out`), variáveis de ambiente locais e arquivos de áudio temporários.

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
