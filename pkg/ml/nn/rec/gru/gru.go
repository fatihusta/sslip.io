// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gru

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"log"
)

var (
	_ nn.Model = &Model{}
)

// Model contains the serializable parameters.
type Model struct {
	nn.BaseModel
	WPart    nn.Param `spago:"type:weights"`
	WPartRec nn.Param `spago:"type:weights"`
	BPart    nn.Param `spago:"type:biases"`
	WRes     nn.Param `spago:"type:weights"`
	WResRec  nn.Param `spago:"type:weights"`
	BRes     nn.Param `spago:"type:biases"`
	WCand    nn.Param `spago:"type:weights"`
	WCandRec nn.Param `spago:"type:weights"`
	BCand    nn.Param `spago:"type:biases"`
	States   []*State `spago:"scope:processor"`
}

// State represent a state of the GRU recurrent network.
type State struct {
	R ag.Node
	P ag.Node
	C ag.Node
	Y ag.Node
}

// New returns a new model with parameters initialized to zeros.
func New(in, out int) *Model {
	m := &Model{
		BaseModel: nn.BaseModel{RCS: false},
	}
	m.WPart, m.WPartRec, m.BPart = newGateParams(in, out)
	m.WRes, m.WResRec, m.BRes = newGateParams(in, out)
	m.WCand, m.WCandRec, m.BCand = newGateParams(in, out)
	return m
}

func newGateParams(in, out int) (w, wRec, b nn.Param) {
	w = nn.NewParam(mat.NewEmptyDense(out, in))
	wRec = nn.NewParam(mat.NewEmptyDense(out, out))
	b = nn.NewParam(mat.NewEmptyVecDense(out))
	return
}

// SetInitialState sets the initial state of the recurrent network.
// It panics if one or more states are already present.
func (m *Model) SetInitialState(state *State) {
	if len(m.States) > 0 {
		log.Fatal("gru: the initial state must be set before any input")
	}
	m.States = append(m.States, state)
}

// Forward performs the forward step for each input node and returns the result.
func (m *Model) Forward(in interface{}) interface{} {
	xs := nn.ToNodes(in)
	ys := make([]ag.Node, len(xs))
	for i, x := range xs {
		s := m.forward(x)
		m.States = append(m.States, s)
		ys[i] = s.Y
	}
	return ys
}

// LastState returns the last state of the recurrent network.
// It returns nil if there are no states.
func (m *Model) LastState() *State {
	n := len(m.States)
	if n == 0 {
		return nil
	}
	return m.States[n-1]
}

// r = sigmoid(wr (dot) x + br + wrRec (dot) yPrev)
// p = sigmoid(wp (dot) x + bp + wpRec (dot) yPrev)
// c = f(wc (dot) x + bc + wcRec (dot) (yPrev * r))
// y = p * c + (1 - p) * yPrev
func (m *Model) forward(x ag.Node) (s *State) {
	g := m.Graph()
	s = new(State)
	yPrev := m.prev()
	s.R = g.Sigmoid(nn.Affine(g, m.BRes, m.WRes, x, m.WResRec, yPrev))
	s.P = g.Sigmoid(nn.Affine(g, m.BPart, m.WPart, x, m.WPartRec, yPrev))
	s.C = g.Tanh(nn.Affine(g, m.BCand, m.WCand, x, m.WCandRec, tryProd(g, yPrev, s.R)))
	s.Y = g.Prod(s.P, s.C)
	if yPrev != nil {
		s.Y = g.Add(s.Y, g.Prod(g.ReverseSub(s.P, g.NewScalar(1.0)), yPrev))
	}
	return
}

func (m *Model) prev() (yPrev ag.Node) {
	s := m.LastState()
	if s != nil {
		yPrev = s.Y
	}
	return
}

// tryProd returns the product if 'a' il not nil, otherwise nil
func tryProd(g *ag.Graph, a, b ag.Node) ag.Node {
	if a != nil {
		return g.Prod(a, b)
	}
	return nil
}
