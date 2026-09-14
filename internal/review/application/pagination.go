package application

const (
	DefaultReviewPage     = 1
	DefaultReviewPageSize = 20
	MaxReviewPageSize     = 100
)

// NormalizeReviewPagination - noto'g'ri yoki bo'sh page/pageSize qiymatlarini standart
// qiymatlarga almashtiradi va pageSize'ni MaxReviewPageSize bilan cheklaydi.
func NormalizeReviewPagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultReviewPage
	}
	if pageSize < 1 {
		pageSize = DefaultReviewPageSize
	}
	if pageSize > MaxReviewPageSize {
		pageSize = MaxReviewPageSize
	}
	return page, pageSize
}
