package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	strRune := []rune(str)
	var isSkipNextRune bool

	var stringBuilder strings.Builder

	for i, r := range strRune {
		if isFirstRuneDigit(i, r) {
			return throwError()
		}

		if isNextRuneNumber(i, r, strRune) {
			return throwError()
		}

		if isSkipNextRune {
			isSkipNextRune = false
			continue
		}

		if isNeededRepeatNewLine(i, r, strRune) {
			stringBuilder.WriteString(strings.Repeat(string('\n'), runeToInt(strRune[i+1])))
			isSkipNextRune = true
			continue
		}

		if isNextSymbolZero(i, strRune) {
			isSkipNextRune = true
			continue
		}

		if isNeededRepeatLetter(i, r, strRune) {
			repeatableCount := runeToInt(r) - 1
			letter := strRune[i-1]
			stringBuilder.WriteString(strings.Repeat(string(letter), repeatableCount))
			continue
		}

		stringBuilder.WriteRune(r)

	}

	return stringBuilder.String(), nil
}

func isNeededRepeatLetter(i int, r rune, strRune []rune) bool {
	return unicode.IsDigit(r) && unicode.IsLetter(strRune[i-1])
}

func isNeededRepeatNewLine(i int, r rune, strRune []rune) bool {
	return r == '\n' && isNextSymbolDigit(i, strRune)
}

func isFirstRuneDigit(i int, r rune) bool {
	return unicode.IsDigit(r) && i == 0
}

func isNextRuneNumber(i int, r rune, strRune []rune) bool {
	if len(strRune) <= i+1 {
		return false
	}
	return unicode.IsDigit(r) && unicode.IsDigit(strRune[i+1])
}

func isNextSymbolDigit(i int, strRune []rune) bool {
	if len(strRune) <= i+1 {
		return false
	}
	nextSymbol := strRune[i+1]

	return unicode.IsDigit(nextSymbol)
}

func isNextSymbolZero(i int, strRune []rune) bool {
	return isNextSymbolDigit(i, strRune) && isZero(strRune[i+1])
}

func isZero(r rune) bool {
	return runeToInt(r) == 0
}

func runeToInt(r rune) int {
	return int(r) - '0'
}

func throwError() (string, error) {
	return "", ErrInvalidString
}
