package filegenerator

import "errors"

var (
	ErrCheckingFileExistence = errors.New("error checking file existence")
	ErrGettingDefaultContent = errors.New("error getting default content")
	ErrWritingFile           = errors.New("error writing file")
)
