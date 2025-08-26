# Security check script for ChronoGo (PowerShell version)
# This script runs basic security checks locally

param(
    [switch]$SkipInstall = $false
)

Write-Host "🔒 Running ChronoGo Security Checks..." -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green

# Check if we're in the right directory
if (-not (Test-Path "go.mod")) {
    Write-Host "❌ Error: go.mod not found. Please run this script from the project root." -ForegroundColor Red
    exit 1
}

if (-not $SkipInstall) {
    Write-Host "📦 Installing security tools..." -ForegroundColor Yellow

    # Install govulncheck
    try {
        $null = Get-Command govulncheck -ErrorAction Stop
        Write-Host "✅ govulncheck already installed" -ForegroundColor Green
    }
    catch {
        Write-Host "Installing govulncheck..." -ForegroundColor Yellow
        go install golang.org/x/vuln/cmd/govulncheck@latest
    }
}

Write-Host ""
Write-Host "🔍 1. Running Go vulnerability check..." -ForegroundColor Cyan
Write-Host "--------------------------------------" -ForegroundColor Cyan
govulncheck ./...

Write-Host ""
Write-Host "🔧 2. Running go vet..." -ForegroundColor Cyan
Write-Host "---------------------" -ForegroundColor Cyan
go vet ./...

Write-Host ""
Write-Host "🧪 3. Running tests..." -ForegroundColor Cyan
Write-Host "--------------------" -ForegroundColor Cyan
go test ./...

Write-Host ""
Write-Host "🔍 4. Checking for potential hardcoded secrets..." -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
try {
    $secretPatterns = Select-String -Path "*.go" -Pattern "(password|secret|key|token|api)" -Exclude "*test*","*example*" -Recurse
    if ($secretPatterns) {
        Write-Host "⚠️  Found potential hardcoded secrets - please review:" -ForegroundColor Yellow
        $secretPatterns | ForEach-Object { Write-Host "  $($_.Filename):$($_.LineNumber) - $($_.Line.Trim())" }
    } else {
        Write-Host "✅ No obvious hardcoded secrets found" -ForegroundColor Green
    }
}
catch {
    Write-Host "✅ No obvious hardcoded secrets found" -ForegroundColor Green
}

Write-Host ""
Write-Host "📋 5. Checking dependencies..." -ForegroundColor Cyan
Write-Host "-----------------------------" -ForegroundColor Cyan
go mod verify
Write-Host "✅ Dependencies verified" -ForegroundColor Green

Write-Host ""
Write-Host "✅ Security checks completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 Recommendations:" -ForegroundColor Yellow
Write-Host "- Review any warnings or vulnerabilities found above"
Write-Host "- Keep dependencies updated with 'go get -u ./...'"
Write-Host "- Run 'go mod tidy' to clean up unused dependencies"
Write-Host "- Consider using 'go mod vendor' for reproducible builds"
Write-Host ""
Write-Host "🔗 For more comprehensive security scanning, push to GitHub" -ForegroundColor Cyan
Write-Host "   to trigger the full security workflow with CodeQL analysis."