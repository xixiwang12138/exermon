package collection

func Map[T, R any](list []T, mapFunc func(T) R) []R {
	ans := make([]R, len(list))
	for i, e := range list {
		ans[i] = mapFunc(e)
	}
	return ans
}

func Filter[T any](list []T, retainFunc func(T) bool) []T {
	ans := make([]T, 0, len(list))
	for _, e := range list {
		if retainFunc(e) {
			ans = append(ans, e)
		}
	}
	return ans
}

func GroupBy[T any, K comparable](list []T, keyFunc func(T) K) map[K][]T {
	ans := make(map[K][]T)
	for _, e := range list {
		k := keyFunc(e)
		ans[k] = append(ans[k], e)
	}
	return ans
}
