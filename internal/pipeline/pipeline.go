package pipeline

import (
	"context"
	"fmt"
	"io"
)

type StageConfig map[string]string

type StageSpec struct {
	Name   string
	Config StageConfig
}

type Stage interface {
	Name() string
	Wrap(ctx context.Context, src io.ReadCloser, cfg StageConfig) (io.ReadCloser, error)
}

type Registry struct {
	stages map[string]Stage
}

func NewRegistry(stages ...Stage) *Registry {
	r := &Registry{stages: make(map[string]Stage, len(stages))}
	for _, stage := range stages {
		r.Register(stage)
	}
	return r
}

func (r *Registry) Register(stage Stage) {
	r.stages[stage.Name()] = stage
}

func (r *Registry) Build(ctx context.Context, src io.ReadCloser, specs []StageSpec) (io.ReadCloser, error) {
	current := src
	for _, spec := range specs {
		stage, ok := r.stages[spec.Name]
		if !ok {
			current.Close()
			return nil, fmt.Errorf("unknown pipeline stage %q", spec.Name)
		}

		next, err := stage.Wrap(ctx, current, spec.Config)
		if err != nil {
			current.Close()
			return nil, fmt.Errorf("build stage %q: %w", spec.Name, err)
		}
		current = next
	}

	return current, nil
}
