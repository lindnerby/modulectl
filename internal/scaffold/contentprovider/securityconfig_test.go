package contentprovider_test

import (
	"testing"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/contentprovider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetDefaultContent_ReturnsError_WhenArgsIsNil(t *testing.T) {
	svc := contentprovider.NewSecurityConfigContentProvider(&objectToYAMLConverterStub{})

	result, err := svc.GetDefaultContent(nil)

	require.ErrorIs(t, err, contentprovider.ErrInvalidArg)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "args")
}

func Test_GetDefaultContent_ReturnsError_WhenModuleNameArgMissing(t *testing.T) {
	svc := contentprovider.NewSecurityConfigContentProvider(&objectToYAMLConverterStub{})

	result, err := svc.GetDefaultContent(types.KeyValueArgs{})

	require.ErrorIs(t, err, contentprovider.ErrMissingArg)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "moduleName")
}

func Test_GetDefaultContent_ReturnsError_WhenModuleNameArgIsEmpty(t *testing.T) {
	svc := contentprovider.NewSecurityConfigContentProvider(&objectToYAMLConverterStub{})

	result, err := svc.GetDefaultContent(types.KeyValueArgs{contentprovider.ArgModuleName: ""})

	require.ErrorIs(t, err, contentprovider.ErrInvalidArg)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "moduleName")
}

func Test_GetDefaultContent_ReturnsConvertedContent(t *testing.T) {
	svc := contentprovider.NewSecurityConfigContentProvider(&objectToYAMLConverterStub{})

	result, err := svc.GetDefaultContent(types.KeyValueArgs{contentprovider.ArgModuleName: "module-name"})

	require.NoError(t, err)
	assert.Equal(t, convertedContent, result)
}

// ***************
// Test Stubs
// ***************

type objectToYAMLConverterStub struct{}

const convertedContent = "content"

func (o *objectToYAMLConverterStub) ConvertToYaml(obj interface{}) string {
	return convertedContent
}
