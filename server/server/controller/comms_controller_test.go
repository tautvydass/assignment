package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCommsController_Close(t *testing.T) {
	c := NewCommsController(zap.NewNop())
	require.NoError(t, c.Close())
}
