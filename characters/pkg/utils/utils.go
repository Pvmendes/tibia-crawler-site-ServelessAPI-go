package utils

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

const (
	htmlTagStart = 60 // Unicode `<`
	htmlTagEnd   = 62 // Unicode `>`
)

func TibiaDataDatetime(date string) string {
	//TODO: Normalization needs to happen above this layer
	date = norm.NFKC.String(date)

	var (
		returnDate time.Time
		err        error
	)

	// If statement to determine if date string is filled or empty
	if date == "" {
		// The string that should be returned is the current timestamp
		returnDate = time.Now()
	} else {
		// timezone use in html: CET/CEST
		loc, _ := time.LoadLocation("Europe/Berlin")

		// format used in datetime on html: Jan 02 2007, 19:20:30 CET
		formatting := "Jan 02 2006, 15:04:05 MST"

		// parsing html in time with location set in loc
		returnDate, err = time.ParseInLocation(formatting, date, loc)

		if err != nil {
			fmt.Println(err)
		}
	}

	// Return of formatted date and time string to functions..
	return returnDate.UTC().Format(time.RFC3339)
}

func TibiaDataDate(date string) string {
	// removing weird spacing and comma
	date = SanitizeStrings(strings.ReplaceAll(date, ",", ""))

	// var time parser
	var tmpDate time.Time
	var err error

	// date formats to parse
	dateFormats := map[string][]string{
		"YearMonthDay": {"January 2 2006", "Jan 02 2006"},
		"YearMonth":    {"January 2006", "Jan 2006", "2006-01", "01/06"},
	}

	for _, layout := range dateFormats["YearMonthDay"] {
		tmpDate, err = time.Parse(layout, date)
		if err == nil {
			// If the parse succeeds, format the date as "2006-01-02"
			return tmpDate.UTC().Format("2006-01-02")
		}
	}

	for _, layout := range dateFormats["YearMonth"] {
		tmpDate, err = time.Parse(layout, date)
		if err == nil {
			// If the parse succeeds, format the date as "2006-01"
			return tmpDate.Format("2006-01")
		}
	}

	return tmpDate.UTC().Format("2006-01-02")
}

func SanitizeStrings(data string) string {
	// replaces weird \u00A0 string to real space
	data = strings.ReplaceAll(data, "\u00A0", " ")
	// replaces weird \u0026 string to amp (&)
	data = strings.ReplaceAll(data, "\u0026", "&")
	// returning string unescaped
	return SanitizeEscapedString(data)
}

func SanitizeEscapedString(data string) string {
	return html.UnescapeString(data)
}

func StringToInteger(data string) int {
	str := strings.ReplaceAll(data, ",", "")

	returnData, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("[warning] TibiaDataStringToInteger: couldn't convert string into int. error: %s", err)
	}

	return returnData
}

func HTMLRemoveLinebreaks(data string) string {
	return strings.ReplaceAll(data, "\n", "")
}

// RemoveHtmlTag replaces all HTML tags with an empty string
//
// Algo from:
// https://stackoverflow.com/questions/55036156/how-to-replace-all-html-tag-with-empty-string-in-golang
func RemoveHtmlTag(s string) string {
	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character and we are not in an HTML tag, save it.
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()
	return s
}
