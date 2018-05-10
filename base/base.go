package main

import (
	"fmt"
	"math"
	"path"
	"runtime"
)

func New(base int, val float64) Number {
	var a Number
	a.Add = func(n Number) {

	}
}

type Number float64

func (number Number) Add(n Number) {

}

type BaseEncoderDecoder interface {
	Encode(int) (string, error)
	Decode(string) (int, error)
}

type Base struct {
	base []rune
	val  []rune
}

func (base Base) Encode(cur float64) (string, error) {
	// Define our available "digits"
	digits := Base{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z',
		// Only URL-safe special characters should be used
		// Underscore removed because it's hard to see in an underlined link
		'-', '~', '\'', '!', '*', '(', ')',
	}
	chars := []rune{}
	for {
		if cur <= 0 {
			break
		}
		mod := math.Mod(cur, float64(len(base)))
		cur = math.Floor(cur / float64(len(base)))
		chars = append(chars, base[int(mod)])
	}
	return string(chars), nil
}

/*
Decode will convert a case-sensitive base-69 value into an integer.
Decoding is 54% faster than encoding.
*/
func (base Base) Decode(encoded string) (float64, error) {
	var retVal float64

	// Define our available "digits"
	tmp := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z',
		// Only URL-safe special characters should be used
		// Underscore removed because it's hard to see in an underlined link
		'-', '~', '\'', '!', '*', '(', ')',
	}
	digits := map[rune]int{}
	for k, v := range tmp {
		digits[v] = k
	}

	runes := []rune(encoded)
	for k, v := range runes {
		retVal += float64(digits[runes[v]]) * math.Pow(float64(len(digits)), float64(len(runes)-(k+1)))
	}
	return retVal, nil
}

func (base Base) Format(s fmt.State, verb rune) {
	switch verb {
	case 'd':
		fmt.Fprintf(s, "%d", base.Decode())
	case 'v':
		switch {
		case s.Flag('+'):
			// Detailed stack trace
			for k, err := range base {
				msg, ok := ErrMsg[err.code]
				if !ok {
					msg = ErrMsg[ErrUnspecified]
				}
				fmt.Fprintf(s, "\n%s", fmt.Sprintf(`(%d) %s %s:%d %s
	Mesg: %s
	Code: %d
	Text: %s
	Http: %d`,
					k,
					err.Error(),
					path.Base(err.caller.File),
					err.caller.Line,
					runtime.FuncForPC(err.caller.Pc).Name(),
					err.msg,
					err.code,
					msg,
					ErrHTTPStatus[err.code],
				))
			}

		case s.Flag('#'):
			// Condensed stack trace
			for k, err := range stack {
				code := ""

				if err.code != ErrUnspecified {
					code = fmt.Sprintf(" (#%d)", err.code)
				}
				msg, ok := ErrMsg[err.code]
				if !ok {
					msg = ErrMsg[ErrUnspecified]
				}
				fmt.Fprintf(s, "%s", fmt.Sprintf(`(%d) %s %s:%d %s - %s%s '%s'; `,
					k,
					err.Error(),
					path.Base(err.caller.File),
					err.caller.Line,
					runtime.FuncForPC(err.caller.Pc).Name(),
					err.msg,
					code,
					msg,
				))
			}
		default:
			// Default struct output
			fmt.Fprintf(s, "%s", stack.Error())
		}
	case 's':
		// Simple error messages
		fmt.Fprintf(s, "%s", stack.Error())
	}
}
