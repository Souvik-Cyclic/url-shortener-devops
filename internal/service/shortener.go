package service

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ShortenerService struct {
	codeToURL map[string]string
	urlToCode map[string]string
	mu        sync.RWMutex
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (s *ShortenerService) Shorten(originalURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL = strings.TrimSpace(originalURL)

	// Check if URL already exists
	if code, exists := s.urlToCode[originalURL]; exists {
		return code
	}

	code := generateCode()
	// Ensure uniqueness
	for _, exists := s.codeToURL[code]; exists; _, exists = s.codeToURL[code] {
		code = generateCode()
	}

	s.codeToURL[code] = originalURL
	s.urlToCode[originalURL] = code
	return code
}

func (s *ShortenerService) GetOriginalURL(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.codeToURL[code]
	return url, exists
}
