package manifest

import (
	"fmt"

	"github.com/kyma-project/modulectl/tools/io"
)

type FileSystem interface {
	FileExists(path string) (bool, error)
	WriteFile(path, content string) error
}

type ManifestService struct {
	fileSystem FileSystem
}

func NewManifestService(fileSystemUtil FileSystem) *ManifestService {
	return &ManifestService{
		fileSystem: fileSystemUtil,
	}
}

func (s *ManifestService) GenerateManifestFile(out io.Out, manifestFilePath string) error {
	defaultCRFileExists, err := s.fileSystem.FileExists(manifestFilePath)
	if err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingManifestFile, manifestFilePath, err)
	}

	if defaultCRFileExists {
		out.Write(fmt.Sprintf("The manifest file already exists, reusing: %s\n", manifestFilePath))
		return nil
	}

	if err := s.writeFile(s.getDefaultContent(), manifestFilePath); err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingManifestFile, manifestFilePath, err)
	}

	out.Write(fmt.Sprintf("Generated a blank manifest file: %s\n", manifestFilePath))

	return nil
}

func (s *ManifestService) getDefaultContent() string {
	return `# This file holds the Manifest of your module, encompassing all resources installed in the cluster once the module is activated.
# It should include the Custom Resource Definition for your module's default CustomResource, if it exists.

`
}

func (s *ManifestService) writeFile(content, path string) error {
	if err := s.fileSystem.WriteFile(path, content); err != nil {
		return fmt.Errorf("%w %s: %w", ErrWritingManifestFile, path, err)
	}

	return nil
}
