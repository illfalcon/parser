package str

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"
)

func ChunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}

func FindSha1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := fmt.Sprintf("%x", h.Sum(nil))
	return bs
}

func NormalizeNewLines(text string) string {
	nl := regexp.MustCompile(`\n+`)
	text = nl.ReplaceAllString(text, "\n")
	return text
}

func NormalizeSpaces(text string) string {
	spaces := regexp.MustCompile(`\s+`)
	text = spaces.ReplaceAllString(text, " ")
	return text
}

func NormalizeOnlySpaces(text string) string {
	onlySpaces := regexp.MustCompile(`[\t ]+`)
	text = onlySpaces.ReplaceAllString(text, " ")
	return text
}

func NormalizePunctuation(text string) string {
	text = strings.ReplaceAll(text, ".", ". ")
	text = strings.ReplaceAll(text, ",", ", ")
	text = strings.ReplaceAll(text, "!", "! ")
	text = strings.ReplaceAll(text, "?", "? ")
	text = strings.ReplaceAll(text, ";", "; ")
	text = NormalizeOnlySpaces(text)
	return text
}
