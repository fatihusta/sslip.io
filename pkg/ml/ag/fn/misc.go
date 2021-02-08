// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	mat "github.com/nlpodyssey/spago/pkg/mat32"
)

// NewTan returns a new UnaryElementwise tangent function.
func NewTan(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  tan,
		df: tanDeriv,
	}
}

// NewTanh returns a new UnaryElementwise hyperbolic tangent function.
func NewTanh(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  tanh,
		df: tanhDeriv,
	}
}

// NewSigmoid returns a new UnaryElementwise sigmoid function.
func NewSigmoid(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  sigmoid,
		df: sigmoidDeriv,
	}
}

// NewHardSigmoid returns a new UnaryElementwise hard sigmoid function.
func NewHardSigmoid(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  hardSigmoid,
		df: hardSigmoidDeriv,
	}
}

// NewHardTanh returns a new UnaryElementwise hard hyperbolic tangent function.
func NewHardTanh(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  hardTanh,
		df: hardTanhDeriv,
	}
}

// NewReLU returns a new UnaryElementwise Rectified Linear Unit (ReLU) function.
func NewReLU(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  relu,
		df: reluDeriv,
	}
}

// NewSoftsign returns a new UnaryElementwise softsign function.
func NewSoftsign(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  softsign,
		df: softsignDeriv,
	}
}

// NewCos returns a new UnaryElementwise cosine function.
func NewCos(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return mat.Cos(v) },
		df: func(i, j int, v mat.Float) mat.Float { return -mat.Sin(v) },
	}
}

// NewSin returns a new UnaryElementwise sine function.
func NewSin(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return mat.Sin(v) },
		df: func(i, j int, v mat.Float) mat.Float { return mat.Cos(v) },
	}
}

// NewExp returns a new UnaryElementwise base-e exponential function.
func NewExp(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return mat.Exp(v) },
		df: func(i, j int, v mat.Float) mat.Float { return mat.Exp(v) },
	}
}

// NewLog returns a new UnaryElementwise natural logarithm function.
func NewLog(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  safeLog,
		df: safeLogDeriv,
	}
}

// NewNeg returns a new UnaryElementwise f(x) = -x function.
func NewNeg(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return -v },
		df: func(i, j int, v mat.Float) mat.Float { return -1.0 },
	}
}

// NewReciprocal returns a new UnaryElementwise reciprocal function.
func NewReciprocal(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return 1.0 / v },
		df: func(i, j int, v mat.Float) mat.Float { return -1.0 / (v * v) },
	}
}

// NewAbs returns a new UnaryElementwise absolute value function.
func NewAbs(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return mat.Abs(v) },
		df: absDeriv,
	}
}

// NewMish returns a new UnaryElementwise Mish function.
//
// Mish is a self-regularized non-monotonic activation function which can be
// mathematically defined as f(x) = x * tanh(softplus(x)).
//
// Reference: "Mish: A Self Regularized Non-Monotonic Neural Activation Function"
// by Diganta Misra, 2019 (https://arxiv.org/pdf/1908.08681.pdf)
func NewMish(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  mish,
		df: mishDeriv,
	}
}

// NewGELU returns a new UnaryElementwise Gaussian Error Linear Unit (GELU) function.
func NewGELU(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  gelu,
		df: geluDeriv,
	}
}

// NewSqrt returns a new UnaryElementwise square root function.
func NewSqrt(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  func(i, j int, v mat.Float) mat.Float { return mat.Sqrt(v) },
		df: func(i, j int, v mat.Float) mat.Float { return 0.5 * mat.Pow(v, -0.5) },
	}
}

func NewSwish(x Operand) *UnaryElementwise {
	return &UnaryElementwise{
		x:  x,
		f:  swish,
		df: swishDeriv,
	}
}

func absDeriv(i, j int, v mat.Float) mat.Float {
	if v < 0 {
		return -1
	} else if v > 0 {
		return 1
	} else {
		return 0 // undefined
	}
}

// safeLog is a simple work-around that make the math.Log() safe for zero or negative values
func safeLog(i, j int, v mat.Float) mat.Float {
	if v > 0.0 {
		return mat.Log(v)
	} else if v == 0.0 {
		return mat.Log(1.0e-08)
	} else {
		panic("ag: invalid log for negative values")
	}
}

func safeLogDeriv(i, j int, v mat.Float) mat.Float {
	if v > 0.0 {
		return 1.0 / v
	} else if v == 0.0 {
		return 1.0 / 1.0e-08
	} else {
		panic("ag: invalid log for negative values")
	}
}

func tan(i, j int, v mat.Float) mat.Float {
	return mat.Tan(v)
}

func tanDeriv(i, j int, v mat.Float) mat.Float {
	return 1.0 / square(i, j, mat.Cos(v))
}

func square(i, j int, v mat.Float) mat.Float {
	return v * v
}

func tanh(i, j int, v mat.Float) mat.Float {
	return mat.Tanh(v)
}

func tanhDeriv(i, j int, v mat.Float) mat.Float {
	return 1.0 - mat.Pow(mat.Tanh(v), 2.0)
}

func sigmoid(i, j int, v mat.Float) mat.Float {
	return 1.0 / (1 + mat.Exp(-v))
}

func sigmoidDeriv(i, j int, v mat.Float) mat.Float {
	fx := sigmoid(i, j, v)
	return fx * (1.0 - fx)
}

func hardSigmoid(i, j int, v mat.Float) mat.Float {
	if v > 2.5 {
		return 1.0
	} else if v < -2.5 {
		return 0.0
	} else {
		return 0.2*v + 0.5
	}
}

func hardSigmoidDeriv(i, j int, v mat.Float) mat.Float {
	if v < 2.5 && v > -2.5 {
		return 0.2
	}
	return 0.0
}

