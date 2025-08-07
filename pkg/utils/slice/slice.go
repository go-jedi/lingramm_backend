package slice

// Clip returns a slice with its capacity trimmed to its length.
// This is useful when you want to reduce memory usage of a slice
// that has excess capacity.
func Clip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}

// Identical reports whether two slices are exactly equal in length and element order.
// Elements are compared in increasing index order. Comparison stops at the first mismatch.
// NaN values in float types are treated as not equal.
func Identical[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// Index returns the index of the first occurrence of v in slice s.
// If v is not found, it returns -1.
func Index[E comparable](s []E, v E) int {
	for i, item := range s {
		if item == v {
			return i
		}
	}
	return -1
}

// Contains reports whether slice s contains the value v.
func Contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}

// Equal reports whether two slices contain the same elements,
// regardless of order. Duplicates are considered, i.e., the count
// of each element must be equal in both slices.
//
// NaN values in float types are treated as not equal.
func Equal[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}

	counter := make(map[E]int, len(s1))
	for _, v := range s1 {
		counter[v]++
	}
	for _, v := range s2 {
		if count, ok := counter[v]; !ok || count == 0 {
			return false
		}
		counter[v]--
		if counter[v] == 0 {
			delete(counter, v)
		}
	}

	return len(counter) == 0
}

// ToSliceOfAny converts a slice of any specific type T into a slice of empty interfaces ([]any).
// Useful when a generic slice must be passed to an interface that accepts []any.
func ToSliceOfAny[T any](s []T) []any {
	res := make([]any, 0, len(s))
	for _, v := range s {
		res = append(res, v)
	}
	return res
}
