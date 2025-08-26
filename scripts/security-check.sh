#!/bin/bash

# Security check script for ChronoGo
# This script runs basic security checks locally

set -e

echo "🔒 Running ChronoGo Security Checks..."
echo "======================================"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run this script from the project root."
    exit 1
fi

echo "📦 Installing security tools..."

# Install govulncheck
if ! command -v govulncheck &> /dev/null; then
    echo "Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

echo ""
echo "🔍 1. Running Go vulnerability check..."
echo "--------------------------------------"
govulncheck ./...

echo ""
echo "� 2. Running go vet..."
echo "---------------------"
go vet ./...

echo ""
echo "🧪 3. Running tests..."
echo "--------------------"
go test ./...

echo ""
echo "� 4. Checking for potential hardcoded secrets..."
echo "------------------------------------------------"
if grep -r --include="*.go" -E "(password|secret|key|token|api)" . | grep -v "test" | grep -v "example" | grep -v "// "; then
    echo "⚠️  Found potential hardcoded secrets - please review above results"
else
    echo "✅ No obvious hardcoded secrets found"
fi

echo ""
echo "� 5. Checking dependencies..."
echo "-----------------------------"
go mod verify
echo "✅ Dependencies verified"

echo ""
echo "✅ Security checks completed!"
echo ""
echo "📝 Recommendations:"
echo "- Review any warnings or vulnerabilities found above"
echo "- Keep dependencies updated with 'go get -u ./...'"
echo "- Run 'go mod tidy' to clean up unused dependencies"
echo "- Consider using 'go mod vendor' for reproducible builds"
echo ""
echo "🔗 For more comprehensive security scanning, push to GitHub"
echo "   to trigger the full security workflow with CodeQL analysis."
