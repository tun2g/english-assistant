# CLAUDE.md - Backend

<<<<<<< Updated upstream This file provides guidance to Claude Code (claude.ai/code) when working
with the Go backend in this repository. ======= This file provides guidance to Claude Code
(claude.ai/code) when working with the Go backend in this repository.

> > > > > > > Stashed changes

## Backend Architecture

### Clean Architecture Implementation

The backend follows clean architecture with clear separation of concerns:

```text
handlers/     # HTTP handlers (Gin controllers)
 auth/   # Authentication endpoints
 oauth/  # OAuth flow management
 translation/ # Translation API endpoints
 user/   # User management endpoints
   video/  # Video processing endpoints

services/     # Business logic layer
 auth/   # Authentication business logic
 jwt/    # JWT token management
 oauth/  # OAuth service integration
 transcript/ # Multi-provider transcript extraction
 translation/ # Google Gemini translation
 user/   # User operations
   video/  # Video processing logic

repositories/ # Data access layer
 base_repository.go     # Generic CRUD operations
 user_repository.go     # User-specific queries
   session_repository.go  # Session management
```

### Transcript Service Architecture

The transcript service uses a provider pattern with multiple extraction strategies:

```text
transcript/ providers/
   innertube/        # YouTube internal API provider
<<<<<<< Updated upstream
   kkdai_youtube/    # Community library provider
=======
   kkdai_youtube/    # Community library provider
>>>>>>> Stashed changes
   youtube_api/      # Official YouTube API provider
    yt_transcript/    # Transcript-specific provider service.go            # Main orchestrator with fallback logic
 types/               # Transcript-specific types
```

### External Integrations

**YouTube Data API** (`pkg/youtube/service.go`): <<<<<<< Updated upstream

=======

> > > > > > > Stashed changes

- Video metadata retrieval
- Channel information
- Playlist management

**Google Gemini AI** (`pkg/gemini/service.go`):

- Content analysis and translation
- Language detection
- Summary generation

## Development Commands

Navigate to `apps/backend/` first:

```bash
# Setup development environment
make setup                 # Install tools + copy config

# Database operations
make docker-up            # Start PostgreSQL container
make migrate-up           # Run database migrations
make migrate-create NAME=migration_name  # Create new migration

# Development
make dev                  # Start with hot reload (air)
make run                  # Start without hot reload

# Testing and quality
make test                 # Run all tests
make test-coverage        # Run tests with HTML coverage report
make lint                 # Run golangci-lint
make format               # Format code with go fmt + goimports

# Build and deploy
make build                # Build binary to bin/
make build-linux          # Cross-compile for Linux
make deploy-prod          # Deploy using DevOps scripts

# Documentation
make swagger              # Generate Swagger docs

# Utilities
make clean                # Clean build artifacts
make db-reset             # Reset database (drop + recreate)
make security             # Run gosec security scan
```

## Generic Repository Pattern

Uses Go generics for type-safe CRUD operations:

```go
// BaseRepositoryInterface provides generic CRUD operations
type BaseRepositoryInterface[T any] interface {
    Create(entity *T) error
    GetByID(id uint) (*T, error)
    Update(entity *T) error
    Delete(id uint) error
    List(req *types.PaginationRequest, opts *QueryOptions) (*types.PaginationResponse[T], error)
    FindBy(field string, value interface{}) (*T, error)
    FindAllBy(field string, value interface{}) ([]*T, error)
}

// UserRepository extends base with user-specific methods
type UserRepositoryInterface interface {
    BaseRepositoryInterface[models.User]
    GetByEmail(email string) (*models.User, error)
    GetActiveUsers(req *types.PaginationRequest) (*types.PaginationResponse[models.User], error)
}
```

## Model Patterns

All models extend the `Auditable` base struct:

```go
type Auditable struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type User struct {
    Auditable
    FirstName string `json:"first_name" gorm:"not null"`
    LastName  string `json:"last_name" gorm:"not null"`
    Email     string `json:"email" gorm:"uniqueIndex;not null"`
    Password  string `json:"-" gorm:"not null"`
    IsActive  bool   `json:"is_active" gorm:"default:true"`
    Role      string `json:"role" gorm:"default:'user'"`
}
```

## Error Handling Pattern

Custom error type with HTTP status codes:

```go
type AppError struct {
    Message string
    Err     error
    Status  int
}

// Usage in services
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
    return nil, errors.NewAppError("Invalid credentials", nil, http.StatusUnauthorized)
}
```

## Dependency Injection Container

All dependencies managed through container (`internal/container/container.go`):

```go
type Container struct {
    Config     *config.Config
    DB         *gorm.DB
    Logger     *logger.Logger

    // Repositories
    UserRepository    repositories.UserRepositoryInterface
    SessionRepository repositories.SessionRepositoryInterface

    // Services
    JWTService        jwtService.ServiceInterface
    AuthService       authService.ServiceInterface
    OAuthService      oauthService.ServiceInterface
    TranscriptService transcriptService.ServiceInterface
    TranslationService translationService.ServiceInterface

    // Handlers
    AuthHandler        auth.HandlerInterface
    OAuthHandler       oauth.HandlerInterface
    TranslationHandler translation.HandlerInterface
    UserHandler        user.HandlerInterface
    VideoHandler       video.HandlerInterface
}
```

## Configuration Management

Uses Viper with YAML configuration files:

```yaml
# configs/app.yaml
server:
  port: 8080
  host: localhost
<<<<<<< Updated upstream

=======

>>>>>>> Stashed changes
database:
  host: localhost
  port: 5434
  name: app_backend_dev
  user: postgres
  password: postgres
<<<<<<< Updated upstream

external:
  youtube_api_key: 'your-youtube-api-key'
  gemini_api_key: 'your-gemini-api-key'
=======

external:
  youtube_api_key: "your-youtube-api-key"
  gemini_api_key: "your-gemini-api-key"
>>>>>>> Stashed changes
```

## Database Migrations

Located in `migrations/` directory:

```bash
# Create new migration
make migrate-create NAME=add_video_transcript_table

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Force specific version
make migrate-force VERSION=2
```

## Testing Patterns

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
<<<<<<< Updated upstream

    // Initialize dependencies
    container := setupTestContainer(t, db)

=======

    // Initialize dependencies
    container := setupTestContainer(t, db)

>>>>>>> Stashed changes
    // Test implementation
    user := &models.User{
        Email:     "test@example.com",
        FirstName: "Test",
        LastName:  "User",
    }
<<<<<<< Updated upstream

=======

>>>>>>> Stashed changes
    err := container.UserService.CreateUser(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

## Key Development Tools

Required tools (installed via `make install-tools`):

- **air**: Hot reload for development
- **migrate**: Database migration tool
- **swag**: Swagger documentation generator
- **golangci-lint**: Comprehensive linter
- **goimports**: Import formatting

## External API Integration

**YouTube Data API Integration**: <<<<<<< Updated upstream

=======

> > > > > > > Stashed changes

```go
// pkg/youtube/service.go
type Service struct {
    apiKey string
    client *http.Client
}

func (s *Service) GetVideoDetails(videoID string) (*VideoDetails, error) {
    // Implementation with proper error handling and rate limiting
}
```

**Google Gemini Integration**: <<<<<<< Updated upstream

=======

> > > > > > > Stashed changes

```go
// pkg/gemini/service.go
type Service struct {
    client *genai.Client
    model  *genai.GenerativeModel
}

func (s *Service) TranslateContent(content string, targetLang string) (*TranslationResponse, error) {
    // Implementation with context and error handling
}
```
