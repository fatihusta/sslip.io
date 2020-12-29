// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqrdist

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
)

var (
	_ nn.Model = &Model{}
)

// Model contains the serializable parameters.
type Model struct {
	nn.BaseModel
	B nn.Param `spago:"type:weights"`
}

// New returns a new model with parameters initialized to zeros.
func New(in, rank int) *Model {
	return &Model{
		BaseModel: nn.BaseModel{RCS: false},
		B:         nn.NewParam(mat.NewEmptyDense(rank, in)),
	}
}

// Forward performs the forward step for each input node and returns the result.
func (m *Model) Forward(in interface{}) interface{} {
	xs := nn.ToNodes(in)
	ys := make([]ag.Node, len(xs))
	for i, x := range xs {
		ys[i] = m.forward(x)
	}
	return ys
}

func (m *Model) forward(x ag.Node) ag.Node {
	g := m.Graph()
	bh := g.Mul(m.B, x)
	return g.Mul(g.T(bh), bh)
}
