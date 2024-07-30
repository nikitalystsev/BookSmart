package errsService

import "errors"

var (
	ErrLibCardAlreadyExist  = errors.New("[!] libCardService error! LibCard already exists")
	ErrLibCardDoesNotExists = errors.New("[!] libCardService error! LibCard does not exist")
	ErrLibCardIsValid       = errors.New("[!] libCardService error! LibCard is valid")
	ErrLibCardIsInvalid     = errors.New("[!] libCardService error! LibCard is invalid")
)
