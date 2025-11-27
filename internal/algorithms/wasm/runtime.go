package wasm

import (
	"context"
	"fmt"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Runtime manages WASM module compilation and caching.
type Runtime struct {
	runtime wazero.Runtime
	cache   map[string]wazero.CompiledModule
	mu      sync.RWMutex
}

// NewRuntime creates a new WASM runtime with WASI support.
func NewRuntime(ctx context.Context) (*Runtime, error) {
	cfg := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true)

	rt := wazero.NewRuntimeWithConfig(ctx, cfg)

	// Instantiate WASI for basic I/O operations
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("instantiate WASI: %w", err)
	}

	return &Runtime{
		runtime: rt,
		cache:   make(map[string]wazero.CompiledModule),
	}, nil
}

// CompileModule compiles WASM binary and caches it by hash.
func (r *Runtime) CompileModule(ctx context.Context, wasmHash string, wasmBinary []byte) (wazero.CompiledModule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check cache first
	if cached, ok := r.cache[wasmHash]; ok {
		return cached, nil
	}

	// Compile module
	compiled, err := r.runtime.CompileModule(ctx, wasmBinary)
	if err != nil {
		return nil, fmt.Errorf("compile WASM module: %w", err)
	}

	// Cache compiled module
	r.cache[wasmHash] = compiled

	return compiled, nil
}

// GetCachedModule returns a cached compiled module if available.
func (r *Runtime) GetCachedModule(wasmHash string) (wazero.CompiledModule, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	module, ok := r.cache[wasmHash]
	return module, ok
}

// RemoveFromCache removes a module from cache (e.g., when algorithm is deleted).
func (r *Runtime) RemoveFromCache(wasmHash string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if module, ok := r.cache[wasmHash]; ok {
		_ = module.Close(context.Background())
		delete(r.cache, wasmHash)
	}
}

// InstantiateModule creates a new instance of a compiled module.
func (r *Runtime) InstantiateModule(ctx context.Context, compiled wazero.CompiledModule, name string) (api.Module, error) {
	cfg := wazero.NewModuleConfig().
		WithName(name).
		WithStartFunctions() // Don't auto-call _start

	module, err := r.runtime.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		return nil, fmt.Errorf("instantiate module: %w", err)
	}

	return module, nil
}

// Close releases all resources.
func (r *Runtime) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Close all cached modules
	for hash, module := range r.cache {
		_ = module.Close(ctx)
		delete(r.cache, hash)
	}

	return r.runtime.Close(ctx)
}

// ValidateModule checks if WASM binary is valid and has required exports.
func (r *Runtime) ValidateModule(ctx context.Context, wasmBinary []byte, kind ModuleKind) error {
	compiled, err := r.runtime.CompileModule(ctx, wasmBinary)
	if err != nil {
		return fmt.Errorf("invalid WASM module: %w", err)
	}
	defer func() { _ = compiled.Close(ctx) }()

	exports := compiled.ExportedFunctions()

	// Check required exports
	requiredExports := []string{"alloc", "dealloc", "evaluate"}
	for _, name := range requiredExports {
		if _, ok := exports[name]; !ok {
			return fmt.Errorf("missing required export: %s", name)
		}
	}

	// handle_feedback is optional but recommended
	if _, ok := exports["handle_feedback"]; !ok {
		// Log warning but don't fail
	}

	return nil
}

// ModuleKind represents the type of WASM algorithm module.
type ModuleKind string

const (
	ModuleKindBandit           ModuleKind = "bandit"
	ModuleKindOptimizer        ModuleKind = "optimizer"
	ModuleKindContextualBandit ModuleKind = "contextual_bandit"
)
