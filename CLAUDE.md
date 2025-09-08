# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this
repository.

## Project Overview

English Learning Assistant is a comprehensive monorepo containing:

- **Backend**: Go REST API with Gin framework, PostgreSQL, JWT auth, and external integrations
  (YouTube API, Google Gemini)
- **Web App**: React 18 SPA with Vite, TanStack Query, React Router, React Hook Form, and Tailwind
  CSS
- **Chrome Extension**: Manifest V3 extension with React 18, Framework7 React, and shared utilities
- **Admin App**: Administrative interface (template structure)
- **Shared Packages**: Common utilities, types, API clients, and UI components

## Development Commands

**Package Manager**: Always use `pnpm` - check `package.json` before running any scripts.

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

### Git Hooks and Code Quality

The repository uses Husky for git hooks with automated code quality checks:

```bash
# Pre-commit hooks automatically run:
# - ESLint with --fix for JS/TS files
# - Prettier formatting for all supported files
# - Go formatting and linting for backend files

# Commit message validation:
# - Uses conventional commits format
# - Must follow: type(scope): description
# - Examples: feat: add user authentication, fix: resolve login bug
```

**Supported commit types**: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert

**Lint-staged Configuration:**

- JavaScript/TypeScript files: ESLint with --fix + Prettier formatting
- JSON/Markdown/YAML files: Prettier formatting only
- Go files: Format and lint using backend Makefile commands

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

### Web App (React + Vite)

Navigate to `apps/web/` first:

```bash
# Development server with hot reload
pnpm dev

# Build for production
pnpm build

# Type checking
pnpm type-check

# Linting with auto-fix
pnpm lint

# Preview production build
pnpm preview
```

### Chrome Extension (React + Framework7)

Navigate to `apps/extension/` first:

```bash
# Development build with watch mode
pnpm dev

# Production build
pnpm build

# Development build only
pnpm build:dev

# Type checking
pnpm type-check

# Clean build artifacts
pnpm clean
```

**Extension Development & Testing:**

```bash
# Load extension for testing:
# 1. Build extension: pnpm build:dev
# 2. Open Chrome -> chrome://extensions/
# 3. Enable "Developer mode"
# 4. Click "Load unpacked" -> select dist/ folder

# Debug extension:
# - Background script: chrome://extensions/ -> "Service Worker" link
# - Content script: Browser DevTools -> Console tab
# - Popup: Right-click extension icon -> "Inspect popup"
```

## Backend Architecture Patterns

### Clean Architecture Implementation

The backend follows clean architecture with clear separation of concerns:

```text
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

```text
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

**Directory Organization:**

- **Shared components**: Place in `src/components/`
- **Child/split components**: Place in parent folder under `components/`
- **Type definitions**: Prefer LOCAL scope over global
  - Feature-specific: `{feature}/interfaces/` or `{feature}/constants/`
  - Component-specific: Define directly in component files
  - Global types: `src/lib/types/` ONLY for app-wide shared types

### TypeScript Rules

**ALWAYS prefer `interface` over `type` for object shapes:**

```typescript
// ✅ Use interface for objects and props
interface UserProfileProps {
  user: User;
  onUpdate: (data: Partial<User>) => void;
}

// ✅ Use type ONLY for: unions, primitives, Zod inference, mapped types
type Status = 'loading' | 'success' | 'error';
type ApiResponse<T> = { status: 'success'; data: T } | { status: 'error'; message: string };
type CreateUserRequest = z.infer<typeof createUserSchema>;
```

**NEVER use `enum` - use `as const` instead:**

```typescript
// ✅ Use as const with SCREAMING_SNAKE_CASE constants
const USER_ROLES = {
  ADMIN: 'admin',
  USER: 'user',
  SYSTEM_ADMIN: 'system_admin',
} as const;
type UserRole = (typeof USER_ROLES)[keyof typeof USER_ROLES];
```

**Code Style Requirements:**

```typescript
// ✅ ALWAYS use curly braces for block statements (if, else, while, for)
if (condition) {
  doSomething();
}
for (let i = 0; i < arr.length; i++) {
  process(arr[i]);
}

// ✅ Arrow functions can skip braces for simple expressions
const fn = () => value;
array.map(item => item.name);
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

**Component Export Patterns:**

- **ALWAYS use named exports** (avoid default exports)
- Use `function` syntax for main component exports, complex handlers, reused functions
- Use arrow functions for inline handlers, array methods, callbacks, small utilities
- Custom hooks must start with `use` and return objects (not arrays)
- Event handlers start with `handle`, callback props start with `on`

### Styling Guidelines

**ALWAYS use Tailwind CSS with CSS-first configuration:**

```css
/* src/styles/globals.css */
@import 'tailwindcss';

@theme {
  --color-brand-primary: oklch(0.6 0.2 250);
  --font-display: 'Inter', sans-serif;
}
```

**MUST use `cn` function from `tailwind-utils.ts` for class concatenation:**

```typescript
// ✅ Correct usage
<div className={cn('base-class', isActive && 'active-class')} />

// ❌ Avoid string concatenation
<div className={`base-class ${isActive ? 'active-class' : ''}`} />
```

**UI Styling Guidelines:**

- Use flat and clean UI design patterns
- Prefer responsive breakpoint prefixes (`md:`, `lg:`, etc.)
- Reference existing components before creating new patterns
- Never modify `src/components/ui/` (shadcn/ui) components without approval

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
// ✅ Generic error messages for security
onError: () => toast.error('Failed to create job posting. Please try again.');

// ❌ Never expose server error details to client
onError: error => toast.error(error.message); // Could leak sensitive info
```

**Error Handling Rules:**

- Log detailed errors server-side for debugging
- Show generic messages client-side for security
- Always handle loading and error states in Client Components

### Form Handling with React Hook Form

**MANDATORY form structure using shadcn/ui components:**

```typescript
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

// 1. Zod schema first
const formSchema = z.object({
  fieldName: z.string().min(1, "This field is required."),
});
type FormValues = z.infer<typeof formSchema>;

// 2. Component with strict structure
function MyForm() {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { fieldName: "" }, // Initialize ALL fields
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="fieldName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Field Label</FormLabel>
              <FormControl>
                <Input {...field} placeholder="Enter value" />
              </FormControl>
              <FormMessage /> {/* Auto-displays validation errors */}
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}
```

**Form Rules:**

- Never modify files in `src/components/ui/` (shadcn/ui components)
- Use `<FormMessage />` for automatic error display (no custom error handling)
- Spread `{...field}` when input component properties match field properties
- Follow strict JSX hierarchy: `Form` → `form` → `FormField` → `FormItem` → `FormControl`

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

**Web App (React + Vite):**

- **Component Architecture**: Pages → Containers → Components pattern
- **State Management**: TanStack Query for server state, React Hook Form for forms
- **Routing**: React Router with protected routes
- **Styling**: Tailwind CSS with shared UI components
- **API Layer**: Axios client in shared package with type-safe endpoints

**Chrome Extension (React + Framework7):**

- **Architecture**: Background scripts, content scripts, popup/options React apps
- **UI Framework**: Framework7 React for mobile-like UI components
- **Structure**: Services, features, UI components with shared utilities
- **Integration**: YouTube integration, OAuth management, notification system

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
