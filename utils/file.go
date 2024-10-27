package utils

import (
	"path/filepath"
	"strings"
)

func DetermineDocumentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".pdf":
		return "PDF"
	case ".doc", ".docx":
		return "DOC"
	case ".mp4", ".avi", ".mkv":
		return "VIDEO"
	default:
		return "PDF"
	}
}
