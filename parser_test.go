package quake_log_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser("test/test.log")
	games, err := p.Parse()
	require.NoError(t, err)
	assert.NotEmpty(t, games)
}
