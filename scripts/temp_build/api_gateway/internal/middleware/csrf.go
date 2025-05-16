package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"
)

var (
	ErrInvalidCSRFToken = errors.New("invalid CSRF token")
	ErrMissingCSRFToken = errors.New("missing CSRF token")
)

// CSRFConfig contains configuration for CSRF middleware
type CSRFConfig struct {
	TokenLength    int           // Довжина згенерованого токена
	CookieName     string        // Ім'я cookie для зберігання токена
	HeaderName     string        // Ім'я заголовка для передачі токена в запитах
	CookieMaxAge   int           // Час життя cookie в секундах
	TokenExpiry    time.Duration // Час життя токена
	CookieHTTPOnly bool          // Чи доступний cookie лише через HTTP
	CookieSecure   bool          // Чи передавати cookie лише через HTTPS
	CookiePath     string        // Шлях cookie (зазвичай "/")
}

// DefaultCSRFConfig повертає налаштування за замовчуванням
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLength:    32,
		CookieName:     "csrf_token",
		HeaderName:     "X-CSRF-Token",
		CookieMaxAge:   86400,          // 24 години
		TokenExpiry:    24 * time.Hour, // 24 години
		CookieHTTPOnly: true,
		CookieSecure:   true, // В продакшені має бути true
		CookiePath:     "/",
	}
}

// CSRFMiddleware представляє middleware для захисту від CSRF атак
type CSRFMiddleware struct {
	config   CSRFConfig
	tokens   map[string]time.Time // Зберігає токени та час їх закінчення
	mutex    sync.RWMutex
	cleanupC chan struct{} // Канал для зупинки горутини очистки
}

// NewCSRFMiddleware створює новий екземпляр CSRFMiddleware
func NewCSRFMiddleware(config CSRFConfig) *CSRFMiddleware {
	m := &CSRFMiddleware{
		config:   config,
		tokens:   make(map[string]time.Time),
		cleanupC: make(chan struct{}),
	}

	// Запускаємо горутину для очистки прострочених токенів
	go m.cleanupRoutine()

	return m
}

// генерує випадковий CSRF токен
func (m *CSRFMiddleware) generateToken() (string, error) {
	bytes := make([]byte, m.config.TokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// валідує токен
func (m *CSRFMiddleware) validateToken(token string) bool {
	if token == "" {
		return false
	}

	m.mutex.RLock()
	expiry, exists := m.tokens[token]
	m.mutex.RUnlock()

	// Перевіряємо, чи токен існує та не є простроченим
	return exists && time.Now().Before(expiry)
}

// Middleware повертає HTTP middleware для захисту від CSRF атак
func (m *CSRFMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Для методів, які не змінюють стан (GET, HEAD, OPTIONS, TRACE),
		// не перевіряємо CSRF токен
		if r.Method == http.MethodGet ||
			r.Method == http.MethodHead ||
			r.Method == http.MethodOptions ||
			r.Method == http.MethodTrace {
			next.ServeHTTP(w, r)
			return
		}

		// Для CORS preflight запитів також пропускаємо перевірку
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			next.ServeHTTP(w, r)
			return
		}

		// Перевіряємо CSRF токен
		csrfToken := r.Header.Get(m.config.HeaderName)
		if csrfToken == "" {
			http.Error(w, ErrMissingCSRFToken.Error(), http.StatusForbidden)
			return
		}

		if !m.validateToken(csrfToken) {
			http.Error(w, ErrInvalidCSRFToken.Error(), http.StatusForbidden)
			return
		}

		// Токен валідний, продовжуємо обробку запиту
		next.ServeHTTP(w, r)
	})
}

// IssueToken генерує новий CSRF токен і встановлює його в cookie
func (m *CSRFMiddleware) IssueToken(w http.ResponseWriter) (string, error) {
	token, err := m.generateToken()
	if err != nil {
		return "", err
	}

	// Зберігаємо токен з часом закінчення
	m.mutex.Lock()
	m.tokens[token] = time.Now().Add(m.config.TokenExpiry)
	m.mutex.Unlock()

	// Встановлюємо токен в cookie
	cookie := &http.Cookie{
		Name:     m.config.CookieName,
		Value:    token,
		MaxAge:   m.config.CookieMaxAge,
		Path:     m.config.CookiePath,
		HttpOnly: m.config.CookieHTTPOnly,
		Secure:   m.config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return token, nil
}

// GetTokenFromCookie отримує CSRF токен з cookie
func (m *CSRFMiddleware) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(m.config.CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// cleanupRoutine регулярно очищає прострочені токени
func (m *CSRFMiddleware) cleanupRoutine() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.cleanupC:
			return
		}
	}
}

// cleanup видаляє прострочені токени
func (m *CSRFMiddleware) cleanup() {
	now := time.Now()
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for token, expiry := range m.tokens {
		if now.After(expiry) {
			delete(m.tokens, token)
		}
	}
}

// Stop зупиняє горутину очистки
func (m *CSRFMiddleware) Stop() {
	close(m.cleanupC)
}
