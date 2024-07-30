package errsService

import "errors"

var (
	ErrReaderAlreadyExist             = errors.New("[!] readerService error! Reader already exists")
	ErrEmptyReaderFio                 = errors.New("[!] readerService error! Empty Reader fio")
	ErrEmptyReaderPhoneNumber         = errors.New("[!] readerService error! Empty Reader phoneNumber")
	ErrInvalidReaderPhoneNumberLen    = errors.New("[!] readerService error! Invalid Reader phoneNumber len")
	ErrInvalidReaderPhoneNumberFormat = errors.New("[!] readerService error! Invalid Reader phoneNumber format")
	ErrInvalidReaderAge               = errors.New("[!] readerService error! Invalid Reader age")
	ErrReaderDoesNotExists            = errors.New("[!] readerService error! Reader does not exist")
)
