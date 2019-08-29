package runtime2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestObj struct{}

func TestTypeString(t *testing.T) {
	var (
		assert = require.New(t)
		table  = []struct {
			x interface{}
			s string
		}{
			{int(0), "int"},
			{uint(0), "uint"},
			{int8(0), "int8"},
			{float32(0), "float32"},
			{float64(0), "float64"},
			{TestObj{}, "runtime2.TestObj"},
		}
	)

	for _, v := range table {
		assert.Equal(v.s, TypeString(v.x))
	}
}

func TestHash(t *testing.T) {
	var table = []interface{}{
		int(0), int(1),
	}
	for _, v := range table {
		Hash(v)
	}
}

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hash(i)
	}
}
