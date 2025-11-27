package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	"github.com/togglr-project/togglr/internal/domain"
)

const (
	// DefaultTimeout for WASM execution
	DefaultTimeout = 100 * time.Millisecond
	// MaxMemoryPages limits memory usage (64KB per page)
	MaxMemoryPages = 16 // 1MB max
)

// Executor handles WASM module execution with memory management.
type Executor struct {
	runtime *Runtime
	timeout time.Duration
}

// NewExecutor creates a new WASM executor.
func NewExecutor(runtime *Runtime) *Executor {
	return &Executor{
		runtime: runtime,
		timeout: DefaultTimeout,
	}
}

// SetTimeout configures execution timeout.
func (e *Executor) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// ModuleInstance represents an instantiated WASM module ready for execution.
type ModuleInstance struct {
	module         api.Module
	alloc          api.Function
	dealloc        api.Function
	evaluate       api.Function
	handleFeedback api.Function // optional
	mu             sync.Mutex
}

// InstantiateModule creates a ready-to-use module instance.
func (e *Executor) InstantiateModule(ctx context.Context, compiled wazero.CompiledModule, instanceName string) (*ModuleInstance, error) {
	module, err := e.runtime.InstantiateModule(ctx, compiled, instanceName)
	if err != nil {
		return nil, err
	}

	alloc := module.ExportedFunction("alloc")
	if alloc == nil {
		_ = module.Close(ctx)
		return nil, fmt.Errorf("missing alloc function")
	}

	dealloc := module.ExportedFunction("dealloc")
	if dealloc == nil {
		_ = module.Close(ctx)
		return nil, fmt.Errorf("missing dealloc function")
	}

	evaluate := module.ExportedFunction("evaluate")
	if evaluate == nil {
		_ = module.Close(ctx)
		return nil, fmt.Errorf("missing evaluate function")
	}

	// handle_feedback is optional
	handleFeedback := module.ExportedFunction("handle_feedback")

	return &ModuleInstance{
		module:         module,
		alloc:          alloc,
		dealloc:        dealloc,
		evaluate:       evaluate,
		handleFeedback: handleFeedback,
	}, nil
}

// Evaluate calls the evaluate function with input and returns output.
func (e *Executor) Evaluate(ctx context.Context, instance *ModuleInstance, input *domain.WASMInput) (*domain.WASMOutput, error) {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Serialize input to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal input: %w", err)
	}

	// Allocate memory for input
	inputPtr, err := e.allocateAndWrite(ctx, instance, inputJSON)
	if err != nil {
		return nil, fmt.Errorf("allocate input: %w", err)
	}
	defer e.deallocate(ctx, instance, inputPtr, uint32(len(inputJSON)))

	// Call evaluate function
	results, err := instance.evaluate.Call(ctx, uint64(inputPtr), uint64(len(inputJSON)))
	if err != nil {
		return nil, fmt.Errorf("call evaluate: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("evaluate returned no results")
	}

	// Parse output pointer
	outputPtr := uint32(results[0])
	if outputPtr == 0 {
		return nil, fmt.Errorf("evaluate returned null pointer")
	}

	// Read output from memory
	outputJSON, outputLen, err := e.readOutput(ctx, instance, outputPtr)
	if err != nil {
		return nil, fmt.Errorf("read output: %w", err)
	}
	defer e.deallocate(ctx, instance, outputPtr, outputLen)

	// Parse output
	var output domain.WASMOutput
	if err := json.Unmarshal(outputJSON, &output); err != nil {
		return nil, fmt.Errorf("unmarshal output: %w", err)
	}

	if output.Error != "" {
		return nil, fmt.Errorf("WASM error: %s", output.Error)
	}

	return &output, nil
}

