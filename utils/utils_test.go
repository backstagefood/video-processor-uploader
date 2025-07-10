package utils

import (
	"bytes"
	"mime/multipart"
	"os"
	"strings"
	"testing"
)

func TestGetEnvVarOrDefault(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		defaultValue   string
		setEnv         bool
		envValue       string
		expectedResult string
	}{
		{
			name:           "Env var exists",
			key:            "TEST_ENV",
			defaultValue:   "default",
			setEnv:         true,
			envValue:       "actual",
			expectedResult: "actual",
		},
		{
			name:           "Env var does not exist",
			key:            "NON_EXISTENT_ENV",
			defaultValue:   "default",
			setEnv:         false,
			expectedResult: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := GetEnvVarOrDefault(tt.key, tt.defaultValue)
			if result != tt.expectedResult {
				t.Errorf("Expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestIsValidVideoFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"Valid MP4", "video.mp4", true},
		{"Valid AVI", "video.avi", true},
		{"Valid MOV", "video.mov", true},
		{"Valid MKV", "video.mkv", true},
		{"Valid WMV", "video.wmv", true},
		{"Valid FLV", "video.flv", true},
		{"Valid WEBM", "video.webm", true},
		{"Invalid TXT", "document.txt", false},
		{"Invalid PDF", "document.pdf", false},
		{"Invalid no extension", "video", false},
		{"Uppercase extension", "video.MP4", true},
		{"Mixed case extension", "video.Mp4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVideoFile(tt.filename)
			if result != tt.expected {
				t.Errorf("For %s, expected %v, got %v", tt.filename, tt.expected, result)
			}
		})
	}
}

func TestSanitizeEmailForPath(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "Simple email",
			email:    "user@example.com",
			expected: "user_example_com",
		},
		{
			name:     "Email with dots",
			email:    "user.name@sub.example.com",
			expected: "user_name_sub_example_com",
		},
		{
			name:     "Email with special chars",
			email:    "user+filter@example.com",
			expected: "user_filter_example_com",
		},
		{
			name:     "Email with numbers",
			email:    "user123@example.com",
			expected: "user123_example_com",
		},
		{
			name:     "Email with mixed case",
			email:    "User.Name@Example.COM",
			expected: "User_Name_Example_COM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeEmailForPath(tt.email)
			if result != tt.expected {
				t.Errorf("For %s, expected %s, got %s", tt.email, tt.expected, result)
			}
		})
	}
}

func TestGetFileSize(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected int64
	}{
		{"Empty file", []byte(""), 0},
		{"Small file", []byte("test content"), int64(len("test content"))},
		{"Medium file", bytes.Repeat([]byte("a"), 1024), 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a multipart file from the content
			file := &bytes.Buffer{}
			writer := multipart.NewWriter(file)
			part, err := writer.CreateFormFile("file", "test.txt")
			if err != nil {
				t.Fatal(err)
			}
			_, err = part.Write(tt.content)
			if err != nil {
				t.Fatal(err)
			}
			writer.Close()

			reader := bytes.NewReader(file.Bytes())
			multipartReader, err := multipart.NewReader(reader, writer.Boundary()).ReadForm(0)
			if err != nil {
				t.Fatal(err)
			}

			fileHeader := multipartReader.File["file"][0]
			fileToTest, err := fileHeader.Open()
			if err != nil {
				t.Fatal(err)
			}
			defer fileToTest.Close()

			size, err := GetFileSize(fileToTest)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if size != tt.expected {
				t.Errorf("Expected size %d, got %d", tt.expected, size)
			}
		})
	}
}

func TestGenerateUniqueKey(t *testing.T) {
	// Test that generated keys are unique
	keys := make(map[string]bool)
	const iterations = 1000

	for i := 0; i < iterations; i++ {
		key := GenerateUniqueKey()
		if keys[key] {
			t.Fatalf("Duplicate key generated: %s", key)
		}
		keys[key] = true

		// Basic format check: timestamp_random
		parts := strings.Split(key, "_")
		if len(parts) < 3 { // 20060102, 150405.000000000, random
			t.Errorf("Key has invalid format: %s", key)
		}
	}
}
