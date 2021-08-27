package chainer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func foo(v int) int {
	return v + 2
}
func bar(v int) int {
	return v * 2
}

func TestChainer(t *testing.T) {
	require.Equal(t, 18, Wrap(7).Then(foo).Then(bar).MustValue())
}
