package edgetts

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	winEpoch           = 11644473600
	trustedClientToken = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	chromiumVersion    = "143.0.3650.75"
	secMSGECVersion    = "1-" + chromiumVersion
	edgeEndpoint       = "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1"
	originHeader       = "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold"
	userAgentHeader    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 Edg/143.0.0.0"
)

// Client gerencia conexões de síntese no Microsoft Edge TTS com alta performance.
type Client struct {
	proxyURL *url.URL
	mu       sync.Mutex
}

// SynthesizeOptions define os parâmetros para a chamada de síntese.
type SynthesizeOptions struct {
	Text        string
	Voice       string
	Rate        string
	Pitch       string
	Volume      string
	Language    string
	Format      string
	BreakComma      string
	BreakPeriod     string
	TelephonyFilter bool
	AutoBreath      bool
}

// NewClient cria uma nova instância do cliente Edge TTS.
func NewClient(proxyStr string) *Client {
	c := &Client{}
	if proxyStr != "" {
		if u, err := url.Parse(proxyStr); err == nil {
			c.proxyURL = u
		}
	}
	return c
}

// extractLangFromVoice extrai o idioma a partir do identificador da voz (ex: "en-US" de "en-US-Andrew...").
func extractLangFromVoice(voice string) string {
	parts := strings.Split(voice, "-")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s-%s", parts[0], parts[1])
	}
	return "pt-BR"
}

// generateSecMSGEC gera o token DRM baseado no timestamp Windows epoch e hash SHA256.
func generateSecMSGEC() string {
	ticks := time.Now().UTC().Unix()
	ticks += winEpoch
	ticks -= ticks % 300
	ticks *= 10000000 // intervalos de 100 nanossegundos

	strToHash := fmt.Sprintf("%d%s", ticks, trustedClientToken)
	hash := sha256.Sum256([]byte(strToHash))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// generateMUID gera um ID aleatório para o cabeçalho Cookie.
func generateMUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}

// dialConnection cria uma conexão WebSocket direta e otimizada com TCP_NODELAY.
func (c *Client) dialConnection(ctx context.Context) (*websocket.Conn, error) {
	connectionID := strings.ReplaceAll(uuid.New().String(), "-", "")
	secMSGEC := generateSecMSGEC()

	wsURL := fmt.Sprintf(
		"%s?TrustedClientToken=%s&Sec-MS-GEC=%s&Sec-MS-GEC-Version=%s&ConnectionId=%s",
		edgeEndpoint, trustedClientToken, secMSGEC, secMSGECVersion, connectionID,
	)

	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  10 * time.Second,
		EnableCompression: true,
		NetDialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	if c.proxyURL != nil {
		dialer.Proxy = http.ProxyURL(c.proxyURL)
	}

	headers := http.Header{}
	headers.Set("Origin", originHeader)
	headers.Set("User-Agent", userAgentHeader)
	headers.Set("Pragma", "no-cache")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Cookie", fmt.Sprintf("muid=%s", generateMUID()))

	conn, resp, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("falha ao conectar ao websocket Edge TTS (status %d): %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("falha ao conectar ao websocket Edge TTS: %w", err)
	}

	return conn, nil
}

// Synthesize gera o áudio diretamente em memória (buffer de bytes).
func (c *Client) Synthesize(ctx context.Context, opts SynthesizeOptions) ([]byte, error) {
	var buf bytes.Buffer
	if err := c.SynthesizeStream(ctx, opts, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SynthesizeStream gera o áudio completo e íntegro com streaming em tempo real.
func (c *Client) SynthesizeStream(ctx context.Context, opts SynthesizeOptions, writer io.Writer) error {
	conn, err := c.dialConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Garantir fechamento imediato se o contexto for cancelado
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	requestID := strings.ReplaceAll(uuid.New().String(), "-", "")
	formatInfo := GetFormatInfo(opts.Format)
	timestampStr := time.Now().UTC().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (Coordinated Universal Time)")

	// 1. Enviar Speech Config
	configMessage := fmt.Sprintf(
		"X-Timestamp:%s\r\n"+
			"Content-Type:application/json; charset=utf-8\r\n"+
			"Path:speech.config\r\n\r\n"+
			`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":"false","wordBoundaryEnabled":"false"},"outputFormat":"%s"}}}}`,
		timestampStr, formatInfo.EdgeFormat,
	)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(configMessage)); err != nil {
		return fmt.Errorf("falha ao enviar speech.config: %w", err)
	}

	// Detectar automaticamente o idioma correto a partir da voz
	lang := opts.Language
	if lang == "" || lang == "pt-BR" {
		lang = extractLangFromVoice(opts.Voice)
	}

	// 2. Enviar SSML Payload Completo com suporte a pausas inteligentes e SSML avançado
	ssmlContent := BuildSSMLWithOptions(opts.Text, SSMLOptions{
		Voice:       opts.Voice,
		Rate:        opts.Rate,
		Pitch:       opts.Pitch,
		Volume:      opts.Volume,
		Lang:        lang,
		BreakComma:  opts.BreakComma,
		BreakPeriod: opts.BreakPeriod,
	})
	slog.Info("Enviando SSML para Edge TTS", "ssml", ssmlContent)
	ssmlMessage := fmt.Sprintf(
		"X-RequestId:%s\r\n"+
			"Content-Type:application/ssml+xml\r\n"+
			"X-Timestamp:%sZ\r\n"+
			"Path:ssml\r\n\r\n"+
			"%s",
		requestID, timestampStr, ssmlContent,
	)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(ssmlMessage)); err != nil {
		return fmt.Errorf("falha ao enviar payload ssml: %w", err)
	}

	// 3. Receber os pacotes de áudio e fazer streaming imediato para o writer
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msgType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) || ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("erro ao receber mensagens do websocket: %w", err)
		}

		switch msgType {
		case websocket.BinaryMessage:
			if len(message) < 2 {
				continue
			}
			headerLen := int(binary.BigEndian.Uint16(message[0:2]))
			if len(message) < 2+headerLen {
				continue
			}

			headerStr := string(message[2 : 2+headerLen])
			if strings.Contains(headerStr, "Path:audio") {
				audioChunk := message[2+headerLen:]
				if len(audioChunk) > 0 {
					if _, writeErr := writer.Write(audioChunk); writeErr != nil {
						return fmt.Errorf("erro ao gravar chunk de áudio: %w", writeErr)
					}
				}
			}

		case websocket.TextMessage:
			textMsg := string(message)
			if strings.Contains(textMsg, "Path:turn.end") {
				slog.Debug("Síntese concluída com sucesso", "request_id", requestID)
				return nil
			}
		}
	}
}
