package postgres

func ToOffsetLimit(page, size int) (int, int) {
	offset := 0
	limit := 1000000000

	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 1000000000
	}

	offset = size * (page - 1)
	limit = size

	return offset, limit
}
