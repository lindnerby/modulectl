package defaultcr

import (
	"fmt"

	"github.com/kyma-project/modulectl/tools/io"
)

type FileSystem interface {
	FileExists(path string) (bool, error)
	WriteFile(path, content string) error
}

type DefaultCRService struct {
	fileSystem FileSystem
}

func NewDefaultCRService(fileSystemUtil FileSystem) *DefaultCRService {
	return &DefaultCRService{
		fileSystem: fileSystemUtil,
	}
}

func (s *DefaultCRService) GenerateDefaultCRFile(out io.Out, defaultCRFilePath string) error {
	defaultCRFileExists, err := s.fileSystem.FileExists(defaultCRFilePath)
	if err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingDefaultCRFile, defaultCRFilePath, err)
	}

	if defaultCRFileExists {
		out.Write(fmt.Sprintf("The default CR file already exists, reusing: %s\n", defaultCRFilePath))
		return nil
	}

	if err := s.writeFile(s.getDefaultContent(), defaultCRFilePath); err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingDefaultCRFile, defaultCRFilePath, err)
	}

	out.Write(fmt.Sprintf("Generated a blank default CR file: %s\n", defaultCRFilePath))

	return nil
}

func (s *DefaultCRService) getDefaultContent() string {
	return `# This is the file that contains the defaultCR for your module, which is the Custom Resource that will be created upon module enablement.
	# Make sure this file contains *ONLY* the Custom Resource (not the Custom Resource Definition, which should be a part of your module manifest)

`
}

func (s *DefaultCRService) writeFile(content, path string) error {
	if err := s.fileSystem.WriteFile(path, content); err != nil {
		return fmt.Errorf("%w %s: %w", ErrWritingDefaultCRFile, path, err)
	}

	return nil
}
