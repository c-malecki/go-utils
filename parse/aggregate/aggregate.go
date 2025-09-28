package aggregate

func AggregateUint32sFromStruct[T any](items []T, extract func(T) []uint32) []uint32 {
	seen := make(map[uint32]struct{})
	var result []uint32

	for _, item := range items {
		for _, id := range extract(item) {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				result = append(result, id)
			}
		}
	}

	return result
}

func AggregateUint32sFromSlices(slices ...[]uint32) []uint32 {
	seen := make(map[uint32]struct{})
	var result []uint32

	for _, item := range slices {
		for _, id := range item {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				result = append(result, id)
			}
		}
	}

	return result
}

func AggregateStringsFromStruct[T any](items []T, extract func(T) []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, item := range items {
		for _, id := range extract(item) {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				result = append(result, id)
			}
		}
	}

	return result
}
