package credential_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/kyma-project/modulectl/internal/service/credential"
)

func TestResolveCredentials_WhenCalledWithInvalidUsernamePasswordFormats_ReturnsError(t *testing.T) {
	_, err := credential.ResolveCredentials(nil, "invalidFormat", "")
	require.ErrorIs(t, err, credential.ErrInvalidCredentialsFormat)

	_, err = credential.ResolveCredentials(nil, ":", "")
	require.ErrorIs(t, err, credential.ErrInvalidCredentialsFormat)

	_, err = credential.ResolveCredentials(nil, ": ", "")
	require.ErrorIs(t, err, credential.ErrInvalidCredentialsFormat)

	_, err = credential.ResolveCredentials(nil, " :", "")
	require.ErrorIs(t, err, credential.ErrInvalidCredentialsFormat)

	_, err = credential.ResolveCredentials(nil, " : ", "")
	require.ErrorIs(t, err, credential.ErrInvalidCredentialsFormat)
}
