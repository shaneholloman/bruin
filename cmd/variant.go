package cmd

import (
	"context"
	"fmt"

	"github.com/bruin-data/bruin/pkg/jinja"
	"github.com/bruin-data/bruin/pkg/pipeline"
)

// variantRendererFactory builds a RenderFunc using a Jinja renderer whose context
// exposes `var.X` for every variable plus `variant` for the variant name.
func variantRendererFactory(vars map[string]any, variantName string) pipeline.RenderFunc {
	r := jinja.NewRenderer(jinja.Context{
		"var":     vars,
		"variant": variantName,
	})
	return r.Render
}

// variantApplyVariablesMutator returns a pipeline mutator that merges the named
// variant's overrides into pl.Variables. Use it when --variant is set so that
// later mutators (e.g. asset parameter rendering) see the variant's values.
func variantApplyVariablesMutator(variantName string) pipeline.PipelineMutator {
	return func(_ context.Context, p *pipeline.Pipeline) (*pipeline.Pipeline, error) {
		if variantName == "" {
			return p, nil
		}
		if len(p.Variants) == 0 {
			return nil, fmt.Errorf("pipeline %q does not declare any variants but --variant=%q was provided", p.Name, variantName)
		}
		if err := p.ApplyVariantVariables(variantName); err != nil {
			return nil, err
		}
		return p, nil
	}
}

// dynamicVariantMutator returns a pipeline mutator whose target variant can be
// switched at call time by mutating *current. When *current is empty the mutator
// is a no-op; when *current is set and the pipeline declares variants, the
// variant's variables are merged so downstream asset mutators see them.
//
// This is the mechanism that lets `validate` iterate variants without registering
// a separate mutator per iteration.
func dynamicVariantMutator(current *string) pipeline.PipelineMutator {
	return func(_ context.Context, p *pipeline.Pipeline) (*pipeline.Pipeline, error) {
		name := *current
		if name == "" || len(p.Variants) == 0 {
			return p, nil
		}
		if err := p.ApplyVariantVariables(name); err != nil {
			return nil, err
		}
		return p, nil
	}
}

// projectHasVariantPipelines probes the discovered pipelines to determine whether
// any of them declare a variants block. The dynamic variant pointer is forced to
// "" while probing so dynamicVariantMutator stays a no-op even if the caller has
// already installed it on the builder. Used to decide whether validation needs
// the variant fan-out path or can stay on the canonical Linter.Lint flow.
func projectHasVariantPipelines(
	ctx context.Context,
	finder func(root string, defs []string) ([]string, error),
	defs []string,
	rootPath string,
	builder pipelineFromPathBuilder,
	currentVariant *string,
) (bool, error) {
	pipelinePaths, err := finder(rootPath, defs)
	if err != nil {
		return false, err
	}
	previous := *currentVariant
	*currentVariant = ""
	defer func() { *currentVariant = previous }()
	for _, p := range pipelinePaths {
		probe, err := builder.CreatePipelineFromPath(ctx, p, pipeline.WithOnlyPipeline())
		if err != nil {
			return false, fmt.Errorf("error reading pipeline at %q: %w", p, err)
		}
		if len(probe.Variants) > 0 {
			return true, nil
		}
	}
	return false, nil
}

// loadPipelinesExpandingVariants discovers pipelines under rootPath and returns
// one *Pipeline per (path, variant) pair: variant pipelines fan out into one
// instance per declared variant; non-variant pipelines yield a single instance.
//
// currentVariant must be the same pointer captured by a previously-installed
// dynamicVariantMutator on builder, so each materialization sees the right
// variant values during asset construction.
func loadPipelinesExpandingVariants(
	ctx context.Context,
	finder func(root string, defs []string) ([]string, error),
	defs []string,
	rootPath string,
	builder pipelineFromPathBuilder,
	currentVariant *string,
) ([]*pipeline.Pipeline, error) {
	pipelinePaths, err := finder(rootPath, defs)
	if err != nil {
		return nil, err
	}
	if len(pipelinePaths) == 0 {
		return nil, fmt.Errorf("no pipelines found in path %q", rootPath)
	}
	out := make([]*pipeline.Pipeline, 0, len(pipelinePaths))
	for _, p := range pipelinePaths {
		// Probe-load with WithOnlyPipeline to read pipeline.yml metadata only.
		// Mutators will run, but currentVariant is empty during probe so the
		// dynamicVariantMutator no-ops.
		previous := *currentVariant
		*currentVariant = ""
		probe, err := builder.CreatePipelineFromPath(ctx, p, pipeline.WithOnlyPipeline())
		*currentVariant = previous
		if err != nil {
			return nil, fmt.Errorf("error reading pipeline at %q: %w", p, err)
		}
		if len(probe.Variants) == 0 {
			*currentVariant = ""
			pl, err := builder.CreatePipelineFromPath(ctx, p, pipeline.WithMutate())
			if err != nil {
				return nil, err
			}
			out = append(out, pl)
			continue
		}
		for _, name := range probe.Variants.Names() {
			*currentVariant = name
			pl, err := builder.CreatePipelineFromPath(ctx, p, pipeline.WithMutate())
			if err != nil {
				return nil, err
			}
			if err := pl.RenderTemplatedFields(func(vars map[string]any) pipeline.RenderFunc {
				return variantRendererFactory(vars, name)
			}); err != nil {
				return nil, err
			}
			out = append(out, pl)
		}
	}
	*currentVariant = ""
	return out, nil
}

