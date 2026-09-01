package audiocache

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	key       string
	data      []byte
	size      int64
	expiresAt time.Time
}

// Cache é um cache em memória LRU (Least Recently Used) thread-safe para áudios sintetizados.
type Cache struct {
	mu        sync.RWMutex
	items     map[string]*list.Element
	evictList *list.List
	maxBytes  int64
	curBytes  int64
	ttl       time.Duration
	enabled   bool
}

// NewCache cria uma nova instância de cache de áudio.
func NewCache(enabled bool, maxMB int, ttl time.Duration) *Cache {
	return &Cache{
		items:     make(map[string]*list.Element),
		evictList: list.New(),
		maxBytes:  int64(maxMB) * 1024 * 1024,
		ttl:       ttl,
		enabled:   enabled,
	}
}

// GenerateKey gera um hash SHA256 único a partir dos parâmetros de síntese.
func GenerateKey(text, voice, format string, speed float64, pitch string, removeFilter bool) string {
	raw := fmt.Sprintf("%s|%s|%s|%.2f|%s|%t", text, voice, format, speed, pitch, removeFilter)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// Get busca um áudio no cache. Retorna os bytes e se foi encontrado (HIT).
func (c *Cache) Get(key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	elem, found := c.items[key]
	if !found {
		return nil, false
	}

	item := elem.Value.(*cacheItem)
	if time.Now().After(item.expiresAt) {
		c.removeElement(elem)
		return nil, false
	}

	c.evictList.MoveToFront(elem)
	return item.data, true
}

// Set armazena um áudio sintetizado no cache LRU.
func (c *Cache) Set(key string, data []byte) {
	if !c.enabled || len(data) == 0 {
		return
	}

	itemSize := int64(len(data))
	if itemSize > c.maxBytes {
		return // Não armazena itens maiores que a capacidade máxima total
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Se já existe, atualiza
	if elem, found := c.items[key]; found {
		oldItem := elem.Value.(*cacheItem)
		c.curBytes -= oldItem.size
		oldItem.data = data
		oldItem.size = itemSize
		oldItem.expiresAt = time.Now().Add(c.ttl)
		c.curBytes += itemSize
		c.evictList.MoveToFront(elem)
		c.evictIfNeeded()
		return
	}

	// Novo item
	item := &cacheItem{
		key:       key,
		data:      data,
		size:      itemSize,
		expiresAt: time.Now().Add(c.ttl),
	}

	elem := c.evictList.PushFront(item)
	c.items[key] = elem
	c.curBytes += itemSize
	c.evictIfNeeded()
}

// removeElement remove um item da lista e do mapa.
func (c *Cache) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	item := elem.Value.(*cacheItem)
	delete(c.items, item.key)
	c.curBytes -= item.size
}

// evictIfNeeded remove itens antigos caso exceda o limite de memória.
func (c *Cache) evictIfNeeded() {
	for c.curBytes > c.maxBytes && c.evictList.Len() > 0 {
		oldest := c.evictList.Back()
		if oldest != nil {
			c.removeElement(oldest)
		}
	}
}

// Stats retorna estatísticas do cache (itens e bytes ocupados).
func (c *Cache) Stats() (int, int64, int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items), c.curBytes, c.maxBytes
}
