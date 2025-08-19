# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

English Learning Assistant is a comprehensive monorepo containing:

- **Backend**: Go REST API with Gin framework, PostgreSQL, JWT auth, and external integrations (YouTube API, Google Gemini)
- **Web App**: React/TypeScript SPA with Vite, TanStack Query, React Router, and Tailwind CSS  
- **Chrome Extension**: Manifest V3 extension with TypeScript, Framework7, and shared utilities
- **Admin App**: Administrative interface (template structure)
- **Shared Packages**: Common utilities, types, API clients, and UI components

## Development Commands

### Monorepo Management (pnpm + Turbo)

```bash
# Install dependencies for entire monorepo
pnpm install

# Build all apps and packages
pnpm build

# Start development servers for all apps
pnpm dev

# Run linting across all packages
pnpm lint

# Run type checking across all packages
pnpm type-check

# Run tests across all packages
pnpm test

# Format code across all packages
pnpm format

# Clean build artifacts
pnpm clean
```

### Backend (Go/Gin API)

Navigate to `apps/backend/` first:

```bash
# Development with live reload
make dev

# Build application
make build

# Run tests with coverage
make test-coverage

# Start database and services
make docker-up

# Run database migrations
make migrate-up

# Generate Swagger docs
make swagger

# Lint Go code
make lint

# Format Go code
make format

# Install development tools
make install-tools

# Setup development environment
make setup
```

### Web App (React/Vite)

Navigate to `apps/web/` first:

```bash
# Development server
pnpm dev

# Build for production
pnpm build

# Type checking
pnpm type-check

# Linting
pnpm lint
```

### Chrome Extension

Navigate to `apps/extension/` first:

```bash
# Development build with watch
pnpm dev

# Production build
pnpm build

# Type checking
pnpm type-check
```

## Backend Architecture Patterns

### Clean Architecture Implementation

The backend follows clean architecture with clear separation of concerns:

```
handlers/     # HTTP handlers (Gin controllers)
  ├── auth/   # Authentication endpoints
  ├── user/   # User management endpoints  
  └── video/  # Video processing endpoints

services/     # Business logic layer
  ├── auth/   # Authentication business logic
  ├── jwt/    # JWT token management
  ├── user/   # User operations
  └── video/  # Video processing logic

repositories/ # Data access layer
  ├── base_repository.go     # Generic CRUD operations
  ├── user_repository.go     # User-specific queries
  └── session_repository.go  # Session management
```

### Generic Repository Pattern

Uses Go generics for type-safe CRUD operations:

```go
// BaseRepositoryInterface provides generic CRUD operations
type BaseRepositoryInterface[T any] interface {
    Create(entity *T) error
    GetByID(id uint) (*T, error
    Update(entity *T) error
    Delete(id uint) error
    List(req *types.PaginationRequest, opts *QueryOptions) (*types.PaginationResponse[T], error)
}

// UserRepository extends base with user-specific methods
type UserRepositoryInterface interface {
    BaseRepositoryInterface[models.User]
    GetByEmail(email string) (*models.User, error)
    GetActiveUsers(req *types.PaginationRequest) (*types.PaginationResponse[models.User], error)
}
```

### Model Patterns

All models extend the `Auditable` base struct with GORM hooks:

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

### Error Handling Pattern

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

### Dependency Injection Container

All dependencies are managed through a container:

```go
type Container struct {
    Config     *config.Config
    DB         *gorm.DB
    Logger     *logger.Logger
    
    // Repositories
    UserRepository    repositories.UserRepositoryInterface
    SessionRepository repositories.SessionRepositoryInterface
    
    // Services  
    JWTService   jwtService.ServiceInterface
    AuthService  authService.ServiceInterface
    
    // Handlers
    AuthHandler  auth.HandlerInterface
    UserHandler  user.HandlerInterface
}
```

## Frontend Coding Conventions

### File Naming Rules

**ALWAYS use kebab-case for ALL filenames and directories:**

```
✅ user-profile.tsx, api-client.ts, auth-types.ts
❌ UserProfile.tsx, userProfile.tsx, user_profile.tsx
```

**File naming by type:**

- **Non-component files** use parent folder suffix: `{name}-{parent-folder}.ts`
  - `constants/app-config-constants.ts`
  - `types/user-profile-types.ts`
  - `lib/dayjs-lib.ts`

- **Component files** rely on folder context: `{descriptive-name}.tsx`
  - `certifications-table.tsx`, `login-form.tsx`

### TypeScript Rules

**ALWAYS prefer `interface` over `type` for object shapes:**

```typescript
// ✅ Use interface for objects and props
interface UserProfileProps {
  user: User;
  onUpdate: (data: Partial<User>) => void;
}

// ✅ Use type for specific cases only
type Status = 'loading' | 'success' | 'error';
type ApiResponse<T> = { status: 'success'; data: T } | { status: 'error'; message: string };
```

**NEVER use `enum` - use `as const` instead:**

```typescript
// ✅ Use as const
const USER_ROLES = {
  ADMIN: 'admin',
  USER: 'user',
  SYSTEM_ADMIN: 'system_admin'
} as const;
type UserRole = typeof USER_ROLES[keyof typeof USER_ROLES];
```

