package opendart

import (
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertOopsContext(t *testing.T, err error, expected map[string]any) {
	t.Helper()

	oopsErr, ok := oops.AsOops(err)
	require.True(t, ok, "expected oops error context")

	context := oopsErr.Context()
	for key, value := range expected {
		assert.Equal(t, value, context[key], "oops context %q", key)
	}
}
