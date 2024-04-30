package filegenerator

import (
	"errors"
	"fmt"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/tools/io"
)

type FileWriter interface {
	WriteFile(path, content string) error
}

type DefaultContentProvider interface {
	GetDefaultContent(args types.KeyValueArgs) (string, error)
}

type FileGeneratorService struct {
	kind                   string
	fileWriter             FileWriter
	defaultContentProvider DefaultContentProvider
}

func NewFileGeneratorService(kind string, fileSystem FileWriter, defaultContentProvider DefaultContentProvider) *FileGeneratorService {
	return &FileGeneratorService{
		kind:                   kind,
		fileWriter:             fileSystem,
		defaultContentProvider: defaultContentProvider,
	}
}

func (s *FileGeneratorService) GenerateFile(out io.Out, path string, args types.KeyValueArgs) error {
	defaultContent, err := s.defaultContentProvider.GetDefaultContent(args)
	if err != nil {
		return errors.Join(ErrGettingDefaultContent, err)
	}

	if err := s.writeFile(defaultContent, path); err != nil {
		return err
	}

	out.Write(fmt.Sprintf("Generated a blank %s file: %s\n", s.kind, path))

	return nil
}

func (s *FileGeneratorService) writeFile(content, path string) error {
	if err := s.fileWriter.WriteFile(path, content); err != nil {
		return fmt.Errorf("%w %s: %w", ErrWritingFile, path, err)
	}

	return nil
}