// renderVariantStringsForRun renders templated string fields on the loaded pipeline
// using a renderer that exposes `var.X` and `variant`. It assumes ApplyVariantVariables
// has already merged the variant's values into pl.Variables (typically via the
// variantApplyVariablesMutator).
func renderVariantStringsForRun(pl *pipeline.Pipeline, variantName string) error {
	if variantName == "" {
		return nil
	}
	return pl.RenderTemplatedFields(func(vars map[string]any) pipeline.RenderFunc {
		return variantRendererFactory(vars, variantName)
	})
}

// requireAndApplyVariant enforces that variant pipelines receive a --variant
// flag and that --variant is rejected on non-variant pipelines, then renders
// templated string fields on the loaded pipeline. It is the shared body used
// by mutating commands like `patch fill-asset-dependencies`.
func requireAndApplyVariant(pl *pipeline.Pipeline, variantName string) error {
	if variantName == "" {
		if len(pl.Variants) > 0 {
			return fmt.Errorf("pipeline %q declares variants %v; --variant is required", pl.Name, pl.Variants.Names())
		}
		return nil
	}
	if len(pl.Variants) == 0 {
		return fmt.Errorf("pipeline %q does not declare any variants but --variant=%q was provided", pl.Name, variantName)
	}
	return pl.RenderTemplatedFields(func(vars map[string]any) pipeline.RenderFunc {
		return variantRendererFactory(vars, variantName)
	})
}

// applyVariant materializes a variant on the given pipeline in-place. It performs
// the validation expected of any command that accepts --variant: a variant pipeline
// requires --variant, and --variant on a non-variant pipeline is rejected.
func applyVariant(pl *pipeline.Pipeline, variantName string) error {
	hasVariants := len(pl.Variants) > 0
	switch {
	case variantName != "" && !hasVariants:
		return fmt.Errorf("pipeline %q does not declare any variants but --variant=%q was provided", pl.Name, variantName)
	case variantName == "" && hasVariants:
		return fmt.Errorf("pipeline %q declares variants %v; --variant is required", pl.Name, pl.Variants.Names())
	case variantName == "":
		return nil
	}
	return pl.MaterializeVariant(variantName, variantRendererFactory)
}

// pipelineFromPathBuilder is the minimal subset of *pipeline.Builder needed for
// lint and other commands that load pipelines via the builder.
type pipelineFromPathBuilder interface {
	CreatePipelineFromPath(ctx context.Context, pathToPipeline string, opts ...pipeline.CreatePipelineOption) (*pipeline.Pipeline, error)
}

// variantPipelineBuilder wraps an underlying pipeline builder so that, for variant
// pipelines, templated string fields are rendered after the pipeline is loaded.
// The underlying builder is expected to already have a variant-variables mutator
// registered so that asset-level mutators see the variant's variable values.
type variantPipelineBuilder struct {
	inner       pipelineFromPathBuilder
	variantName string
}

func (vb *variantPipelineBuilder) CreatePipelineFromPath(ctx context.Context, pipelinePath string, opts ...pipeline.CreatePipelineOption) (*pipeline.Pipeline, error) {
	pl, err := vb.inner.CreatePipelineFromPath(ctx, pipelinePath, opts...)
	if err != nil {
		return nil, err
	}
	if vb.variantName == "" || len(pl.Variants) == 0 {
		return pl, nil
	}
	if err := pl.RenderTemplatedFields(func(vars map[string]any) pipeline.RenderFunc {
		return variantRendererFactory(vars, vb.variantName)
	}); err != nil {
		return nil, err
	}
	return pl, nil
}
