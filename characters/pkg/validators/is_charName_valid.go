package validators

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// characterNameRegex is used to check if the character name provided is valid
	// nowadays Tibia only accepts a-zA-Z, but we have to consider old names
	characterNameRegex = regexp.MustCompile(`[^\s'\p{L}\-\.\+]`)
)

func IsCharNameValid(charName string) error {

	lenName := utf8.RuneCountInString(charName)

	switch {
	case lenName == 0: // Name is an empty string
		return ErrorCharacterNameEmpty
	case lenName < MinRunesAllowedInACharacterName: // Name is too small
		return ErrorCharacterNameTooSmall
	case lenName > MaxRunesAllowedInACharacterName: // Name is too big
		return ErrorCharacterNameTooBig
	}

	if strings.TrimSpace(charName) == "" {
		return ErrorCharacterNameIsOnlyWhiteSpace
	}

	strs := strings.Fields(charName)
	for _, str := range strs {
		if utf8.RuneCountInString(str) > MaxRunesAllowedInACharacterNameWord {
			return ErrorCharacterWordTooBig
		}

		if utf8.RuneCountInString(str) < MinRunesAllowedInACharacterNameWord {
			return ErrorCharacterWordTooSmall
		}
	}

	matched := characterNameRegex.MatchString(charName)
	if matched {
		return ErrorCharacterNameInvalid
	}

	return nil
}
