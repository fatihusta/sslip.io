// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hyperbolic

import "github.com/nlpodyssey/spago/mat"

// Hyperbolic defines an hyperbolic decay depending on the time step
//     lr = lr / (1 + rate*t).
type Hyperbolic[T mat.DType] struct {
	init  T
	final T
	rate  T
}

// New returns a new Hyperbolic decay optimizer.
func New[T mat.DType](init, final, rate T) *Hyperbolic[T] {
	if init < final {
		panic("decay: the initial learning rate must be >= than the final one")
	}
	return &Hyperbolic[T]{
		init:  init,
		final: final,
		rate:  rate,
	}
}

// Decay calculates the decay of the learning rate lr at time t.
func (d *Hyperbolic[T]) Decay(lr T, t int) T {
	if t > 1 && d.rate > 0.0 && lr > d.final {
		return d.init / (1.0 + d.rate*T(t))
	}
	return lr
}