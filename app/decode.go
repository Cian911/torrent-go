package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// decodeAt decodes ONE value starting at s[i].
// It returns (value, nextIndexAfterValue, error).
func decodeAt(s string, i int) (any, int, error) {
	if i < 0 || i >= len(s) {
		return nil, i, fmt.Errorf("unexpected end of input at %d", i)
	}

	switch {
	case s[i] == 'i':
		return decodeIntAt(s, i)
	case s[i] == 'l':
		return decodeListAt(s, i)
	case unicode.IsDigit(rune(s[i])):
		return decodeStringAt(s, i)
	default:
		return nil, i, fmt.Errorf("invalid bencode token %q at index %d", s[i], i)
	}
}

func decodeIntAt(s string, i int) (any, int, error) {
	// Format: i<base10>e
	if s[i] != 'i' {
		return nil, i, errors.New("decodeIntAt called on non-int")
	}

	j := i + 1
	if j >= len(s) {
		return nil, i, fmt.Errorf("unterminated int at %d", i)
	}

	// Find terminating 'e'
	end := j
	for end < len(s) && s[end] != 'e' {
		end++
	}
	if end >= len(s) {
		return nil, i, fmt.Errorf("unterminated int at %d", i)
	}

	numStr := s[j:end]
	if numStr == "" {
		return nil, i, fmt.Errorf("empty int at %d", i)
	}

	// Bencode integer rules:
	// - "i0e" valid
	// - "-0" invalid
	// - no leading zeros: "03" invalid, "-03" invalid
	if numStr == "-0" {
		return nil, i, fmt.Errorf("invalid int -0 at %d", i)
	}
	if numStr[0] == '0' && len(numStr) > 1 {
		return nil, i, fmt.Errorf("invalid leading zero in int %q at %d", numStr, i)
	}
	if numStr[0] == '-' {
		if len(numStr) == 1 {
			return nil, i, fmt.Errorf("invalid int %q at %d", numStr, i)
		}
		if numStr[1] == '0' && len(numStr) > 2 {
			return nil, i, fmt.Errorf("invalid leading zero in int %q at %d", numStr, i)
		}
	}

	v, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return nil, i, fmt.Errorf("invalid int %q at %d: %w", numStr, i, err)
	}

	return v, end + 1, nil
}

func decodeStringAt(s string, i int) (any, int, error) {
	// Format: <len>:<bytes>
	j := i
	for j < len(s) && unicode.IsDigit(rune(s[j])) {
		j++
	}
	if j == i {
		return nil, i, fmt.Errorf("string length missing at %d", i)
	}
	if j >= len(s) || s[j] != ':' {
		return nil, i, fmt.Errorf("expected ':' after string length at %d", i)
	}

	n, err := strconv.Atoi(s[i:j])
	if err != nil || n < 0 {
		return nil, i, fmt.Errorf("invalid string length %q at %d", s[i:j], i)
	}

	start := j + 1
	end := start + n
	if end > len(s) {
		return nil, i, fmt.Errorf("string overruns input at %d (need %d bytes)", i, n)
	}

	return s[start:end], end, nil
}

func decodeListAt(s string, i int) (any, int, error) {
	// Format: l<item><item>...e
	if s[i] != 'l' {
		return nil, i, errors.New("decodeListAt called on non-list")
	}

	out := make([]any, 0)
	idx := i + 1

	for {
		if idx >= len(s) {
			return nil, i, fmt.Errorf("unterminated list at %d", i)
		}
		if s[idx] == 'e' {
			return out, idx + 1, nil
		}

		val, next, err := decodeAt(s, idx)
		if err != nil {
			return nil, i, err
		}

		out = append(out, val)
		idx = next
	}
}

