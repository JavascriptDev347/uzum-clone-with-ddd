package application

const (
	DefaultGalleryPage     = 1
	DefaultGalleryPageSize = 20
	MaxGalleryPageSize     = 100
)

// NormalizeGalleryPagination - noto'g'ri yoki bo'sh page/pageSize qiymatlarini standart
// qiymatlarga almashtiradi va pageSize'ni MaxGalleryPageSize bilan cheklaydi.
func NormalizeGalleryPagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultGalleryPage
	}
	if pageSize < 1 {
		pageSize = DefaultGalleryPageSize
	}
	if pageSize > MaxGalleryPageSize {
		pageSize = MaxGalleryPageSize
	}
	return page, pageSize
}
