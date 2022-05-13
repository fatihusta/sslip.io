// Copyright 2022 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat_test

import (
	"testing"

	"github.com/nlpodyssey/spago/mat"
	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	t.Run("float32", testFloat[float32])
	t.Run("float64", testFloat[float64])
}

func testFloat[T mat.DType](t *testing.T) {
	v := mat.Float[T](T(42))
	assert.Equal(t, float32(42), v.Float32())
	assert.Equal(t, float64(42), v.Float64())
}

func TestDTFloat(t *testing.T) {
	t.Run("it returns the correct value according to the type", func(t *testing.T) {
		f := fakeFloat{f32: 32, f64: 64}
		assert.Equal(t, float32(32), mat.DTFloat[float32](f))
		assert.Equal(t, float64(64), mat.DTFloat[float64](f))
	})

	t.Run("it panics with nil", func(t *testing.T) {
		assert.Panics(t, func() { mat.DTFloat[float32](nil) })
		assert.Panics(t, func() { mat.DTFloat[float64](nil) })
	})
}

type fakeFloat struct {
	f32 float32
	f64 float64
}

func (f fakeFloat) Float32() float32 { return f.f32 }
func (f fakeFloat) Float64() float64 { return f.f64 }
