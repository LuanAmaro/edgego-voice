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
	reHTMLTags         = regexp.MustCompile(`(?i)</?([a-zA-Z0-9_:-]+)(?:\s+[^>]*)?/?>`)
	reMultipleNewlines = regexp.MustCompile(`\n{2,}`)
	reMultipleSpaces   = regexp.MustCompile(`[ \t]{2,}`)

	ssmlAllowedTags = map[string]struct{}{
		"break":            {},
		"say-as":           {},
		"emphasis":         {},
		"sub":              {},
		"phoneme":          {},
		"prosody":          {},
		"voice":            {},
		"speak":            {},
		"s":                {},
		"p":                {},
		"lang":             {},
		"audio":            {},
		"mstts:express-as": {},
		"express-as":       {},
		"sfx":              {},
	}
)

// CleanText remove Markdown, tags HTML não-SSML, emojis e caracteres indesejados para sintetização TTS.
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

	// 3. Remover tags HTML não-SSML mantendo tags SSML legítimas
	text = removeNonSSMLTags(text)

	// 4. Normalização de espaçamento
	text = reMultipleNewlines.ReplaceAllString(text, "\n\n")
	text = reMultipleSpaces.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// removeNonSSMLTags remove tags como <div>, <span>, <script>, mantendo <break>, <emphasis>, <say-as>, etc.
func removeNonSSMLTags(text string) string {
	return reHTMLTags.ReplaceAllStringFunc(text, func(match string) string {
		submatches := reHTMLTags.FindStringSubmatch(match)
		if len(submatches) > 1 {
			tagName := strings.ToLower(submatches[1])
			if _, ok := ssmlAllowedTags[tagName]; ok {
				return match
			}
		}
		return ""
	})
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
