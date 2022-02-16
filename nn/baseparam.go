// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nn

import (
	"sync"

	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
)

var _ Param[float32] = &BaseParam[float32]{}

// BaseParam is the default implementation satisfying the Param interface.
type BaseParam[T mat.DType] struct {
	name         string
	pType        ParamsType    // lazy initialization
	value        mat.Matrix[T] // store the results of a forward evaluation.
	grad         mat.Matrix[T]
	payload      *Payload[T] // additional data used for example by gradient-descend optimization methods
	hasGrad      bool
	requiresGrad bool
	// Allows thread-safe locking for operations on value.
	valueMu sync.RWMutex
	// Allows thread-safe locking for operations on grad.
	gradMu sync.RWMutex
	// Allows thread-safe locking for operations on payload.
	payloadMu sync.RWMutex
}

// ParamOption allows to configure a new Param with your specific needs.
type ParamOption[T mat.DType] func(*BaseParam[T])

// RequiresGrad is an option to specify whether a Param should be trained or not.
func RequiresGrad[T mat.DType](value bool) ParamOption[T] {
	return func(p *BaseParam[T]) {
		p.requiresGrad = value
	}
}

// NewParam returns a new param.
func NewParam[T mat.DType](value mat.Matrix[T], opts ...ParamOption[T]) Param[T] {
	p := &BaseParam[T]{
		name:         "",        // lazy initialization
		pType:        Undefined, // lazy initialization
		value:        value,
		grad:         nil, // lazy initialization
		hasGrad:      false,
		requiresGrad: true, // true by default, can be modified with the options
		payload:      nil,  // lazy initialization
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// SetName set the params name (can be empty string).
func (p *BaseParam[_]) SetName(name string) {
	p.name = name
}

// SetType set the params type (weights, biases, undefined).
func (p *BaseParam[_]) SetType(pType ParamsType) {
	p.pType = pType
}

// Name returns the params name (can be empty string).
func (p *BaseParam[_]) Name() string {
	return p.name
}

// Type returns the params type (weights, biases, undefined).
func (p *BaseParam[_]) Type() ParamsType {
	return p.pType
}

// Value returns the value of the delegate itself.
func (p *BaseParam[T]) Value() mat.Matrix[T] {
	p.valueMu.RLock()
	defer p.valueMu.RUnlock()
	return p.value
}

// ReplaceValue replaces the value of the parameter and clears the gradients and
// the support structure.
func (p *BaseParam[T]) ReplaceValue(value mat.Matrix[T]) {
	p.ClearPayload()
	p.ZeroGrad()

	p.valueMu.Lock()
	defer p.valueMu.Unlock()
	p.value = value
}

// ScalarValue returns the scalar value of the node.
// It panics if the value is not a scalar.
// Note that it is not possible to start the backward step from a scalar value.
func (p *BaseParam[T]) ScalarValue() T {
	p.valueMu.RLock()
	defer p.valueMu.RUnlock()
	return p.value.Scalar()
}

// Grad returns the gradients accumulated during the backward pass.
func (p *BaseParam[T]) Grad() mat.Matrix[T] {
	p.gradMu.RLock()
	defer p.gradMu.RUnlock()
	return p.grad
}

// PropagateGrad accumulate the gradients
func (p *BaseParam[T]) PropagateGrad(grad mat.Matrix[T]) {
	if !p.requiresGrad {
		return
	}
	p.gradMu.Lock()
	defer p.gradMu.Unlock()
	if p.grad == nil {
		p.grad = p.value.ZerosLike()
	}
	p.grad.AddInPlace(grad)
	p.hasGrad = true
}

// HasGrad returns true if there are accumulated gradients.
func (p *BaseParam[_]) HasGrad() bool {
	p.gradMu.RLock()
	defer p.gradMu.RUnlock()
	return p.hasGrad
}

// RequiresGrad returns true if the param requires gradients.
func (p *BaseParam[_]) RequiresGrad() bool {
	return p.requiresGrad
}

// SetRequiresGrad is an option to specify whether a Param should be trained or not.
func (p *BaseParam[_]) SetRequiresGrad(value bool) {
	p.requiresGrad = value
}

// ZeroGrad clears the gradients.
func (p *BaseParam[_]) ZeroGrad() {
	if p.grad == nil {
		return
	}
	p.gradMu.Lock()
	defer p.gradMu.Unlock()
	mat.ReleaseMatrix(p.grad)
	p.grad = nil
	p.hasGrad = false
}

// ApplyDelta updates the value applying the delta.
func (p *BaseParam[T]) ApplyDelta(delta mat.Matrix[T]) {
	p.valueMu.Lock()
	defer p.valueMu.Unlock()
	p.value.SubInPlace(delta)
}

// Payload returns the optimizer support structure (can be nil).
func (p *BaseParam[T]) Payload() *Payload[T] {
	p.payloadMu.RLock()
	defer p.payloadMu.RUnlock()
	return p.payload
}

// SetPayload is a thread safe operation to set the given Payload on the
// receiver Param.
func (p *BaseParam[T]) SetPayload(payload *Payload[T]) {
	p.payloadMu.Lock()
	defer p.payloadMu.Unlock()
	p.payload = payload
}

// ClearPayload clears the support structure.
func (p *BaseParam[_]) ClearPayload() {
	if p.payload == nil {
		return
	}
	p.payloadMu.Lock()
	defer p.payloadMu.Unlock()
	p.payload.ClearData()
	p.payload = nil
}

// Graph returns always nil since the "pure" parameter is not associated with any graph.
func (p *BaseParam[T]) Graph() *ag.Graph[T] {
	panic("nn: attempting to access Graph on a not reified param.")
}

// ID returns always -1 since the "pure" parameter is not associated with any graph.
func (p *BaseParam[_]) ID() int {
	panic("nn: attempting to access the ID of a not reified param.")
}

// TimeStep returns always 0 since the "pure" parameter is not associated with any graph.
func (p *BaseParam[_]) TimeStep() int {
	panic("nn: attempting to access the TimeStep of a not reified param.")
}