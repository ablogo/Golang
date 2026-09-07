package utils

func GetDefaultIfNull[T comparable](value T, default_value T) (result T) {
	var zero_value T
	if value == zero_value {
		result = default_value
	} else {
		result = value
	}
	return
}
