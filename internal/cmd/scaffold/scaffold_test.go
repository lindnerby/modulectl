package scaffold_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/cmd/scaffold"
	"github.com/kyma-project/modulectl/internal/testutils/builder"
)

func Test_RunScaffold_ReturnsError_WhenOutIsNil(t *testing.T) {
	opts := builder.NewScaffoldOptionsBuilder().WithOut(nil).Build()

	result := scaffold.RunScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Out")
}

func Test_RunScaffold_Succeeds(t *testing.T) {
	opts := builder.NewScaffoldOptionsBuilder().Build()

	result := scaffold.RunScaffold(opts)

	require.NoError(t, result)
}
