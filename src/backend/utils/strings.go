package utils

import (
	"log/slog"
	"unicode"
	"unicode/utf8"
)

func ToSentenceString(text string) string {
	if len(text) == 0 {
		return ""
	}

	rune, size := utf8.DecodeRuneInString(text)

	if rune == utf8.RuneError {
		slog.Error("Failed to decode the rune from a string. Returning unchanged string.", "string", text)
		return text
	}

	return string(unicode.ToUpper(rune)) + text[size:]
}