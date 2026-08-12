package memory

import (
	"context"
	"errors"
	"sync"
)

// Erros comuns de repositório
var (
	ErrNotFound      = errors.New("item não encontrado")
	ErrAlreadyExists = errors.New("item já existe")
)

// Entity representa o modelo de domínio armazenado
type Entity struct {
	ID   string
	Name string
}

// MemoryRepository implementa a interface de repositório usando memória/mapa thread-safe
type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]Entity
}

// NewMemoryRepository inicializa um novo repositório em memória
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make(map[string]Entity),
	}
}

// Save adiciona ou atualiza um item no repositório
func (r *MemoryRepository) Save(ctx context.Context, item Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.ID] = item
	return nil
}

// FindByID busca um item pelo ID
func (r *MemoryRepository) FindByID(ctx context.Context, id string) (Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return Entity{}, ErrNotFound
	}

	return item, nil
}

// FindAll retorna todos os itens armazenados
func (r *MemoryRepository) FindAll(ctx context.Context) ([]Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Entity, 0, len(r.items))
	for _, item := range r.items {
		list = append(list, item)
	}

	return list, nil
}

// Delete remove um item pelo ID
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return ErrNotFound
	}

	delete(r.items, id)
	return nil
}
