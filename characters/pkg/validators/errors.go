package validators

import (
	"errors"
)

// Error represents a validation error
type Error struct {
	error
}

var (
	// ErrorCharacterNameEmpty will be sent if the request contains an empty character name
	// Code: 10001
	ErrorCharacterNameEmpty = Error{errors.New("the provided character name is an empty string")}

	// ErrorCharacterNameTooSmall will be sent if the request contains a character name of length < MinRunesAllowedInACharacterName
	// Code: 10002
	ErrorCharacterNameTooSmall = Error{errors.New("the provided character name is too small")}

	// ErrorCharacterNameInvalid will be sent if the request contains an invalid character name
	// Code: 10003
	ErrorCharacterNameInvalid = Error{errors.New("the provided character name is invalid")}

	// ErrorCharacterNameIsOnlyWhiteSpace will be sent if the request contains a name that consists of only whitespaces
	// Code: 10004
	ErrorCharacterNameIsOnlyWhiteSpace = Error{errors.New("the provided character name consists only of whitespaces")}

	// ErrorCharacterNameTooBig will be sent if the request contains a character name of length > MaxRunesAllowedInACharacterName
	// Code: 10005
	ErrorCharacterNameTooBig = Error{errors.New("the provided character name is too big")}

	// ErrorCharacterWordTooBig will be sent if the request contains a word with length > MaxRunesAllowedInACharacterNameWord in the character name
	// Code: 10006
	ErrorCharacterWordTooBig = Error{errors.New("the provided character name has a word too big")}

	// ErrorCharacterWordTooSmall will be sent if the request contains a word with length < MinRunesAllowedInACharacterNameWord in the character name
	// Code: 10007
	ErrorCharacterWordTooSmall = Error{errors.New("the provided character name has a word too small")}
)

// Code will return the code of the error
func (e Error) Code() int {
	switch e {
	case ErrorCharacterNameEmpty:
		return 10001
	case ErrorCharacterNameTooSmall:
		return 10002
	case ErrorCharacterNameInvalid:
		return 10003
	case ErrorCharacterNameIsOnlyWhiteSpace:
		return 10004
	case ErrorCharacterNameTooBig:
		return 10005
	case ErrorCharacterWordTooBig:
		return 10006
	case ErrorCharacterWordTooSmall:
		return 10007
	default:
		return 0
	}
}
