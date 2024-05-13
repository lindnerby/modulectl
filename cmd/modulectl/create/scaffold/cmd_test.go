package scaffold_test

import (
	"testing"

	scaffoldcmd "github.com/kyma-project/modulectl/cmd/modulectl/create/scaffold"
	scaffoldsvc "github.com/kyma-project/modulectl/internal/service/scaffold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCmd_ReturnsError_WhenScaffoldServiceIsNil(t *testing.T) {
	_, err := scaffoldcmd.NewCmd(nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "scaffoldService")
}

func Test_NewCmd_Succceeds(t *testing.T) {
	svc := &scaffoldServiceStub{called: false}
	_, err := scaffoldcmd.NewCmd(svc)

	require.NoError(t, err)
	require.False(t, svc.called)
}

func Test_Execute_CallsScaffoldService(t *testing.T) {
	svc := &scaffoldServiceStub{called: false}
	cmd, _ := scaffoldcmd.NewCmd(svc)

	err := cmd.Execute()

	require.NoError(t, err)
	require.True(t, svc.called)
}

// ***************
// Test Stubs
// ***************

type scaffoldServiceStub struct {
	called bool
}

func (s *scaffoldServiceStub) CreateScaffold(_ scaffoldsvc.Options) error {
	s.called = true
	return nil
}
