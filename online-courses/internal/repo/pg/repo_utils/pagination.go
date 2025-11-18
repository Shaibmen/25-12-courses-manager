package repoutils

func Pagination(page int) (int, int) {

	const LIMIT_COUNT = 25
	offset := (page - 1) * LIMIT_COUNT

	return LIMIT_COUNT, offset
}
