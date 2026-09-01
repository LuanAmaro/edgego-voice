package audiocache

import (
	"testing"
	"time"
)

func TestAudioCache(t *testing.T) {
	cache := NewCache(true, 1, 1*time.Minute) // 1MB

	key := GenerateKey("Olá mundo", "pt-BR-ThalitaMultilingualNeural", "mp3", 1.0, "+0Hz", false)
	data := []byte("audio-bytes-fake-12345")

	// 1. Get miss
	_, found := cache.Get(key)
	if found {
		t.Fatalf("Esperava cache miss, mas encontrou")
	}

	// 2. Set
	cache.Set(key, data)

	// 3. Get hit
	cachedData, found := cache.Get(key)
	if !found {
		t.Fatalf("Esperava cache hit")
	}
	if string(cachedData) != string(data) {
		t.Fatalf("Dados no cache divergentes")
	}

	// 4. Test eviction
	items, bytes, _ := cache.Stats()
	if items != 1 || bytes != int64(len(data)) {
		t.Fatalf("Stats inválidas: items=%d, bytes=%d", items, bytes)
	}
}
