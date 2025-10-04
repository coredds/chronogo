# PowerShell script to set up the gh-pages branch for GitHub Pages deployment

Write-Host "Setting up gh-pages branch for GitHub Pages..." -ForegroundColor Green

try {
    # Save current branch
    $currentBranch = git branch --show-current
    
    # Create and switch to orphan gh-pages branch
    git checkout --orphan gh-pages
    
    # Remove all files from the working directory
    git rm -rf .
    
    # Copy docs content to root
    Copy-Item -Path "docs\*" -Destination "." -Recurse -Force
    Copy-Item -Path "docs\.nojekyll" -Destination "." -Force
    
    # Remove docs directory
    Remove-Item -Path "docs" -Recurse -Force
    
    # Add and commit files
    git add .
    git commit -m "Initial GitHub Pages deployment"
    
    # Push to remote
    git push -u origin gh-pages
    
    # Switch back to original branch
    git checkout $currentBranch
    
    Write-Host "✅ gh-pages branch set up successfully!" -ForegroundColor Green
    Write-Host "📝 Next steps:" -ForegroundColor Yellow
    Write-Host "   1. Go to your GitHub repository Settings"
    Write-Host "   2. Navigate to Pages section"
    Write-Host "   3. Set Source to 'Deploy from a branch'"
    Write-Host "   4. Select 'gh-pages' branch and '/ (root)' folder"
    Write-Host "   5. Save the settings"
    Write-Host ""
    Write-Host "🌐 Your site will be available at: https://coredds.github.io/chronogo/" -ForegroundColor Cyan
}
catch {
    Write-Host "❌ Error occurred: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Switching back to original branch..." -ForegroundColor Yellow
    git checkout $currentBranch
}
