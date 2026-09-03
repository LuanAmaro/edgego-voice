package textcleaner

import (
	"testing"
)

func TestCleanText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Remover Emojis",
			input:    "Olá mundo! 👋🚀 Como você está? 😊",
			expected: "Olá mundo! Como você está?",
		},
		{
			name:     "Remover Links Markdown",
			input:    "Acesse o [Google](https://google.com) para mais informações.",
			expected: "Acesse o Google para mais informações.",
		},
		{
			name:     "Remover Imagens Markdown",
			input:    "Veja esta imagem: ![Logotipo da Empresa](logo.png)",
			expected: "Veja esta imagem: Logotipo da Empresa",
		},
		{
			name:     "Remover Formatação Markdown (Negrito, Itálico)",
			input:    "Este texto tem **negrito**, *itálico*, __sublinhado__ e ~~riscado~~.",
			expected: "Este texto tem negrito, itálico, sublinhado e riscado.",
		},
		{
			name:     "Remover Blocos de Código",
			input:    "Veja o código abaixo:\n```go\nfunc main() {}\n```\nFim do código.",
			expected: "Veja o código abaixo:\n\nFim do código.",
		},
		{
			name:     "Remover Código Inline",
			input:    "Execute o comando `npm install` no terminal.",
			expected: "Execute o comando npm install no terminal.",
		},
		{
			name:     "Remover Headers Markdown",
			input:    "# Título Principal\n## Subtítulo\nTexto normal.",
			expected: "Título Principal\nSubtítulo\nTexto normal.",
		},
		{
			name:     "Remover Tags HTML",
			input:    "Olá <b>mundo</b>, <span class=\"bold\">bem-vindo</span>!",
			expected: "Olá mundo, bem-vindo!",
		},
		{
			name:     "Preservar Tags SSML Legítimas",
			input:    "Atenção! <break time=\"400ms\"/> Este é um teste com <emphasis level=\"strong\">ênfase</emphasis>.",
			expected: "Atenção! <break time=\"400ms\"/> Este é um teste com <emphasis level=\"strong\">ênfase</emphasis>.",
		},
		{
			name:     "String Vazia",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanText(tt.input)
			if got != tt.expected {
				t.Errorf("CleanText(%q) = %q, esperado %q", tt.input, got, tt.expected)
			}
		})
	}
}
