# Agent Instructions for Movie Stream App

## Build/Lint/Test Commands

### Frontend (React/TypeScript)
- **Build**: `cd Client/frontend && npm run build`
- **Lint**: `cd Client/frontend && npm run lint`
- **Test all**: `cd Client/frontend && npm run test`
- **Test single**: `cd Client/frontend && npm run test -- AddMovie.test.tsx`
- **Dev server**: `cd Client/frontend && npm run dev`

### Backend (Go)
- **Build**: `cd Server/StreamMoviesServer && go build`
- **Run**: `cd Server/StreamMoviesServer && go run main.go`
- **Test**: `cd Server/StreamMoviesServer && go test ./...`

## Code Style Guidelines

### TypeScript/React
- **Imports**: Group by type (React, third-party, local); absolute paths `@/*`
- **Components**: Functional with hooks; PascalCase naming
- **Types**: Strict TypeScript; PascalCase interfaces, snake_case properties
- **Naming**: camelCase variables/functions, PascalCase types/components
- **Error handling**: Let errors bubble up; minimal try-catch
- **Formatting**: ESLint + TypeScript rules; no semicolons

### Go
- **Imports**: Standard library, third-party, local order
- **Naming**: PascalCase exported, camelCase unexported
- **Error handling**: Early returns; context timeouts for DB ops
- **Structure**: Gin handlers as functions; dependency injection

### General
- **Testing**: Vitest frontend, Go testing backend
- **Commits**: Run lint/tests before committing
- **Security**: Never log/commit secrets; validate inputs