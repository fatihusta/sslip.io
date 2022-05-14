// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/mat"
)

// SubScalar is an element-wise subtraction function with a scalar value.
type SubScalar[T mat.DType, O Operand[T]] struct {
	x1       O
	x2       O // scalar
	operands []O
}

// NewSubScalar returns a new SubScalar Function.
func NewSubScalar[T mat.DType, O Operand[T]](x1 O, x2 O) *SubScalar[T, O] {
	return &SubScalar[T, O]{
		x1:       x1,
		x2:       x2,
		operands: []O{x1, x2},
	}
}

// Operands returns the list of operands.
func (r *SubScalar[T, O]) Operands() []O {
	return r.operands
}

// Forward computes the output of the node.
func (r *SubScalar[T, O]) Forward() mat.Matrix {
	return r.x1.Value().SubScalar(r.x2.Value().Scalar().Float64())
}

// Backward computes the backward pass.
func (r *SubScalar[T, O]) Backward(gy mat.Matrix) {
	if !(mat.SameDims(r.x1.Value(), gy) || mat.VectorsOfSameSize(r.x1.Value(), gy)) {
		panic("fn: matrices with not compatible size")
	}
	if r.x1.RequiresGrad() {
		r.x1.AccGrad(gy) // equals to gy.ProdScalar(1.0)
	}
	if r.x2.RequiresGrad() {
		neg := gy.ProdScalar(-1)
		defer mat.ReleaseMatrix(neg)
		gx := neg.Sum()
		defer mat.ReleaseMatrix(gx)
		r.x2.AccGrad(gx)
	}
}
