package compdescbuilder

import (
	"fmt"
	"ocm.software/ocm/api/ocm/compdesc"
	ocmv1 "ocm.software/ocm/api/ocm/compdesc/meta/v1"
)

type Service struct {
}

// NewService
func NewService() Service {

}

type CompDescBuilder struct {
	descriptor *compdesc.ComponentDescriptor
}

const (
	versionV1, versionV2 = "v1", "v2"
	providerName         = "kyma-project.io"
	labelKey             = providerName + "/built-by"
	labelValue           = "modulectl"
)

func New() *CompDescBuilder {
	componentDescriptor := &compdesc.ComponentDescriptor{}
	componentDescriptor.SetName(moduleName)
	componentDescriptor.SetVersion(moduleVersion)
	componentDescriptor.Metadata.ConfiguredVersion = versionV2

	providerLabel, err := ocmv1.NewLabel(labelKey, labelValue, ocmv1.WithVersion(versionV1))
	if err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}
	componentDescriptor.Provider = ocmv1.Provider{Name: providerName, Labels: ocmv1.Labels{*providerLabel}}

	compdesc.DefaultResources(componentDescriptor)
	if err = compdesc.Validate(componentDescriptor); err != nil {
		return nil, fmt.Errorf("failed to validate component descriptor: %w", err)
	}

	return &CompDescBuilder{
		descriptor: componentDescriptor,
	}
}

func (c *CompDescBuilder) WithName(name string) *CompDescBuilder {
	c.descriptor.SetName(name)
	return c
}

func (c *CompDescBuilder) Build() *compdesc.ComponentDescriptor {
	return c.descriptor
}
