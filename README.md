# English Learning Assistant

A comprehensive monorepo containing a full-stack English learning platform with backend API, web
application, Chrome extension, and shared packages.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Development](#development)
- [Applications](#applications)
- [Deployment](#deployment)
- [Contributing](#contributing)

## 🚀 Overview

The English Learning Assistant is a modern, full-stack platform designed to help users learn English
through interactive video content, AI-powered features, and comprehensive learning tools.

### Key Features

- **AI-Powered Learning**: Integration with Google Gemini for intelligent content analysis
- **Video Integration**: YouTube API integration for transcript extraction and video processing
- **Cross-Platform**: Web application and Chrome extension for seamless learning experiences
- **Authentication**: JWT-based authentication with session management
- **Responsive Design**: Modern UI built with React, Tailwind CSS, and Framework7

### Tech Stack

- **Backend**: Go (Gin framework), PostgreSQL, Redis, GORM
- **Frontend**: React 18, TypeScript, Vite, TanStack Query
- **Extension**: Chrome Extension Manifest V3, React, Framework7
- **Infrastructure**: Docker, Kubernetes, Nginx
- **Monorepo**: pnpm workspaces, Turbo
- **Code Quality**: ESLint, Prettier, Husky, lint-staged

## 🏗️ Architecture

```text
english-learning-assistant/
├── apps/
│   ├── backend/           # Go REST API server
│   ├── web/              # React web application
│   ├── extension/        # Chrome extension
│   └── admin/           # Admin interface
├── packages/
│   ├── shared/          # Common utilities and types
│   ├── ui/              # Shared UI components
│   └── eslint-config/   # Shared ESLint configuration
├── devops/              # Deployment and infrastructure
└── docs/               # Documentation
```

### Backend Architecture

- **Clean Architecture**: Layered structure with handlers → services → repositories
- **Dependency Injection**: Container pattern for dependency management
- **Generic Repository**: Type-safe CRUD operations using Go generics
- **API Documentation**: Auto-generated Swagger/OpenAPI specs

### Frontend Architecture

- **Component-Based**: Pages → Containers → Components hierarchy
- **State Management**: TanStack Query for server state, React Hook Form for forms
- **Type Safety**: Full TypeScript coverage with shared types
- **Styling**: Tailwind CSS with shared design system

## ⚡ Quick Start

### Prerequisites

- **Node.js** >= 18.0.0
- **pnpm** >= 8.0.0
- **Go** >= 1.21
- **Docker** and Docker Compose
- **PostgreSQL** (for local development)

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd english-learning-assistant
   ```

2. **Install dependencies**

   ```bash
   pnpm install
   ```

3. **Set up the backend**

   ```bash
   cd apps/backend
   make setup
   make docker-up     # Start PostgreSQL and Redis
   make migrate-up    # Run database migrations
   ```

4. **Start development servers**

   ```bash
   # From project root - starts all applications
   pnpm dev

   # Or start individual applications
   cd apps/backend && make dev    # Backend API
   cd apps/web && pnpm dev       # Web application
   cd apps/extension && pnpm dev # Chrome extension
   ```

5. **Access applications**
   - Backend API: <http://localhost:8080>
   - Web App: <http://localhost:3000>
   - API Documentation: <http://localhost:8080/docs>

## 🛠️ Development

### Monorepo Commands

All commands should be run from the project root:

```bash
# Install dependencies
pnpm install

# Build all applications
pnpm build

# Start all development servers
pnpm dev

# Run linting across all packages
pnpm lint

# Run type checking across all packages
pnpm type-check

# Format code across all packages
pnpm format

# Run tests across all packages
pnpm test

# Clean build artifacts
pnpm clean
```

### Backend Commands

Navigate to `apps/backend/` first:

```bash
# Development with live reload
make dev

# Build application
make build

# Run tests with coverage
make test-coverage

# Database operations
make migrate-up      # Apply migrations
make migrate-down    # Rollback migrations

# Generate Swagger documentation
make swagger

# Code quality
make lint           # Lint Go code
make format         # Format Go code
```

### Frontend Commands

For web app (`apps/web/`) or extension (`apps/extension/`):

```bash
# Development server
pnpm dev

# Production build
pnpm build

# Type checking
pnpm type-check

# Linting
pnpm lint

# Preview build
pnpm preview
```

### Git Hooks

The repository uses Husky for automated code quality checks:

- **Pre-commit**: ESLint, Prettier, Go formatting and linting
- **Commit message**: Conventional commits format validation

Commit message format: `type(scope): description`

Examples:

- `feat: add user authentication`
- `fix: resolve video loading issue`
- `docs: update API documentation`

## 📱 Applications

### Backend API

**Location**: `apps/backend/`

Go-based REST API with:

- JWT authentication with refresh tokens
- PostgreSQL database with GORM
- Redis for session storage
- YouTube Data API integration
- Google Gemini AI integration
- Comprehensive middleware stack
- Auto-generated API documentation

**Key Endpoints**:

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/user/profile` - User profile
- `POST /api/v1/videos/analyze` - Video analysis

### Web Application

**Location**: `apps/web/`

React-based single-page application featuring:

- Modern React 18 with hooks and Suspense
- TanStack Query for server state management
- React Router for navigation
- React Hook Form with Zod validation
- Tailwind CSS for styling
- Responsive design for all screen sizes

**Key Features**:

- User dashboard and profile management
- Video learning interface
- Progress tracking
- AI-powered content recommendations

### Chrome Extension

**Location**: `apps/extension/`

Manifest V3 Chrome extension with:

- YouTube integration for transcript extraction
- OAuth 2.0 authentication
- Framework7 React UI components
- Background service worker
- Content script injection
- Cross-extension messaging

**Installation**:

1. Build: `cd apps/extension && pnpm build:dev`
2. Open Chrome → `chrome://extensions/`
3. Enable "Developer mode"
4. Click "Load unpacked" → select `dist/` folder

### Shared Packages

**@english/shared** (`packages/shared/`):

- API client utilities
- Common types and interfaces
- Storage abstractions
- Constants and configuration

**@english/ui** (`packages/ui/`):

- Reusable React components
- Tailwind CSS design system
- Form components and utilities

**@english/eslint-config** (`packages/eslint-config/`):

- Shared ESLint configurations
- TypeScript and React rules
- Code style enforcement

## 🚀 Deployment

### Production Build

```bash
# Build all applications for production
pnpm build

# Backend production build
cd apps/backend && make build

# Create optimized Docker images
docker build -t english-backend ./apps/backend
```

### Docker Deployment

Production deployment uses Docker Compose:

```bash
# Start production services
cd devops/docker
docker-compose -f docker-compose.prod.yml up -d
```

**Services**:

- **App**: Go backend server
- **PostgreSQL**: Primary database
- **Redis**: Session storage and caching
- **Nginx**: Reverse proxy and load balancer

### Environment Configuration

Required environment variables for production:

```bash
# Database
DB_HOST=postgres
DB_NAME=english_app
DB_USER=postgres
DB_PASSWORD=your_secure_password

# JWT
JWT_SECRET=your_jwt_secret
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_HOURS=168

# External APIs
YOUTUBE_API_KEY=your_youtube_api_key
GEMINI_API_KEY=your_gemini_api_key
```

### Infrastructure

The application is designed for deployment on:

- **Kubernetes**: Production orchestration
- **Docker**: Containerized deployment
- **Cloud Providers**: AWS, GCP, Azure compatible
- **CDN**: Static asset distribution

## 🤝 Contributing

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt` and `golangci-lint`
- **TypeScript/React**: ESLint + Prettier with shared configuration
- **File Naming**: Always use kebab-case for files and directories
- **Commits**: Use conventional commit format

### Development Workflow

1. Create feature branch: `git checkout -b feat/your-feature`
2. Make changes following established patterns
3. Run tests: `pnpm test`
4. Run linting: `pnpm lint`
5. Commit changes: `git commit -m "feat: add your feature"`
6. Push branch and create pull request

### Code Quality

All contributions must pass:

- TypeScript type checking
- ESLint linting rules
- Go linting with golangci-lint
- Unit and integration tests
- Pre-commit hooks validation

### Architecture Guidelines

- Follow clean architecture principles
- Use dependency injection patterns
- Implement proper error handling
- Write comprehensive tests
- Document public APIs
- Follow security best practices

## 📄 License

This project is proprietary software. All rights reserved.

## 🔗 Links

- [API Documentation](http://localhost:8080/docs)
- [Project Board](link-to-project-management)
- [Design System](link-to-design-docs)

---

For detailed development instructions, see the `CLAUDE.md` file in the project root and individual
application directories.
