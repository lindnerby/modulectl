package manifest

import "fmt"

type FileSystem interface {
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

func (s *ManifestService) GetDefaultManifestContent() string {
	return `# This file holds the Manifest of your module, encompassing all resources installed in the cluster once the module is activated.
# It should include the Custom Resource Definition for your module's default CustomResource, if it exists.

`
}

func (s *ManifestService) WriteManifestFile(content, path string) error {
	if err := s.fileSystem.WriteFile(path, content); err != nil {
		return fmt.Errorf("%w %s: %w", ErrWritingManifestFile, path, err)
	}

	return nil
}
