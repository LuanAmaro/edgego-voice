package audiocache

import (
	"sync"
)

// call representa uma chamada de síntese em andamento.
type call[T any] struct {
	wg     sync.WaitGroup
	val    T
	err    error
	shared bool
}

// Group deduplica execuções concorrentes com a mesma chave.
// Se 50 chamadas solicitarem o mesmo áudio ao mesmo tempo, apenas 1 síntese é executada,
// e todas as 50 goroutines recebem o mesmo resultado sem duplicar memória ou requisições ao Edge TTS.
type Group[T any] struct {
	mu sync.Mutex
	m  map[string]*call[T]
}

// NewGroup cria uma nova instância de Group para deduplicação.
func NewGroup[T any]() *Group[T] {
	return &Group[T]{
		m: make(map[string]*call[T]),
	}
}

// Do executa e retorna o resultado da função fornecida, garantindo que apenas uma execução
// esteja em andamento para uma determinada chave ao mesmo tempo.
func (g *Group[T]) Do(key string, fn func() (T, error)) (T, error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[T])
	}

	if c, ok := g.m[key]; ok {
		c.shared = true
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}

	c := new(call[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	shared := c.shared
	g.mu.Unlock()

	return c.val, c.err, shared
}
