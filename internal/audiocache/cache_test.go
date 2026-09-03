package audiocache

import (
	"fmt"
	"sync"
	"sync/atomic"
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

	// 4. Test stats
	items, bytes, _ := cache.Stats()
	if items != 1 || bytes != int64(len(data)) {
		t.Fatalf("Stats inválidas: items=%d, bytes=%d", items, bytes)
	}
}

func TestAudioCache_SLRUProtection(t *testing.T) {
	// Cria cache com 1MB total (~65KB por shard)
	cache := NewCache(true, 1, 1*time.Hour)

	frequentKey := GenerateKey("Saudação Frequente de Atendimento", "pt-BR-FranciscaNeural", "mp3", 1.0, "+0Hz", false)
	frequentData := []byte("audio-frequente-muito-importante")

	cache.Set(frequentKey, frequentData)
	// Segundo hit promove para Protected Segment
	cache.Get(frequentKey)

	// Simula enxurrada de itens únicos (scan pollution)
	shard := cache.getShard(frequentKey)
	dummySize := int(shard.maxBytes / 4)
	dummyData := make([]byte, dummySize)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("random-flood-key-%d", i)
		cache.Set(key, dummyData)
	}

	// O item frequente deve ter sido protegido pelo Segmented LRU!
	cachedData, found := cache.Get(frequentKey)
	if !found {
		t.Fatalf("Item frequente foi expulso indevidamente do cache pelo flood (Scan Pollution)")
	}
	if string(cachedData) != string(frequentData) {
		t.Fatalf("Conteúdo corrompido")
	}
}

func TestSingleflight(t *testing.T) {
	group := NewGroup[[]byte]()
	var executionCount int32

	var wg sync.WaitGroup
	numWorkers := 20
	key := "chamada-concorrente-telefonia"

	results := make([][]byte, numWorkers)
	sharedFlags := make([]bool, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			val, err, shared := group.Do(key, func() ([]byte, error) {
				atomic.AddInt32(&executionCount, 1)
				time.Sleep(50 * time.Millisecond) // Simula latência de rede do Edge TTS
				return []byte("audio-sintetizado-uma-vez-so"), nil
			})
			if err != nil {
				t.Errorf("worker %d erro: %v", workerID, err)
			}
			results[workerID] = val
			sharedFlags[workerID] = shared
		}(i)
	}

	wg.Wait()

	// Garante que mesmo com 20 requisições simultâneas, a função de síntese foi chamada APENAS 1 VEZ!
	if count := atomic.LoadInt32(&executionCount); count != 1 {
		t.Fatalf("Esperava exatamente 1 execução de síntese para 20 chamadas simultâneas, obteve %d", count)
	}

	for i := 0; i < numWorkers; i++ {
		if string(results[i]) != "audio-sintetizado-uma-vez-so" {
			t.Fatalf("Worker %d recebeu dados incorretos", i)
		}
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache(true, 10, 1*time.Hour)
	key := GenerateKey("Texto de benchmark", "pt-BR-FranciscaNeural", "mp3", 1.0, "+0Hz", false)
	cache.Set(key, make([]byte, 1024))

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cache.Get(key)
		}
	})
}

func BenchmarkCacheGetMiss(b *testing.B) {
	cache := NewCache(true, 10, 1*time.Hour)
	key := "chave-inexistente-para-medir-miss"

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cache.Get(key)
		}
	})
}

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache(true, 50, 1*time.Hour)
	data := make([]byte, 2048)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := fmt.Sprintf("bench-key-%d", i%1000)
			cache.Set(k, data)
			i++
		}
	})
}

func BenchmarkCacheMixed90Read10Write(b *testing.B) {
	cache := NewCache(true, 50, 1*time.Hour)
	data := make([]byte, 1024)

	// Pré-popular chaves frequentes
	for i := 0; i < 50; i++ {
		k := fmt.Sprintf("frequent-prompt-%d", i)
		cache.Set(k, data)
		cache.Get(k) // promove para protegido
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				cache.Set(fmt.Sprintf("write-key-%d", i%200), data)
			} else {
				_, _ = cache.Get(fmt.Sprintf("frequent-prompt-%d", i%50))
			}
			i++
		}
	})
}

func BenchmarkSingleflight(b *testing.B) {
	group := NewGroup[[]byte]()
	data := []byte("audio-sintetizado")

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = group.Do("mesma-chave-concorrente", func() ([]byte, error) {
				return data, nil
			})
		}
	})
}

func BenchmarkGenerateKey(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = GenerateKey("Olá! Seja bem-vindo ao suporte telefônico.", "pt-BR-FranciscaNeural", "mp3", 1.0, "+0Hz", false)
	}
}
