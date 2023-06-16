package internal

func Map[T any, R any](a []T, f func(t T) R) []R {
	r := make([]R, 0, len(a))
	for idx := range a {
		r = append(r, f(a[idx]))
	}
	return r
}

func MapKeysAsSlice[K comparable, V any](m map[K]V) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