func hardTanh(i, j int, v mat.Float) mat.Float {
	if v > 1.0 {
		return 1.0
	} else if v < -1.0 {
		return -1.0
	} else {
		return v
	}
}

func hardTanhDeriv(i, j int, v mat.Float) mat.Float {
	if v < 1.0 && v > -1.0 {
		return 1.0
	}
	return 0.0
}

func relu(i, j int, v mat.Float) mat.Float {
	return mat.Max(0.0, v)
}

func reluDeriv(i, j int, v mat.Float) mat.Float {
	if v >= 0.0 {
		return 1.0
	}
	return 0.0
}

func softsign(i, j int, v mat.Float) mat.Float {
	return v / (1.0 + mat.Abs(v))
}

func softsignDeriv(i, j int, v mat.Float) mat.Float {
	return mat.Pow(1.0-mat.Abs(softsign(i, j, v)), 2.0)
}

func celu(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	if v <= 0 {
		return alpha[0] * (mat.Exp(v/alpha[0]) - 1)
	} else if v > 0 {
		return v
	}
	return 0
}

func celuDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	if v <= 0 {
		return mat.Exp(v / alpha[0])
	} else if v > 0 {
		return 1
	}
	return 0
}

func elu(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	if v <= 0 {
		return alpha[0] * (mat.Exp(v) - 1)
	} else if v > 0 {
		return v
	}
	return 0
}

func eluDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	if v <= 0 {
		return alpha[0] * mat.Exp(v)
	} else if v > 0 {
		return 1
	}
	return 0
}

func leakyReLU(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	if v <= 0 {
		return alpha[0] * v // slope * v
	} else if v > 0 {
		return v
	}
	return 0
}

func leakyReLUDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	if v <= 0 {
		return alpha[0] // slope
	} else if v > 0 {
		return 1
	}
	return 0
}

// alpha[0] is the alpha
// alpha[1] is the scale
func selu(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	scale := alpha[1]
	if v <= 0 {
		return scale * alpha[0] * (mat.Exp(v) - 1)
	} else if v > 0 {
		return scale * v
	}
	return 0
}

// alpha[0] is the alpha
// alpha[1] is the scale
func seluDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	scale := alpha[1]
	if v <= 0 {
		return scale * alpha[0] * mat.Exp(v)
	} else if v > 0 {
		return scale
	}
	return 0
}

func softPlus(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	threshold := alpha[1]
	beta := alpha[0]
	if v <= threshold {
		return (1 / beta) * mat.Log(1+mat.Exp(beta*v))
	} else if v > threshold {
		return v
	}
	return 0
}

func softPlusDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	threshold := alpha[1]
	beta := alpha[0]
	if v <= threshold {
		exp := mat.Exp(v * beta)
		return exp / (exp + 1)
	} else if v > threshold {
		return 1
	}
	return 0
}

func softShrink(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	lambda := alpha[0]
	if v < -lambda {
		return v + lambda
	} else if v > lambda {
		return v - lambda
	}
	return 0
}

func softShrinkDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	lambda := alpha[0]
	if v < -lambda {
		return 1
	} else if v > lambda {
		return 1
	}
	return 0
}

func threshold(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	value := alpha[1]
	threshold := alpha[0]
	if v <= threshold {
		return value
	} else if v > threshold {
		return v
	}
	return 0
}

func thresholdDeriv(i, j int, v mat.Float, alpha ...mat.Float) mat.Float {
	threshold := alpha[0]
	if v <= threshold {
		return 0
	} else if v > threshold {
		return 1
	}
	return 0
}

func swish(i, j int, v mat.Float) mat.Float {
	return v * (1.0 / (1 + mat.Exp(-v)))
}

func swishDeriv(i, j int, v mat.Float) mat.Float {
	return swishBDeriv(i, j, v, 1.0)
}

func swishB(i, j int, v mat.Float, beta ...mat.Float) mat.Float {
	return v * (1.0 / (1 + mat.Exp(beta[0]*-v)))
}

func swishBDeriv(i, j int, v mat.Float, beta ...mat.Float) mat.Float {
	prod := v * beta[0]
	exp := mat.Exp(prod)
	return exp * (exp + prod + 1) / ((exp + 1) * (exp + 1))
}

func swishBBetaDeriv(v mat.Float, beta mat.Float) mat.Float {
	prod := v * beta
	exp := mat.Exp(-prod)
	return (v * v * exp) / ((exp + 1) * (exp + 1))
}

// Reference: "Mish: A Self Regularized Non-Monotonic Neural Activation Function" by Diganta Misra, 2019.
// (https://arxiv.org/pdf/1908.08681.pdf)
func mish(i, j int, v mat.Float) mat.Float {
	return v * mat.Tanh(mat.Log(1+mat.Exp(v)))
}

func mishDeriv(i, j int, v mat.Float) mat.Float {
	exp := mat.Exp(v)
	exp2 := mat.Exp(2 * v)
	exp3 := mat.Exp(3 * v)
	omega := 4.0*(v+1.0) + 4.0*exp2 + exp3 + exp*(4.0*v+6.0)
	delta := 2*exp + exp2 + 2.0
	return exp * (omega / (delta * delta))
}

func gelu(i, j int, v mat.Float) mat.Float {
	return 0.5 * v * (1.0 + mat.Tanh(mat.Sqrt(2/mat.Pi)*(v+0.044715*mat.Pow(v, 3.0))))
}

func geluDeriv(i, j int, x mat.Float) mat.Float {
	x3 := mat.Pow(x, 3)
	return 0.5*mat.Tanh(0.0356774*x3+0.797885*x) +
		(0.0535161*x3+0.398942*x)*
			mat.Pow(1.0/mat.Cosh(0.0356774*x3+0.797885*x), 2) + 0.5
}
