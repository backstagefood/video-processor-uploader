package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func GetEnvVarOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func IsValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func SanitizeEmailForPath(email string) string {
	// Replace @ with _
	sanitized := strings.ReplaceAll(email, "@", "_")

	// You might want to add additional sanitization:
	// - Replace dots (.) if they cause issues
	sanitized = strings.ReplaceAll(sanitized, ".", "_")

	// Remove any other potentially problematic characters
	sanitized = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			return r
		}
		return '_' // Replace other special chars with underscore
	}, sanitized)

	return sanitized
}

func GetFileSize(file multipart.File) (int64, error) {
	// Salva a posição atual do cursor
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Move para o final do arquivo para obter o tamanho
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Retorna o cursor para a posição original
	_, err = file.Seek(currentPos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return size, nil
}
