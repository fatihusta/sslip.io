package ag

import (
	"fmt"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/mat/mattest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewVariable(t *testing.T) {
	t.Run("float32", testNewVariable[float32])
	t.Run("float64", testNewVariable[float64])
}

func testNewVariable[T mat.DType](t *testing.T) {
	testCases := []struct {
		value        mat.Matrix[T]
		requiresGrad bool
	}{
		{mat.NewScalar[T](42), true},
		{mat.NewScalar[T](42), false},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("NewVariable(%g, %v)", tc.value, tc.requiresGrad)
		t.Run(name, func(t *testing.T) {
			v := NewVariable[T](tc.value, tc.requiresGrad)
			require.NotNil(t, v)
			assert.Empty(t, v.Name())
			assert.Same(t, tc.value, v.Value())
			assert.Nil(t, v.Grad())
			assert.False(t, v.HasGrad())
			assert.Equal(t, tc.requiresGrad, v.RequiresGrad())
		})
	}
}

func TestNewVariableWithName(t *testing.T) {
	t.Run("float32", testNewVariableWithName[float32])
	t.Run("float64", testNewVariableWithName[float64])
}

func testNewVariableWithName[T mat.DType](t *testing.T) {
	testCases := []struct {
		value        mat.Matrix[T]
		requiresGrad bool
		name         string
	}{
		{mat.NewScalar[T](42), true, "foo"},
		{mat.NewScalar[T](42), false, "bar"},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("NewVariableWithName(%g, %v, %#v)", tc.value, tc.requiresGrad, tc.name)
		t.Run(name, func(t *testing.T) {
			v := NewVariableWithName[T](tc.value, tc.requiresGrad, tc.name)
			require.NotNil(t, v)
			assert.Equal(t, tc.name, v.Name())
			assert.Same(t, tc.value, v.Value())
			assert.Nil(t, v.Grad())
			assert.False(t, v.HasGrad())
			assert.Equal(t, tc.requiresGrad, v.RequiresGrad())
		})
	}
}

func TestNewScalar(t *testing.T) {
	t.Run("float32", testNewScalar[float32])
	t.Run("float64", testNewScalar[float64])
}

func testNewScalar[T mat.DType](t *testing.T) {
	v := NewScalar[T](42)
	require.NotNil(t, v)
	assert.Empty(t, v.Name())
	mattest.AssertMatrixEquals[T](t, mat.NewScalar[T](42), v.Value())
	assert.Nil(t, v.Grad())
	assert.False(t, v.HasGrad())
	assert.False(t, v.RequiresGrad())
}

func TestNewScalarWithName(t *testing.T) {
	t.Run("float32", testNewScalarWithName[float32])
	t.Run("float64", testNewScalarWithName[float64])
}

func testNewScalarWithName[T mat.DType](t *testing.T) {
	v := NewScalarWithName[T](42, "foo")
	require.NotNil(t, v)
	assert.Equal(t, "foo", v.Name())
	mattest.AssertMatrixEquals[T](t, mat.NewScalar[T](42), v.Value())
	assert.Nil(t, v.Grad())
	assert.False(t, v.HasGrad())
	assert.False(t, v.RequiresGrad())
}

func TestConstant(t *testing.T) {
	t.Run("float32", testConstant[float32])
	t.Run("float64", testConstant[float64])
}

func testConstant[T mat.DType](t *testing.T) {
	v := Constant[T](42)
	require.NotNil(t, v)
	assert.Equal(t, "42", v.Name())
	mattest.AssertMatrixEquals[T](t, mat.NewScalar[T](42), v.Value())
	assert.Nil(t, v.Grad())
	assert.False(t, v.HasGrad())
	assert.False(t, v.RequiresGrad())
}

func TestVariable_Gradients(t *testing.T) {
	t.Run("float32", testVariableGradients[float32])
	t.Run("float64", testVariableGradients[float64])
}

func testVariableGradients[T mat.DType](t *testing.T) {
	t.Run("with requires gradient true", func(t *testing.T) {
		v := NewVariable[T](mat.NewScalar[T](42), true)
		require.Nil(t, v.Grad())
		assert.False(t, v.HasGrad())

		v.AccGrad(mat.NewScalar[T](5))
		mattest.RequireMatrixEquals[T](t, mat.NewScalar[T](5), v.Grad())
		assert.True(t, v.HasGrad())

		v.AccGrad(mat.NewScalar[T](10))
		mattest.RequireMatrixEquals[T](t, mat.NewScalar[T](15), v.Grad())
		assert.True(t, v.HasGrad())

		v.ZeroGrad()
		require.Nil(t, v.Grad())
		assert.False(t, v.HasGrad())
	})

	t.Run("with requires gradient false", func(t *testing.T) {
		v := NewVariable[T](mat.NewScalar[T](42), false)
		require.Nil(t, v.Grad())
		assert.False(t, v.HasGrad())

		v.AccGrad(mat.NewScalar[T](5))
		require.Nil(t, v.Grad())
		assert.False(t, v.HasGrad())

		v.ZeroGrad()
		require.Nil(t, v.Grad())
		assert.False(t, v.HasGrad())
	})
}