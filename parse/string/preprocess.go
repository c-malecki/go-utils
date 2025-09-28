package string

import (
	"slices"
	"strings"
)

type PreprocessStringFn func(s string) string

func PreStrSlice(arr []string, pre ...PreprocessStringFn) []string {
	if len(arr) == 0 {
		return []string{}
	}
	for i, s := range arr {
		for _, fn := range pre {
			arr[i] = fn(s)
		}
	}
	return arr
}

func StrToSlice(s string, pre ...PreprocessStringFn) []string {
	if len(s) == 0 {
		return []string{}
	}
	for _, fn := range pre {
		s = fn(s)
	}
	return []string{s}
}

func StrToPtr(s string, pre ...PreprocessStringFn) *string {
	if len(s) == 0 {
		return nil
	}
	for _, fn := range pre {
		s = fn(s)
	}
	return &s
}

func PreStrPtr(s *string, pre ...PreprocessStringFn) *string {
	if s == nil {
		return s
	}
	str := *s
	for _, fn := range pre {
		str = fn(str)
	}
	return &str
}

func ConcatToStrPtr(arr []string, pre ...PreprocessStringFn) *string {
	if len(arr) == 0 {
		return nil
	}

	var strs []string
	for _, v := range arr {
		for _, fn := range pre {
			v = fn(v)
		}
		strs = append(strs, v)
	}
	s := strings.Join(strs, ";")
	return &s
}

func AppendIfNotExists(base []string, merge []string, pre ...PreprocessStringFn) []string {
	for _, s := range merge {
		for _, fn := range pre {
			s = fn(s)
		}
		if !slices.Contains(base, s) {
			base = append(base, s)
		}
	}
	return base
}

func AppendToConcatedStrPtr(base *string, s string, pre ...PreprocessStringFn) *string {
	split := strings.Split(*base, ";")
	for _, fn := range pre {
		s = fn(s)
	}
	if !slices.Contains(split, s) {
		split = append(split, s)
	}
	str := strings.Join(split, ";")
	return &str
}

func ConcatedStrPtrToStrSlice(s *string, pre ...PreprocessStringFn) []string {
	if s == nil {
		return []string{}
	}
	split := strings.Split(*s, ";")
	strs := make([]string, 0, len(split))
	for _, v := range split {
		for _, fn := range pre {
			v = fn(v)
		}
		strs = append(strs, v)
	}
	return strs
}
