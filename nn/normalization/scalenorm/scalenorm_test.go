// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scalenorm

import (
	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/nn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_Forward(t *testing.T) {
	t.Run("float32", testModelForward[float32])
	t.Run("float64", testModelForward[float64])
}

func testModelForward[T mat.DType](t *testing.T) {
	model := newTestModel[T]()
	g := ag.NewGraph[T](ag.WithMode[T](ag.Training))

	// == Forward
	x1 := g.NewVariable(mat.NewVecDense([]T{1.0, 2.0, 0.0, 4.0}), true)
	x2 := g.NewVariable(mat.NewVecDense([]T{3.0, 2.0, 1.0, 6.0}), true)
	x3 := g.NewVariable(mat.NewVecDense([]T{6.0, 2.0, 5.0, 1.0}), true)

	y := nn.Reify(model, g).Forward(x1, x2, x3)

	assert.InDeltaSlice(t, []T{0.1091089451, -0.0872871560, 0.0, 0.6982972487}, y[0].Value().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []T{0.2121320343, -0.0565685424, 0.0424264068, 0.6788225099}, y[1].Value().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []T{0.3692744729, -0.0492365963, 0.1846372364, 0.0984731927}, y[2].Value().Data(), 1.0e-06)

	// == Backward
	y[0].PropagateGrad(mat.NewVecDense([]T{-1.0, -0.2, 0.4, 0.6}))
	y[1].PropagateGrad(mat.NewVecDense([]T{-0.3, 0.1, 0.7, 0.9}))
	y[2].PropagateGrad(mat.NewVecDense([]T{0.3, -0.4, 0.7, -0.8}))
	g.BackwardAll()

	assert.InDeltaSlice(t, []T{-0.1246959373, -0.0224452687, 0.0261861468, 0.0423966187}, x1.Grad().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []T{-0.0554937402, -0.0256821183, 0.0182716392, 0.033262303}, x2.Grad().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []T{0.0020142244, 0.0043641529, 0.0121412971, -0.0815201374}, x3.Grad().Data(), 1.0e-06)
}

func newTestModel[T mat.DType]() *Model[T] {
	model := New[T](4)
	model.Gain.Value().SetData([]T{0.5, -0.2, 0.3, 0.8})
	return model
}