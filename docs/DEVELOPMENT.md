# Development Guide

## Overview

This document provides guidelines and instructions for contributing to the Cutter Project. It covers coding standards, development workflow, testing practices, and deployment procedures.

## Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Docker (optional)
- Git

## Development Environment Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/cutterproject.git
cd cutterproject
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Environment Variables

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```bash
# Server Configuration
GO_SERVER=:8080

# Database Configuration
POSTGRES_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable

# Redis Configuration
REDIS_URL=redis://localhost:6379

# Logging Configuration
LOG_LEVEL=info
```

### 4. Set Up Database

```bash
# Create database
createdb cutterproject

# Run migrations
psql $POSTGRES_URL -f db/migrations/001_create_users_table.sql
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The application will be available at `http://localhost:8080`.

## Coding Standards

### Go Guidelines

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use [gofmt](https://golang.org/cmd/gofmt/) for code formatting
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names
- Keep functions small and focused on a single responsibility
- Use interfaces to define contracts between packages
- Prefer composition over inheritance

### Project Structure

- Follow the existing project structure
- Place new files in appropriate directories
- Use descriptive file names
- Keep related functionality together

### Error Handling

- Use custom error types defined in `internal/exception/errors.go`
- Wrap errors with context using `fmt.Errorf` or custom error wrapping
- Handle errors explicitly and appropriately
- Log errors with sufficient context
- Return appropriate HTTP status codes for API errors

### Logging

- Use structured logging with Zap
- Include relevant context in log messages
- Use appropriate log levels (Debug, Info, Warn, Error)
- Avoid logging sensitive information
- Log errors with stack traces when appropriate

### Testing

- Write unit tests for all new functionality
- Use table-driven tests for test cases with similar structure
- Mock external dependencies for unit tests
- Write integration tests for critical paths
- Aim for high test coverage

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Implement your feature or fix
- Write tests for your changes
- Ensure all tests pass
- Format your code with `gofmt`
- Run linters to check for issues

### 3. Commit Changes

```bash
git add .
git commit -m "feat: add your feature description"
```

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages:

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

### 4. Push Changes

```bash
git push origin feature/your-feature-name
```

### 5. Create a Pull Request

- Go to the GitHub repository
- Click "New Pull Request"
- Select your feature branch
- Provide a clear description of your changes
- Link any relevant issues
- Request a review from a team member

### 6. Address Review Comments

- Make requested changes
- Commit and push updates
- Respond to review comments

### 7. Merge Pull Request

- Ensure all checks pass
- Resolve any conflicts
- Merge the pull request
- Delete the feature branch

## Testing

### Unit Tests

Run unit tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with coverage and generate a coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests

Run integration tests:

```bash
go test -tags=integration ./...
```

### Benchmark Tests

Run benchmark tests:

```bash
go test -bench=. ./...
```

### Test Database

For integration tests, you can set up a test database:

```bash
# Create test database
createdb cutterproject_test

# Set test database URL in environment
export TEST_POSTGRES_URL=postgres://username:password@localhost:5432/cutterproject_test?sslmode=disable

# Run tests with test database
go test -tags=integration ./...
```

## Building and Deployment

### Local Build

```bash
go build -o bin/cutterproject cmd/main.go
```

### Docker Build

```bash
docker build -t cutterproject .
```

### Docker Compose

```bash
docker-compose up --build
```

### Production Deployment

1. Build the application:

```bash
CGO_ENABLED=0 GOOS=linux go build -o bin/cutterproject cmd/main.go
```

2. Build Docker image:

```bash
docker build -t cutterproject:latest .
```

3. Push to registry:

```bash
docker tag cutterproject:latest your-registry/cutterproject:latest
docker push your-registry/cutterproject:latest
```

4. Deploy to your environment using your preferred deployment method.

## Debugging

### Logging

- Use structured logging with Zap
- Set log level to `debug` for detailed logging
- Include request IDs for tracing requests through the system

### Debugging Tools

- Use Delve for debugging Go applications:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/main.go
```

- Use pprof for profiling:

```bash
go install github.com/google/pprof@latest
```

### Common Issues

#### Database Connection Issues

- Check database URL in `.env` file
- Ensure database is running and accessible
- Check database credentials

#### Redis Connection Issues

- Check Redis URL in `.env` file
- Ensure Redis is running and accessible
- Check Redis credentials

#### Port Already in Use

- Check if another application is using the same port
- Change the port in `.env` file
- Kill the process using the port:

```bash
lsof -ti:8080 | xargs kill -9
```

## Performance Optimization

### Database Optimization

- Use connection pooling
- Optimize database queries
- Use indexes for frequently queried fields
- Consider read replicas for read-heavy workloads

### Caching

- Use Redis for caching frequently accessed data
- Implement cache invalidation strategies
- Consider cache warming for critical data

### Application Optimization

- Use efficient data structures
- Minimize memory allocations
- Use streaming for large data transfers
- Implement rate limiting to prevent abuse

## Security Considerations

### Authentication and Authorization

- Use JWT for authentication
- Implement proper token validation
- Use role-based access control
- Validate user permissions for sensitive operations

### Input Validation

- Validate all input data
- Use proper validation for user inputs
- Sanitize user inputs to prevent injection attacks
- Use parameterized queries for database operations

### Security Headers

- Implement security-related HTTP headers
- Use HTTPS in production
- Implement CORS policies
- Use Content Security Policy (CSP)

### Sensitive Data

- Never log sensitive data
- Use secure storage for secrets
- Implement proper encryption for sensitive data
- Use secure password hashing

## Contributing Guidelines

### Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Welcome newcomers and help them learn
- Be patient and understanding

### Pull Request Process

- Ensure your code follows the project's coding standards
- Write tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Request a review from a team member

### Issue Reporting

- Use the GitHub issue tracker
- Provide clear and detailed descriptions
- Include steps to reproduce
- Include relevant error messages and logs
- Suggest potential solutions if possible

## Resources

### Documentation

- [Go Documentation](https://golang.org/doc/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)

### Tools

- [Go Tools](https://golang.org/cmd/)
- [Delve Debugger](https://github.com/go-delve/delve)
- [gofmt](https://golang.org/cmd/gofmt/)
- [golint](https://github.com/golang/lint)
- [gosec](https://github.com/securego/gosec)

### Community

- [Go Forum](https://forum.golangbridge.org/)
- [Go Slack](https://gophers.slack.com/)
- [Go Reddit](https://www.reddit.com/r/golang/)
- [Go Twitter](https://twitter.com/golang)