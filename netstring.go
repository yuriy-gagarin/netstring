package netstring

import (
	"fmt"
	"strconv"
)

const (
	TokenLength = iota
	TokenSeparator
	TokenData
	TokenEnd
)

func SplitNetstring(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	i := 0
	length := 0
	expectedToken := TokenLength
	buf := make([]byte, 0)

	for i < len(data) {
		switch expectedToken {

		case TokenLength:

			// skip empty start, if any
			for '0' > data[i] || data[i] > '9' {
				i++
				if i >= len(data) {
					return 0, nil, nil
				}
			}

			for '0' <= data[i] && data[i] <= '9' {
				length = length*10 + int(data[i]) - 48
				i++
			}

			if length == 0 {
				// data is broken, try to continue
				return i, nil, fmt.Errorf("Broken data: no length, instead got %v", data[i])
			}

			expectedToken = TokenSeparator

		case TokenSeparator:

			if data[i] != ':' {
				return 0, nil, fmt.Errorf("Broken data: missing separator, instead got %v", data[0])
			}

			i++
			expectedToken = TokenData

		case TokenData:

			size := min(length, len(data)-i)
			buf = append(data[:0:0], data[i:i+size]...)
			length -= size
			if length < 1 {
				expectedToken = TokenEnd
			}
			i += size

		case TokenEnd:

			if data[i] != ',' {
				// data is broken, try to continue
				return i, nil, fmt.Errorf("Broken data: missing end, instead got %v", data[i])
			}
			i++

			if len(buf) == 0 {
				return i, nil, nil
			}

			return i, buf, nil

		}
	}

	return 0, nil, nil
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Encode(data []byte) []byte {
	length := strconv.FormatInt(int64(len(data)), 10)
	return []byte(length + ":" + string(data) + ",")
}
