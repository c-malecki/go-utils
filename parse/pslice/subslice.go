package pslice

func SubSlice[T any](items []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var subs [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		subs = append(subs, items[i:end])
	}
	return subs
}
