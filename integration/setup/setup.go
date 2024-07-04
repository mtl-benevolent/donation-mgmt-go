package setup

import (
	"context"
	"testing"
)

type Builder interface {
	Name() string
	Type() string
	Execute(ctx context.Context) (any, error)
}

type Setup struct {
	builders []Builder
}

type SetupResult struct {
	entities map[string]any
}

func NewSetup() *Setup {
	return &Setup{
		builders: make([]Builder, 0),
	}
}

func (s *Setup) Execute(ctx context.Context, t *testing.T) *SetupResult {
	results := &SetupResult{
		entities: make(map[string]any),
	}

	for _, builder := range s.builders {
		entity, err := builder.Execute(ctx)
		if err != nil {
			t.Errorf("Failed to execute setup. Unable to build %s (name: %s): %v", builder.Type(), builder.Name(), err)
			t.FailNow()
			return results
		}

		results.entities[builder.Name()] = entity
	}

	return results
}

func (s *SetupResult) Get(name string) (any, bool) {
	entity, ok := s.entities[name]
	if !ok {
		return nil, false
	}

	return entity, true
}

func GetEntity[T any](results *SetupResult, name string) (T, bool) {
	var def T

	entity, ok := results.Get(name)
	if !ok {
		return def, false
	}

	t, ok := entity.(T)
	return t, ok
}
