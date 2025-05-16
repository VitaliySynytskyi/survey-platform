package auth

import (
	"sync"
	"time"
)

// TokenBlacklist представляє механізм для блокування використаних або скомпрометованих токенів
type TokenBlacklist struct {
	blacklist map[string]time.Time // Ключ - токен, значення - час, коли токен можна видалити з blacklist
	mutex     sync.RWMutex
}

// NewTokenBlacklist створює новий екземпляр TokenBlacklist і запускає горутину для очистки
func NewTokenBlacklist() *TokenBlacklist {
	bl := &TokenBlacklist{
		blacklist: make(map[string]time.Time),
	}

	// Запускаємо горутину для очистки застарілих токенів
	go bl.cleanupRoutine()

	return bl
}

// Add додає токен до чорного списку на заданий період часу
func (bl *TokenBlacklist) Add(token string, ttl time.Duration) {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	// Час, коли токен можна видалити з blacklist (поточний час + TTL)
	expiry := time.Now().Add(ttl)
	bl.blacklist[token] = expiry
}

// Contains перевіряє, чи міститься токен у чорному списку
func (bl *TokenBlacklist) Contains(token string) bool {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	expiry, exists := bl.blacklist[token]

	// Якщо токен є і не прострочений, повертаємо true
	if exists && time.Now().Before(expiry) {
		return true
	}

	// Якщо токен прострочений, видаляємо його при перевірці
	// (це додаткова очистка до основної горутини)
	if exists {
		bl.mutex.RUnlock() // Розблоковуємо читання перед блокуванням запису
		bl.mutex.Lock()
		delete(bl.blacklist, token)
		bl.mutex.Unlock()
		bl.mutex.RLock() // Відновлюємо блокування для читання
	}

	return false
}

// cleanupRoutine регулярно очищає прострочені токени з чорного списку
func (bl *TokenBlacklist) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		bl.cleanup()
	}
}

// cleanup видаляє прострочені токени з чорного списку
func (bl *TokenBlacklist) cleanup() {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	now := time.Now()
	for token, expiry := range bl.blacklist {
		if now.After(expiry) {
			delete(bl.blacklist, token)
		}
	}
}
