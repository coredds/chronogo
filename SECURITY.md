# Security Policy

## Supported Versions

We actively support the following versions of ChronoGo with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.5.x   | :white_check_mark: |
| 0.4.x   | :white_check_mark: |
| < 0.4   | :x:                |

## Security Features

ChronoGo implements comprehensive security measures:

### Automated Security Scanning

- **GitHub Actions Security Workflow**: Automated security scanning on every push and pull request
- **Weekly Security Scans**: Scheduled vulnerability scans every Monday
- **CodeQL Analysis**: Static application security testing for Go code
- **Dependency Review**: Automated dependency vulnerability scanning for pull requests

### Security Tools Integration

1. **govulncheck**: Go vulnerability database scanning
2. **CodeQL**: Advanced static analysis for security patterns
3. **Dependency Review Action**: GitHub's dependency vulnerability scanner
4. **go vet**: Static analysis for common Go programming errors

### Local Security Testing

Run security checks locally using:

```bash
# Unix/Linux/macOS
make security
./scripts/security-check.sh

# Windows PowerShell
.\scripts\security-check.ps1
```

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in ChronoGo, please follow these steps:

### 1. Do NOT create a public GitHub issue

For security vulnerabilities, please do not create public GitHub issues as this could expose the vulnerability to potential attackers.

### 2. Contact us privately

Please report security vulnerabilities by emailing us directly at:
- **Email**: [security@coredds.com](mailto:security@coredds.com)
- **Subject**: [SECURITY] ChronoGo Security Vulnerability Report

### 3. Include detailed information

When reporting a vulnerability, please include:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact assessment
- Any proof-of-concept code (if applicable)
- Your contact information for follow-up

### 4. Response timeline

We are committed to responding to security reports promptly:

- **Initial Response**: Within 48 hours of receiving your report
- **Vulnerability Assessment**: Within 7 days
- **Fix Development**: Timeline depends on severity (critical issues prioritized)
- **Disclosure**: Coordinated disclosure after fix is available

## Security Best Practices for Users

When using ChronoGo in your applications:

### 1. Keep Dependencies Updated

```bash
# Check for dependency updates
go get -u ./...
go mod tidy
```

### 2. Use Go Modules

Always use Go modules for dependency management:
```bash
go mod init your-project
go mod download
```

### 3. Vulnerability Scanning

Regularly scan your projects:
```bash
# Install govulncheck if not already installed
go install golang.org/x/vuln/cmd/govulncheck@latest

# Scan your project
govulncheck ./...
```

### 4. Static Analysis

Use static analysis tools:
```bash
# Built-in Go static analysis
go vet ./...

# Optional: Use golangci-lint for comprehensive linting
golangci-lint run
```

## Security Development Lifecycle

### Code Review Process

- All code changes require review by at least one maintainer
- Security-sensitive changes require additional security review
- Automated security scanning runs on all pull requests

### Dependency Management

- Dependencies are regularly updated and scanned for vulnerabilities
- New dependencies are evaluated for security implications
- Dependabot automatically creates pull requests for security updates

### Release Security

- All releases undergo security scanning before publication
- Security fixes are prioritized and released as patch versions
- Release notes clearly document any security-related changes

## Security Contacts

- **Primary Security Contact**: security@coredds.com
- **GitHub Security Advisories**: [ChronoGo Security Advisories](https://github.com/coredds/ChronoGo/security/advisories)

## Acknowledgments

We appreciate the security research community's efforts in responsibly disclosing vulnerabilities. Contributors who report valid security issues will be acknowledged in our security advisories (with their permission).

---

**Last Updated**: August 2025  
**Version**: 1.0