// HandleFeedback calls the handle_feedback function if available.
func (e *Executor) HandleFeedback(ctx context.Context, instance *ModuleInstance, input *domain.WASMFeedbackInput) (*domain.WASMFeedbackOutput, error) {
	if instance.handleFeedback == nil {
		// No feedback handler, return empty output
		return &domain.WASMFeedbackOutput{}, nil
	}

	instance.mu.Lock()
	defer instance.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Serialize input to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal feedback input: %w", err)
	}

	// Allocate memory for input
	inputPtr, err := e.allocateAndWrite(ctx, instance, inputJSON)
	if err != nil {
		return nil, fmt.Errorf("allocate feedback input: %w", err)
	}
	defer e.deallocate(ctx, instance, inputPtr, uint32(len(inputJSON)))

	// Call handle_feedback function
	results, err := instance.handleFeedback.Call(ctx, uint64(inputPtr), uint64(len(inputJSON)))
	if err != nil {
		return nil, fmt.Errorf("call handle_feedback: %w", err)
	}

	if len(results) == 0 {
		return &domain.WASMFeedbackOutput{}, nil
	}

	// Parse output pointer
	outputPtr := uint32(results[0])
	if outputPtr == 0 {
		return &domain.WASMFeedbackOutput{}, nil
	}

	// Read output from memory
	outputJSON, outputLen, err := e.readOutput(ctx, instance, outputPtr)
	if err != nil {
		return nil, fmt.Errorf("read feedback output: %w", err)
	}
	defer e.deallocate(ctx, instance, outputPtr, outputLen)

	// Parse output
	var output domain.WASMFeedbackOutput
	if err := json.Unmarshal(outputJSON, &output); err != nil {
		return nil, fmt.Errorf("unmarshal feedback output: %w", err)
	}

	return &output, nil
}

// allocateAndWrite allocates memory and writes data to it.
func (e *Executor) allocateAndWrite(ctx context.Context, instance *ModuleInstance, data []byte) (uint32, error) {
	size := uint32(len(data))

	results, err := instance.alloc.Call(ctx, uint64(size))
	if err != nil {
		return 0, fmt.Errorf("alloc call failed: %w", err)
	}

	if len(results) == 0 {
		return 0, fmt.Errorf("alloc returned no results")
	}

	ptr := uint32(results[0])
	if ptr == 0 {
		return 0, fmt.Errorf("alloc returned null pointer")
	}

	// Write data to memory
	mem := instance.module.Memory()
	if mem == nil {
		return 0, fmt.Errorf("no memory exported")
	}

	if !mem.Write(ptr, data) {
		return 0, fmt.Errorf("memory write failed")
	}

	return ptr, nil
}

// deallocate frees allocated memory.
func (e *Executor) deallocate(ctx context.Context, instance *ModuleInstance, ptr, size uint32) {
	if ptr == 0 {
		return
	}
	_, _ = instance.dealloc.Call(ctx, uint64(ptr), uint64(size))
}

// readOutput reads a null-terminated or length-prefixed string from memory.
// Output format: first 4 bytes are length (little-endian), followed by JSON data.
func (e *Executor) readOutput(ctx context.Context, instance *ModuleInstance, ptr uint32) ([]byte, uint32, error) {
	mem := instance.module.Memory()
	if mem == nil {
		return nil, 0, fmt.Errorf("no memory exported")
	}

	// Read length (first 4 bytes)
	lenBytes, ok := mem.Read(ptr, 4)
	if !ok {
		return nil, 0, fmt.Errorf("failed to read output length")
	}

	length := uint32(lenBytes[0]) |
		uint32(lenBytes[1])<<8 |
		uint32(lenBytes[2])<<16 |
		uint32(lenBytes[3])<<24

	if length == 0 || length > 1024*1024 { // Max 1MB
		return nil, 0, fmt.Errorf("invalid output length: %d", length)
	}

	// Read data
	data, ok := mem.Read(ptr+4, length)
	if !ok {
		return nil, 0, fmt.Errorf("failed to read output data")
	}

	return data, length + 4, nil
}

// Close releases the module instance.
func (mi *ModuleInstance) Close(ctx context.Context) error {
	return mi.module.Close(ctx)
}

// HasFeedbackHandler returns true if the module has a handle_feedback function.
func (mi *ModuleInstance) HasFeedbackHandler() bool {
	return mi.handleFeedback != nil
}
