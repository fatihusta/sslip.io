// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crf

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
)

var (
	_ nn.Module = &Model{}
)

// Model contains the serializable parameters.
type Model struct {
	nn.BaseModel
	Size             int
	TransitionScores nn.Param    `type:"weights"`
	Scores           [][]ag.Node `scope:"processor"`
}

// New returns a new convolution Model, initialized according to the given configuration.
func New(size int) *Model {
	return &Model{
		BaseModel:        nn.BaseModel{FullSeqProcessing: true},
		Size:             size,
		TransitionScores: nn.NewParam(mat.NewEmptyDense(size+1, size+1)), // +1 for start and end transitions
	}
}

func (m *Model) InitProc() {
	m.Scores = nn.Separate(m.GetGraph(), m.TransitionScores) // TODO: lazy initialization
}

// Predict performs the forward step for each input and returns the result.
func (m *Model) Predict(emissionScores []ag.Node) []int {
	return Viterbi(m.TransitionScores.Value(), emissionScores)
}

// Forward is not available for the CRF. Use Predict() instead.
func (m *Model) Forward(_ ...ag.Node) []ag.Node {
	panic("crf: Forward() not available. Use Predict() instead.")
}

// NegativeLogLoss computes the negative log loss with respect to the targets.
func (m *Model) NegativeLogLoss(emissionScores []ag.Node, target []int) ag.Node {
	goldScore := m.goldScore(emissionScores, target)
	totalScore := m.totalScore(emissionScores)
	return m.GetGraph().Sub(totalScore, goldScore)
}

func (m *Model) goldScore(emissionScores []ag.Node, target []int) ag.Node {
	g := m.GetGraph()
	goldScore := g.At(emissionScores[0], target[0], 0)
	goldScore = g.Add(goldScore, m.Scores[0][target[0]+1]) // start transition
	prevIndex := target[0] + 1
	for i := 1; i < len(target); i++ {
		goldScore = g.Add(goldScore, g.AtVec(emissionScores[i], target[i]))
		goldScore = g.Add(goldScore, m.Scores[prevIndex][target[i]+1])
		prevIndex = target[i] + 1
	}
	goldScore = g.Add(goldScore, m.Scores[prevIndex][0]) // end transition
	return goldScore
}

func (m *Model) totalScore(predicted []ag.Node) ag.Node {
	g := m.GetGraph()
	totalVector := m.totalScoreStart(predicted[0])
	for i := 1; i < len(predicted); i++ {
		totalVector = m.totalScoreStep(totalVector, nn.SeparateVec(g, predicted[i]))
	}
	totalVector = m.totalScoreEnd(totalVector)
	return g.Log(g.ReduceSum(g.Concat(totalVector...)))

}

func (m *Model) totalScoreStart(stepVec ag.Node) []ag.Node {
	firstTransitionScores := m.Scores[0]
	scores := make([]ag.Node, m.Size)
	g := m.GetGraph()
	for i := 0; i < m.Size; i++ {
		scores[i] = g.Add(g.AtVec(stepVec, i), firstTransitionScores[i+1])
	}
	return scores
}

func (m *Model) totalScoreEnd(stepVec []ag.Node) []ag.Node {
	scores := make([]ag.Node, m.Size)
	g := m.GetGraph()
	for i := 0; i < m.Size; i++ {
		vecTrans := g.Add(stepVec[i], m.Scores[i+1][0])
		scores[i] = g.Add(scores[i], g.Exp(vecTrans))
	}
	return scores
}

func (m *Model) totalScoreStep(totalVec []ag.Node, stepVec []ag.Node) []ag.Node {
	scores := make([]ag.Node, m.Size)
	g := m.GetGraph()
	for i := 0; i < m.Size; i++ {
		nodei := totalVec[i]
		transitionScores := m.Scores[i+1]
		for j := 0; j < m.Size; j++ {
			vecSum := g.Add(nodei, stepVec[j])
			vecTrans := g.Add(vecSum, transitionScores[j+1])
			scores[j] = g.Add(scores[j], g.Exp(vecTrans))
		}
	}
	for i := 0; i < m.Size; i++ {
		scores[i] = g.Log(scores[i])
	}
	return scores
}
