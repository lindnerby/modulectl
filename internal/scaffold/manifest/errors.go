package manifest

import "errors"

var (
	ErrGeneratingManifestFile = errors.New("error generating manifest file")
	ErrWritingManifestFile    = errors.New("error writing manifest file")
)
