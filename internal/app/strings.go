package app

import "strings"

func StringInSlice(ss []string, s string) bool {
	return StringInSliceFn(ss, s, func(a, b string) bool { return a == b })
}

func StringInSliceEqualFold(ss []string, s string) bool {
	return StringInSliceFn(ss, s, strings.EqualFold)
}

func StringInSliceFn(ss []string, s string, compareFn func(string, string) bool) bool {
	for _, v := range ss {
		if compareFn(v, s) {
			return true
		}
	}
	return false
}

func LookupByKeyEqualFold(m map[string]string, s string) (string, bool) {
	return LookupByKeyFn(m, s, strings.EqualFold)
}

func LookupByKeyFn(m map[string]string, s string, compareFn func(string, string) bool) (string, bool) {
	for k, v := range m {
		if compareFn(k, s) {
			return v, true
		}
	}
	return "", false
}
