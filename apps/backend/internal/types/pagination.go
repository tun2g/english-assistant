package types

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDir  string `json:"sort_dir" form:"sort_dir" binding:"oneof=asc desc"`
	Search   string `json:"search" form:"search"`
}

// PaginationResponse represents paginated response
type PaginationResponse[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// PaginationMetadata contains pagination metadata
type PaginationMetadata struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

// NewPaginationRequest creates a new pagination request with defaults
func NewPaginationRequest() *PaginationRequest {
	return &PaginationRequest{
		Page:     1,
		PageSize: 10,
		SortBy:   "id",
		SortDir:  "desc",
	}
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the limit for database queries
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

// GetOrderBy returns the order by clause for database queries
func (p *PaginationRequest) GetOrderBy() string {
	if p.SortBy == "" {
		p.SortBy = "id"
	}
	if p.SortDir == "" {
		p.SortDir = "desc"
	}
	return p.SortBy + " " + p.SortDir
}

// NewPaginationResponse creates a new paginated response
func NewPaginationResponse[T any](data []T, req *PaginationRequest, totalRecords int64) *PaginationResponse[T] {
	totalPages := int((totalRecords + int64(req.PageSize) - 1) / int64(req.PageSize))
	
	return &PaginationResponse[T]{
		Data: data,
		Pagination: PaginationMetadata{
			CurrentPage:  req.Page,
			PageSize:     req.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
			HasNext:      req.Page < totalPages,
			HasPrev:      req.Page > 1,
		},
	}
}