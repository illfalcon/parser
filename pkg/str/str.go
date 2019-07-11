package str

import (
	"crypto/sha1"
	"fmt"
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
