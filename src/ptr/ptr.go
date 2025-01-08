package ptr

func Wrap[T any](value T) *T {
	return &value
}

func Unwrap[T any](ptr *T, defaultVal T) T {
	if ptr == nil {
		return defaultVal
	}

	return *ptr
}

func UnwrapWithDefault[T any](ptr *T) T {
	var defVal T

	return Unwrap(ptr, defVal)
}
