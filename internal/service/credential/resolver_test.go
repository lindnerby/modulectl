package credential_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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

//
//func Test_GetCredentials_ReturnUserPasswordWhenGiven(t *testing.T) {
//	userPasswordCreds := "user1:pass1"
//	creds := registry.GetCredentials(cpi.DefaultContext(), false, userPasswordCreds, "ghcr.io")
//
//	require.Equal(t, "user1", creds.GetProperty("username"))
//	require.Equal(t, "pass1", creds.GetProperty("password"))
//}
//
//func Test_GetCredentials_ReturnNilWhenInsecure(t *testing.T) {
//	creds := registry.GetCredentials(cpi.DefaultContext(), true, "", "ghcr.io")
//
//	require.Equal(t, credentials.NewCredentials(nil), creds)
//}
//
//func Test_UserPass_ReturnsCorrectUsernameAndPassword(t *testing.T) {
//	user, pass := registry.ParseUserPass("user1:pass1")
//	require.Equal(t, "user1", user)
//	require.Equal(t, "pass1", pass)
//}
//
//func Test_UserPass_ReturnsCorrectUsername(t *testing.T) {
//	user, pass := registry.ParseUserPass("user1:")
//	require.Equal(t, "user1", user)
//	require.Empty(t, pass)
//}
//
//func Test_UserPass_ReturnsCorrectPassword(t *testing.T) {
//	user, pass := registry.ParseUserPass(":pass1")
//	require.Empty(t, user)
//	require.Equal(t, "pass1", pass)
//}
