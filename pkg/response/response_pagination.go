package response

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResult struct {
	Items      any        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func NewPaginatedResult(items any, total int64, page, pageSize int) PaginatedResult {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // ceil, math.Ceil'siz
	if totalPages < 1 {
		totalPages = 1
	}
	return PaginatedResult{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}
