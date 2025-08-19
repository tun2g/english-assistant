package repositories

import (
	"app-backend/internal/types"
	
	"gorm.io/gorm"
)

// QueryOptions contains options for querying with pagination
type QueryOptions struct {
	Conditions   map[string]interface{} // WHERE conditions
	SearchFields []string               // Fields to search in
	Preloads     []string               // Associations to preload
}

// BaseRepositoryInterface defines common CRUD operations
type BaseRepositoryInterface[T any] interface {
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	Update(entity *T) error
	Delete(id uint) error
	List(req *types.PaginationRequest, opts *QueryOptions) (*types.PaginationResponse[T], error)
	FindBy(field string, value interface{}) (*T, error)
	FindAllBy(field string, value interface{}) ([]*T, error)
}

// BaseRepository provides common database operations
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create inserts a new entity
func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// GetByID retrieves an entity by its ID
func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update saves an entity
func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete soft deletes an entity by ID
func (r *BaseRepository[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// List retrieves entities with pagination and optional conditions/search
func (r *BaseRepository[T]) List(req *types.PaginationRequest, opts *QueryOptions) (*types.PaginationResponse[T], error) {
	var entities []T
	var total int64
	var entity T

	// Set defaults if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Handle nil options
	if opts == nil {
		opts = &QueryOptions{}
	}

	// Build base query
	query := r.db.Model(&entity)

	// Apply WHERE conditions
	if opts.Conditions != nil {
		for field, value := range opts.Conditions {
			query = query.Where(field+" = ?", value)
		}
	}

	// Apply search conditions
	if req.Search != "" && len(opts.SearchFields) > 0 {
		searchQuery := r.db.Model(&entity)
		
		// Apply existing conditions to search query too
		if opts.Conditions != nil {
			for field, value := range opts.Conditions {
				searchQuery = searchQuery.Where(field+" = ?", value)
			}
		}
		
		// Add search conditions
		searchQuery = searchQuery.Where("1=0") // Start with false condition
		for _, field := range opts.SearchFields {
			searchQuery = searchQuery.Or(field+" ILIKE ?", "%"+req.Search+"%")
		}
		query = searchQuery
	}

	// Apply preloads
	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	// Count total records (create a separate query for counting to avoid issues with preloads)
	countQuery := r.db.Model(&entity)
	if opts.Conditions != nil {
		for field, value := range opts.Conditions {
			countQuery = countQuery.Where(field+" = ?", value)
		}
	}
	if req.Search != "" && len(opts.SearchFields) > 0 {
		countQuery = countQuery.Where("1=0")
		for _, field := range opts.SearchFields {
			countQuery = countQuery.Or(field+" ILIKE ?", "%"+req.Search+"%")
		}
	}
	
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Get paginated results with sorting
	err := query.
		Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&entities).Error
		
	if err != nil {
		return nil, err
	}

	return types.NewPaginationResponse(entities, req, total), nil
}

// FindBy finds a single entity by a specific field
func (r *BaseRepository[T]) FindBy(field string, value interface{}) (*T, error) {
	var entity T
	err := r.db.Where(field+" = ?", value).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAllBy finds all entities by a specific field
func (r *BaseRepository[T]) FindAllBy(field string, value interface{}) ([]*T, error) {
	var entities []*T
	err := r.db.Where(field+" = ?", value).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// GetDB returns the database instance for custom queries
func (r *BaseRepository[T]) GetDB() *gorm.DB {
	return r.db
}