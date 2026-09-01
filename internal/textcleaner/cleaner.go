package textcleaner

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	reMarkdownImages   = regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`)
	reMarkdownLinks    = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	reMarkdownFormat   = regexp.MustCompile(`(\*\*|__|\*|_|~~)`)
	reCodeBlocks       = regexp.MustCompile("(?s)```.*?```")
	reInlineCode       = regexp.MustCompile("`([^`]+)`")
	reMarkdownHeaders  = regexp.MustCompile(`(?m)^\s*#{1,6}\s+`)
	reHTMLTags         = regexp.MustCompile(`</?[^>]+(>|$)`)
	reMultipleNewlines = regexp.MustCompile(`\n{2,}`)
	reMultipleSpaces   = regexp.MustCompile(`[ \t]{2,}`)
)

// CleanText remove Markdown, tags HTML, emojis e caracteres indesejados para sintetização TTS.
func CleanText(text string) string {
	if text == "" {
		return ""
	}

	// 1. Remover Emojis
	text = removeEmojis(text)

	// 2. Remover formatações Markdown mantendo o texto legível
	// IMPORTANTE: Imagens devem ser tratadas antes de links para não deixar '!' órfão
	text = reMarkdownImages.ReplaceAllString(text, "$1")
	text = reMarkdownLinks.ReplaceAllString(text, "$1")
	text = reCodeBlocks.ReplaceAllString(text, "")
	text = reInlineCode.ReplaceAllString(text, "$1")
	text = reMarkdownFormat.ReplaceAllString(text, "")
	text = reMarkdownHeaders.ReplaceAllString(text, "")
	text = reHTMLTags.ReplaceAllString(text, "")

	// 3. Normalização de espaçamento
	text = reMultipleNewlines.ReplaceAllString(text, "\n\n")
	text = reMultipleSpaces.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// removeEmojis filtra runes que pertencem a faixas de emojis e símbolos especiais.
func removeEmojis(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if !isEmoji(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F: // Emoticons
		return true
	case r >= 0x1F300 && r <= 0x1F5FF: // Misc Symbols and Pictographs
		return true
	case r >= 0x1F680 && r <= 0x1F6FF: // Transport and Map
		return true
	case r >= 0x1F1E0 && r <= 0x1F1FF: // Regional indicator symbols (Flags)
		return true
	case r >= 0x2600 && r <= 0x26FF:   // Miscellaneous Symbols
		return true
	case r >= 0x2700 && r <= 0x27BF:   // Dingbats
		return true
	case r >= 0xFE00 && r <= 0xFE0F:   // Variation Selectors
		return true
	case r >= 0x1F900 && r <= 0x1F9FF: // Supplemental Symbols and Pictographs
		return true
	case r >= 0x1FA00 && r <= 0x1FA6F: // Chess Symbols, Symbols and Pictographs Extended-A
		return true
	case r >= 0x1FA70 && r <= 0x1FAFF: // Symbols and Pictographs Extended-B
		return true
	case r >= 0x1F000 && r <= 0x1F02F: // Mahjong Tiles
		return true
	case r >= 0x1F0A0 && r <= 0x1F0FF: // Playing Cards
		return true
	case r == 0x200D:                  // Zero Width Joiner
		return true
	default:
		return unicode.Is(unicode.So, r) && (r > 127)
	}
}
