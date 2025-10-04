# Contributing to chronogo

Thank you for your interest in contributing to chronogo! We welcome contributions from the community and are pleased to have you join us.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to see if the problem has already been reported. When creating a bug report, please include:

- A clear and descriptive title
- Steps to reproduce the behavior
- Expected behavior
- Actual behavior
- Go version and operating system
- Code samples that demonstrate the issue

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- A clear and descriptive title
- A detailed description of the proposed feature
- Use cases that would benefit from this feature
- Any relevant examples from other libraries

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`go test ./...`)
6. Run go fmt (`go fmt ./...`)
7. Run go vet (`go vet ./...`)
8. Commit your changes (`git commit -m 'Add amazing feature'`)
9. Push to the branch (`git push origin feature/amazing-feature`)
10. Open a Pull Request

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/coredds/chronogo.git
   cd chronogo
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

4. Run tests with coverage:
   ```bash
   go test -cover ./...
   ```

5. Run benchmarks:
   ```bash
   go test -bench=. ./...
   ```

## Coding Standards

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `go fmt` to format your code
- Use `go vet` to check for common errors
- Write clear, self-documenting code

### Testing

- Write tests for all new functionality
- Maintain or improve test coverage
- Include both unit tests and integration tests where appropriate
- Test edge cases and error conditions

### Documentation

- Update documentation for any API changes
- Include examples in documentation
- Write clear commit messages
- Update CHANGELOG.md for notable changes

## Commit Message Format

Use the following format for commit messages:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

Examples:
- `feat(datetime): add timezone normalization`
- `fix(parse): handle edge case in RFC3339 parsing`
- `docs(readme): update installation instructions`

## Release Process

chronogo follows [Semantic Versioning](https://semver.org/):

- MAJOR version when you make incompatible API changes
- MINOR version when you add functionality in a backwards compatible manner
- PATCH version when you make backwards compatible bug fixes

## Getting Help

If you need help, you can:

- Open an issue for bugs or feature requests
- Start a discussion for questions or ideas
- Check the documentation and examples

## License

By contributing to chronogo, you agree that your contributions will be licensed under the MIT License.