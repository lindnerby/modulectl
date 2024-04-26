package stub

import "errors"

type FileExistsStub struct{}

func (*FileExistsStub) FileExists(_ string) (bool, error) {
	return true, nil
}

type FileDoesNotExistStub struct{}

func (*FileDoesNotExistStub) FileExists(_ string) (bool, error) {
	return false, nil
}

var ErrSomeOSError = errors.New("some OS error")

type ErrorStub struct{}

func (*ErrorStub) FileExists(_ string) (bool, error) {
	return false, ErrSomeOSError
}
