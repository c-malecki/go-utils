package pstring

import "fmt"

type NumberValue interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func SafeNumPtrToStr[T NumberValue](v *T) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", *v)
}

func SafeStrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
