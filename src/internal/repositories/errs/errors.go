package errs

import "errors"

var (
	ErrNotFound = errors.New("[-] Repository error! Объекта нет в базе данных")
)
