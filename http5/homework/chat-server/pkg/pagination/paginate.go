package pagination

func Paginate[T any](items []T, offset, limit int) []T {
	startIndex := offset
	endIndex := startIndex + limit

	if endIndex > len(items) {
		endIndex = len(items)
	}

	if startIndex > len(items) {
		startIndex = len(items)
	}

	return items[startIndex:endIndex]
}
