package defaultcr

import "errors"

var (
	ErrGeneratingDefaultCRFile = errors.New("error generating default CR file")
	ErrWritingDefaultCRFile    = errors.New("error writing default CR file")
)
