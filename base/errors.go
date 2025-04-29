package base

import "errors"

var (
	ErrInvalidStringParse      = errors.New("String cannot be parsed to this datatype")
	ErrInvalidTimeString       = errors.New("String cannot be parsed to time.Duration")
	ErrExtractingTableData     = errors.New("Error when extracting table data")
	ErrNoGoodByNameFound       = errors.New("Error no good with that name found")
	ErrFailedToCreateDirectory = errors.New("Failed to create directory")
	ErrFileNotExists           = errors.New("File does not exist")
	ErrFailedToReadFile        = errors.New("Failed to read file")
	ErrFailedJSONParse         = errors.New("Failed to parse JSON")
	ErrFailedToWriteFile       = errors.New("Failed to write file")
)
