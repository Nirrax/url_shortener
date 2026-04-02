package utils

import (
	"fmt"
	"math"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type Encoder interface {
	Encode(n uint64) string
	Decode(s string) (uint64, error)
}

type Base62Encoder struct{}

func NewBase62Encoder() *Base62Encoder {
	return &Base62Encoder{}
}

func (e *Base62Encoder) Encode(num uint64) string {
	if num == 0 {
		return "0"
	}
	result := []byte{}
	for num > 0 {
		result = append([]byte{charset[num%62]}, result...)
		num /= 62
	}
	return string(result)
}

func (e *Base62Encoder) Decode(value string) (uint64, error) {
	var result uint64
	for _, c := range value {
		if result > math.MaxUint64/62 {
			return 0, fmt.Errorf("value overflows uint64")
		}
		result *= 62
		switch {
		case c >= '0' && c <= '9':
			result += uint64(c - '0')
		case c >= 'A' && c <= 'Z':
			result += uint64(c-'A') + 10
		case c >= 'a' && c <= 'z':
			result += uint64(c-'a') + 36
		default:
			return 0, fmt.Errorf("invalid character: %c", c)
		}
	}
	return result, nil
}
