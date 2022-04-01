// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import "github.com/nlpodyssey/spago/mat"

// ReverseSubScalar is the element-wise subtraction function over two values.
type ReverseSubScalar[T mat.DType, O Operand[T]] struct {
	x1       O
	x2       O // scalar
	operands []O
}

// NewReverseSubScalar returns a new ReverseSubScalar Function.
func NewReverseSubScalar[T mat.DType, O Operand[T]](x1 O, x2 O) *ReverseSubScalar[T, O] {
	return &ReverseSubScalar[T, O]{
		x1:       x1,
		x2:       x2,
		operands: []O{x1, x2},
	}
}

// Operands returns the list of operands.
func (r *ReverseSubScalar[T, O]) Operands() []O {
	return r.operands
}

// Forward computes the output of the function.
func (r *ReverseSubScalar[T, O]) Forward() mat.Matrix[T] {
	x1v, x2v := r.x1.Value(), r.x2.Value()
	return mat.NewInitDense(x1v.Rows(), x1v.Columns(), x2v.Scalar()).Sub(x1v)
}

// Backward computes the backward pass.
func (r *ReverseSubScalar[T, O]) Backward(gy mat.Matrix[T]) {
	if !(mat.SameDims(r.x1.Value(), gy) || mat.VectorsOfSameSize(r.x1.Value(), gy)) {
		panic("fn: matrices with not compatible size")
	}
	if r.x1.RequiresGrad() {
		gx := gy.ProdScalar(-1.0)
		defer mat.ReleaseMatrix(gx)
		r.x1.PropagateGrad(gx)
	}
	if r.x2.RequiresGrad() {
		var gx T = 0.0
		for i := 0; i < gy.Rows(); i++ {
			for j := 0; j < gy.Columns(); j++ {
				gx += gy.At(i, j)
			}
		}
		scalar := mat.NewScalar(gx)
		defer mat.ReleaseDense(scalar)
		r.x2.PropagateGrad(scalar)
	}
}
