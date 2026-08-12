package memory_test

import (
	"context"
	"testing"

	"infrastructure/repository/memory"
)

func TestMemoryRepository_SaveAndFindByID(t *testing.T) {
	repo := memory.NewMemoryRepository()
	ctx := context.Background()

	entity := memory.Entity{ID: "1", Name: "Teste"}

	// Teste de salvamento
	err := repo.Save(ctx, entity)
	if err != nil {
		t.Fatalf("esperava erro nil ao salvar, obteve: %v", err)
	}

	// Teste de busca
	found, err := repo.FindByID(ctx, "1")
	if err != nil {
		t.Fatalf("esperava erro nil ao buscar, obteve: %v", err)
	}

	if found.Name != entity.Name {
		t.Errorf("esperava nome %s, obteve %s", entity.Name, found.Name)
	}
}
