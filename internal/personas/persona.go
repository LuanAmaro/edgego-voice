package personas

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	ErrPersonaNotFound = errors.New("persona não encontrada")
	ErrPersonaExists   = errors.New("já existe uma persona com este identificador")
	reSlug             = regexp.MustCompile(`^[a-z0-9-_]+$`)
)

// Persona define as configurações de uma persona de voz dedicada.
type Persona struct {
	ID           string    `json:"id"`            // Slug único (ex: "atendimento-whatsapp")
	Name         string    `json:"name"`          // Nome legível (ex: "Atendimento WhatsApp")
	Description  string    `json:"description"`   // Descrição da finalidade
	Voice        string    `json:"voice"`         // Voz neural (ex: "pt-BR-ThalitaMultilingualNeural")
	Format       string    `json:"format"`        // Formato (mp3, opus, wav, etc.)
	Speed        float64   `json:"speed"`         // Multiplicador de velocidade (0.25 a 2.0)
	Pitch        string    `json:"pitch"`         // Ajuste de tom (ex: "+0Hz")
	RemoveFilter bool      `json:"remove_filter"` // Se true, não sanitiza markdown/emojis
	APIKey       string    `json:"api_key"`       // Chave exclusiva da persona (opcional)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Manager gerencia o ciclo de vida e persistência das personas de áudio.
type Manager struct {
	mu       sync.RWMutex
	filePath string
	personas map[string]Persona
}

// NewManager cria e inicializa o gerenciador de personas a partir de um arquivo JSON.
func NewManager(filePath string) *Manager {
	if filePath == "" {
		filePath = "personas.json"
	}

	m := &Manager{
		filePath: filePath,
		personas: make(map[string]Persona),
	}

	if err := m.load(); err != nil {
		slog.Warn("Não foi possível carregar personas.json, criando padrões iniciais", "error", err)
		m.seedDefaults()
		_ = m.save()
	}

	return m
}

func (m *Manager) seedDefaults() {
	now := time.Now().UTC()
	defaults := []Persona{
		{
			ID:           "whatsapp-suporte",
			Name:         "WhatsApp Suporte (Opus)",
			Description:  "Otimizado para envio de mensagens de voz no WhatsApp com baixa latência e formato Opus.",
			Voice:        "pt-BR-ThalitaMultilingualNeural",
			Format:       "opus",
			Speed:        1.0,
			Pitch:        "+0Hz",
			RemoveFilter: false,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "assistente-vendas",
			Name:         "Assistente de Vendas (MP3)",
			Description:  "Voz enérgica e clara para atendimento comercial e chatbots web.",
			Voice:        "pt-BR-AntonioNeural",
			Format:       "mp3",
			Speed:        1.05,
			Pitch:        "+0Hz",
			RemoveFilter: false,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "podcast-narrador",
			Name:         "Narrador de Alta Qualidade (WAV)",
			Description:  "Áudio PCM sem perdas para vídeos, podcasts e narrações profissionais.",
			Voice:        "pt-BR-FranciscaNeural",
			Format:       "wav",
			Speed:        0.95,
			Pitch:        "+0Hz",
			RemoveFilter: false,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	for _, p := range defaults {
		m.personas[p.ID] = p
	}
}

func (m *Manager) load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	var list []Persona
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	m.personas = make(map[string]Persona, len(list))
	for _, p := range list {
		m.personas[p.ID] = p
	}

	slog.Info("Personas carregadas com sucesso", "total", len(m.personas))
	return nil
}

func (m *Manager) save() error {
	list := make([]Persona, 0, len(m.personas))
	for _, p := range m.personas {
		list = append(list, p)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("falha ao serializar personas: %w", err)
	}

	if err := os.WriteFile(m.filePath, data, 0644); err != nil {
		return fmt.Errorf("falha ao salvar personas.json: %w", err)
	}

	return nil
}

// List retorna todas as personas cadastradas.
func (m *Manager) List() []Persona {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]Persona, 0, len(m.personas))
	for _, p := range m.personas {
		list = append(list, p)
	}
	return list
}

// Get obtém uma persona pelo ID.
func (m *Manager) Get(id string) (Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.personas[id]
	if !ok {
		return Persona{}, ErrPersonaNotFound
	}
	return p, nil
}

// Slugify normaliza strings para identificadores seguros (sem acentos, espaços viram hífens, sem caracteres especiais).
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "õ", "o", "ô", "o", "ö", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u",
		"ç", "c", "ñ", "n",
	)
	s = replacer.Replace(s)

	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		} else if r == ' ' || r == '/' || r == '.' || r == ':' {
			b.WriteRune('-')
		}
	}
	res := regexp.MustCompile(`-+`).ReplaceAllString(b.String(), "-")
	return strings.Trim(res, "-_")
}

// Create cadastra uma nova persona.
func (m *Manager) Create(p Persona) (Persona, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Se ID estiver vazio, gerar a partir do Nome
	if strings.TrimSpace(p.ID) == "" && strings.TrimSpace(p.Name) != "" {
		p.ID = Slugify(p.Name)
	} else {
		p.ID = Slugify(p.ID)
	}

	if p.ID == "" {
		return Persona{}, errors.New("o identificador (ID) da persona é obrigatório")
	}

	if _, exists := m.personas[p.ID]; exists {
		return Persona{}, ErrPersonaExists
	}

	if p.Name == "" {
		p.Name = p.ID
	}
	if p.Voice == "" {
		p.Voice = "pt-BR-ThalitaMultilingualNeural"
	}
	if p.Format == "" {
		p.Format = "mp3"
	}
	if p.Speed <= 0 {
		p.Speed = 1.0
	}
	if p.Pitch == "" {
		p.Pitch = "+0Hz"
	}

	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	m.personas[p.ID] = p
	if err := m.save(); err != nil {
		delete(m.personas, p.ID)
		return Persona{}, err
	}

	slog.Info("Nova persona criada com sucesso", "id", p.ID, "voice", p.Voice, "format", p.Format)
	return p, nil
}

// Update atualiza os dados de uma persona existente.
func (m *Manager) Update(id string, update Persona) (Persona, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id = strings.ToLower(strings.TrimSpace(id))
	p, ok := m.personas[id]
	if !ok {
		return Persona{}, ErrPersonaNotFound
	}

	if update.Name != "" {
		p.Name = update.Name
	}
	p.Description = update.Description
	if update.Voice != "" {
		p.Voice = update.Voice
	}
	if update.Format != "" {
		p.Format = update.Format
	}
	if update.Speed > 0 {
		p.Speed = update.Speed
	}
	if update.Pitch != "" {
		p.Pitch = update.Pitch
	}
	p.RemoveFilter = update.RemoveFilter
	p.APIKey = update.APIKey
	p.UpdatedAt = time.Now().UTC()

	m.personas[id] = p
	if err := m.save(); err != nil {
		return Persona{}, err
	}

	slog.Info("Persona atualizada", "id", id)
	return p, nil
}

// Delete remove uma persona.
func (m *Manager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id = strings.ToLower(strings.TrimSpace(id))
	if _, ok := m.personas[id]; !ok {
		return ErrPersonaNotFound
	}

	delete(m.personas, id)
	return m.save()
}
