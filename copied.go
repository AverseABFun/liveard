package main

func delete[S ~[]E, E any](s S, i, j int) S { // Copied from slices package to reduce memory footprint
	_ = s[i:j:len(s)] // bounds check

	if i == j {
		return s
	}

	oldlen := len(s)
	s = append(s[:i], s[j:]...)
	clear(s[len(s):oldlen]) // zero/nil out the obsolete elements, for GC
	return s
}

func index[S ~[]E, E comparable](s S, v E) int { // Copied from slices package to reduce memory footprint
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

const intSize = 32 << (^uint(0) >> 63)

func Atoi(s string) (int, error) { // Copied from strconv package to reduce memory footprint
	const fnAtoi = "Atoi"

	sLen := len(s)
	if intSize == 32 && (0 < sLen && sLen < 10) ||
		intSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.
		s0 := s
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
			if len(s) < 1 {
				return 0, syntaxError(fnAtoi, s0)
			}
		}

		n := 0
		for _, ch := range []byte(s) {
			ch -= '0'
			if ch > 9 {
				return 0, syntaxError(fnAtoi, s0)
			}
			n = n*10 + int(ch)
		}
		if s0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Slow path for invalid, big, or underscored integers.
	i64, err := ParseInt(s, 10, 0)
	if nerr, ok := err.(*NumError); ok {
		nerr.Func = fnAtoi
	}
	return int(i64), err
}

func ParseInt(s string, base int, bitSize int) (i int64, err error) {
	const fnParseInt = "ParseInt"

	if s == "" {
		return 0, syntaxError(fnParseInt, s)
	}

	// Pick off leading sign.
	s0 := s
	neg := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	var un uint64
	un, err = ParseUint(s, base, bitSize)
	if err != nil && err.(*NumError).Err != ErrRange {
		err.(*NumError).Func = fnParseInt
		err.(*NumError).Num = cloneString(s0)
		return 0, err
	}

	if bitSize == 0 {
		bitSize = intSize
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), rangeError(fnParseInt, s0)
	}
	if neg && un > cutoff {
		return -int64(cutoff), rangeError(fnParseInt, s0)
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

func ParseUint(s string, base int, bitSize int) (uint64, error) {
	const fnParseUint = "ParseUint"

	if s == "" {
		return 0, syntaxError(fnParseUint, s)
	}

	base0 := base == 0

	s0 := s
	switch {
	case 2 <= base && base <= 36:
		// valid base; nothing to do

	case base == 0:
		// Look for octal, hex prefix.
		base = 10
		if s[0] == '0' {
			switch {
			case len(s) >= 3 && lower(s[1]) == 'b':
				base = 2
				s = s[2:]
			case len(s) >= 3 && lower(s[1]) == 'o':
				base = 8
				s = s[2:]
			case len(s) >= 3 && lower(s[1]) == 'x':
				base = 16
				s = s[2:]
			default:
				base = 8
				s = s[1:]
			}
		}

	default:
		return 0, baseError(fnParseUint, s0, base)
	}

	if bitSize == 0 {
		bitSize = intSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, bitSizeError(fnParseUint, s0, bitSize)
	}

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	// Use compile-time constants for common cases.
	var cutoff uint64
	switch base {
	case 10:
		cutoff = maxUint64/10 + 1
	case 16:
		cutoff = maxUint64/16 + 1
	default:
		cutoff = maxUint64/uint64(base) + 1
	}

	maxVal := uint64(1)<<uint(bitSize) - 1

	underscores := false
	var n uint64
	for _, c := range []byte(s) {
		var d byte
		switch {
		case c == '_' && base0:
			underscores = true
			continue
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= lower(c) && lower(c) <= 'z':
			d = lower(c) - 'a' + 10
		default:
			return 0, syntaxError(fnParseUint, s0)
		}

		if d >= byte(base) {
			return 0, syntaxError(fnParseUint, s0)
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, rangeError(fnParseUint, s0)
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return maxVal, rangeError(fnParseUint, s0)
		}
		n = n1
	}

	if underscores && !underscoreOK(s0) {
		return 0, syntaxError(fnParseUint, s0)
	}

	return n, nil
}

func lower(c byte) byte {
	return c | ('x' - 'X')
}

func syntaxError(fn, str string) *NumError {
	return &NumError{fn, cloneString(str), ErrSyntax}
}

func baseError(fn, str string, base int) *NumError {
	return &NumError{fn, cloneString(str), NewError("invalid base " + Itoa(base))}
}

func bitSizeError(fn, str string, bitSize int) *NumError {
	return &NumError{fn, cloneString(str), NewError("invalid bit size " + Itoa(bitSize))}
}

func rangeError(fn, str string) *NumError {
	return &NumError{fn, cloneString(str), ErrRange}
}

const maxUint64 = 1<<64 - 1

type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat, ParseComplex)
	Num  string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *NumError) Error() string { // Modified slightly so that I didn't have to include a giant chain of functions to quote e.Num
	return "strconv." + e.Func + ": " + "parsing " + e.Num + ": " + e.Err.Error()
}

func cloneString(x string) string { return string([]byte(x)) }

var ErrSyntax = NewError("invalid syntax")
var ErrRange = NewError("value out of range")

func underscoreOK(s string) bool {
	// saw tracks the last character (class) we saw:
	// ^ for beginning of number,
	// 0 for a digit or base prefix,
	// _ for an underscore,
	// ! for none of the above.
	saw := '^'
	i := 0

	// Optional sign.
	if len(s) >= 1 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}

	// Optional base prefix.
	hex := false
	if len(s) >= 2 && s[0] == '0' && (lower(s[1]) == 'b' || lower(s[1]) == 'o' || lower(s[1]) == 'x') {
		i = 2
		saw = '0' // base prefix counts as a digit for "underscore as digit separator"
		hex = lower(s[1]) == 'x'
	}

	// Number proper.
	for ; i < len(s); i++ {
		// Digits are always okay.
		if '0' <= s[i] && s[i] <= '9' || hex && 'a' <= lower(s[i]) && lower(s[i]) <= 'f' {
			saw = '0'
			continue
		}
		// Underscore must follow digit.
		if s[i] == '_' {
			if saw != '0' {
				return false
			}
			saw = '_'
			continue
		}
		// Underscore must also be followed by digit.
		if saw == '_' {
			return false
		}
		// Saw non-digit, non-underscore.
		saw = '!'
	}
	return saw != '_'
}

func Itoa(i int) string { // Almost entirely rewritten
	if i < 10 {
		return digits[i : i+1]
	}
	var out = ""
	for i >= 1 {
		i /= 10
		println(i)
		out += digits[i%10 : (i%10)+1]
	}
	return out
}

func getDigits(i int) []int {
	var out = []int{}
	for i >= 1 {
		i /= 10
		out = append(out, i%10)
	}
	return out
}

var digits = "0123456789"

func TrimSuffix(s, suffix string) string {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func NewError(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
