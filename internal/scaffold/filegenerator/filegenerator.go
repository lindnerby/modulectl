package filegenerator

import (
	"errors"
	"fmt"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/tools/io"
)

type FileSystem interface {
	FileExists(path string) (bool, error)
	WriteFile(path, content string) error
}

type DefaultContentProvider interface {
	GetDefaultContent(args types.KeyValueArgs) (string, error)
}

type FileGeneratorService struct {
	kind                   string
	fileSystem             FileSystem
	defaultContentProvider DefaultContentProvider
}

func NewFileGeneratorService(kind string, fileSystem FileSystem, defaultContentProvider DefaultContentProvider) *FileGeneratorService {
	return &FileGeneratorService{
		kind:                   kind,
		fileSystem:             fileSystem,
		defaultContentProvider: defaultContentProvider,
	}
}

func (s *FileGeneratorService) GenerateFile(out io.Out, path string, args types.KeyValueArgs) error {
	fileExists, err := s.fileSystem.FileExists(path)
	if err != nil {
		return fmt.Errorf("%w %s: %w", ErrCheckingFileExistence, path, err)
	}

	if fileExists {
		out.Write(fmt.Sprintf("The %s file already exists, reusing: %s\n", s.kind, path))
		return nil
	}

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
	if err := s.fileSystem.WriteFile(path, content); err != nil {
		return fmt.Errorf("%w %s: %w", ErrWritingFile, path, err)
	}

	return nil
}
