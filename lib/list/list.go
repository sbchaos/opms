package list

func Remove[T comparable](list []T, item T) []T {
	output := make([]T, 0, len(list)-1)
	for _, v := range list {
		if v != item {
			output = append(output, v)
		}
	}
	return output
}
