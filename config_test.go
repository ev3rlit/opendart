package opendart

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequiresAPIKey(t *testing.T) {
	client, err := New(Config{})
	require.Error(t, err)
	assert.Nil(t, client)
	assertOopsContext(t, err, map[string]any{"field": "APIKey"})
}

func TestNewTrimsBaseURL(t *testing.T) {
	client, err := New(Config{APIKey: "test-key"}, WithBaseURL("https://example.com/"))
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://example.com", client.resty.BaseURL)
}