### React Patterns

**Component Structure:**

- **Pages**: Top-level route components (`src/pages/`)
- **Containers**: Presentation components with data fetching (`src/containers/`)
- **Components**: Reusable UI components (`src/components/`)

**Import Conventions:**

```typescript
// ✅ Direct React imports
import { useState, useEffect } from 'react';

// ✅ Type-only imports
import type { User } from './types';
import { someFunction, type SomeType } from './module';
```

**Server vs Client Components:**

- **Default to Server Components** - only use `'use client';` when needed
- Use client components for: event handlers, state, lifecycle effects, browser APIs

### Styling Guidelines

**ALWAYS use Tailwind CSS with CSS-first configuration:**

```css
/* src/styles/globals.css */
@import "tailwindcss";

@theme {
  --color-brand-primary: oklch(0.6 0.2 250);
  --font-display: "Inter", sans-serif;
}
```

**MUST use `cn` function for class concatenation:**

```typescript
// ✅ Correct usage
<div className={cn('base-class', isActive && 'active-class')} />

// ❌ Avoid concatenation
<div className={`base-class ${isActive ? 'active-class' : ''}`} />
```

## Development Patterns

### API Integration

**TanStack Query Usage:**

```typescript
// ✅ Use query key constants
import { QUERY_KEY } from '@/lib/constants/query-key-constants';

const { data } = useQuery({
  queryKey: [QUERY_KEY.AUTH.LOGOUT],
  queryFn: logoutUser,
});

// ✅ Mutations with proper invalidation
const mutation = useMutation({
  mutationFn: updateUser,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: [QUERY_KEY.USER.PROFILE] });
  },
});
```

### Error Handling

**Generic error messages on client-side:**

```typescript
// ✅ Generic error messages
onError: () => toast.error('Failed to create job posting. Please try again.')

// ❌ Never expose server error details
onError: (error) => toast.error(error.message)
```

### Type Safety Between Backend and Frontend

**Shared types pattern** - Keep DTOs aligned:

```typescript
// Backend (Go)
type AuthResponse struct {
    User         *UserResponse `json:"user"`
    AccessToken  string        `json:"access_token"`
    RefreshToken string        `json:"refresh_token"`
    TokenType    string        `json:"token_type"`
    ExpiresIn    int           `json:"expires_in"`
}

// Frontend (TypeScript)
interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
}
```

## Architecture

### Backend Architecture

- **Clean Architecture**: Layered structure with handlers → services → repositories
- **Dependency Injection**: Container pattern manages all dependencies
- **Configuration**: Viper-based config with YAML files and environment variables
- **Database**: GORM with PostgreSQL, automated migrations
- **Authentication**: JWT-based with access/refresh tokens and session management
- **Middleware**: Request logging, CORS, authentication, error handling, recovery
- **External APIs**: YouTube Data API and Google Gemini AI integration
- **API Documentation**: Swagger/OpenAPI auto-generation

### Frontend Architecture

- **Component Architecture**: Pages → Containers → Components pattern
- **State Management**: TanStack Query for server state, React Hook Form for forms
- **Routing**: React Router with protected routes
- **Styling**: Tailwind CSS with shared UI components
- **API Layer**: Axios client in shared package with type-safe endpoints

### Shared Package Structure

- **@english/shared**: API clients, constants, types, utilities, storage abstractions
- **@english/ui**: Reusable UI components with Tailwind and Radix UI
- **@english/eslint-config**: Shared ESLint configurations for all apps

### Key Patterns

- **Interface-based Design**: All services and repositories implement interfaces
- **Error Handling**: Centralized error types and structured logging
- **Type Safety**: Full TypeScript coverage with shared types across frontend/backend
- **Configuration Management**: Environment-specific configs with sensible defaults

## Development Setup

1. Install prerequisites: Node.js 18+, pnpm 8+, Go 1.21+, Docker
2. Clone repository and run `pnpm install`
3. For backend development:
   - Run `cd apps/backend && make setup`
   - Start services: `make docker-up`
   - Run migrations: `make migrate-up`
   - Start dev server: `make dev`
4. For frontend development:
   - Start web app: `cd apps/web && pnpm dev`
   - For extension: `cd apps/extension && pnpm dev`

## Database

- **Primary DB**: PostgreSQL with GORM
- **Migrations**: Located in `apps/backend/migrations/`
- **Models**: User, Session, VideoTranscript with base model patterns
- **Connection**: Configurable via config files or environment variables

## Testing

- **Backend**: Go testing with coverage reports (`make test-coverage`)
- **Frontend**: Configuration exists for testing with Turbo pipeline
- **Linting**: ESLint with shared configs, golangci-lint for Go
- **Type Checking**: TypeScript strict mode across all packages

## External Services

- **YouTube Data API**: Video metadata and transcript retrieval
- **Google Gemini**: AI-powered content analysis and processing
- **Redis**: Session storage and caching (optional)

## Deployment

- **DevOps**: Kubernetes configs, Docker Compose for production, Nginx configs
- **Scripts**: Automated deployment scripts in `devops/scripts/`
- **Environments**: Development, staging, production configurations