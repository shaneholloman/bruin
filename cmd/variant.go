package cmd

import (
	"context"

	"github.com/bruin-data/bruin/pkg/pipeline"
)

// variantInjectingBuilder wraps a pipeline builder so every CreatePipelineFromPath
// call has WithVariant(variantName) appended. It exists for the lint command's
// single-variant path: lint.Linter calls the builder once per discovered pipeline
// and we need each call to materialize the requested variant. For non-variant
// pipelines the builder no-ops on WithVariant, so this works for mixed projects.
type variantInjectingBuilder struct {
	inner interface {
		CreatePipelineFromPath(ctx context.Context, path string, opts ...pipeline.CreatePipelineOption) (*pipeline.Pipeline, error)
	}
	variantName string
}

func (b *variantInjectingBuilder) CreatePipelineFromPath(ctx context.Context, path string, opts ...pipeline.CreatePipelineOption) (*pipeline.Pipeline, error) {
	return b.inner.CreatePipelineFromPath(ctx, path, append(opts, pipeline.WithVariant(b.variantName))...)
}
