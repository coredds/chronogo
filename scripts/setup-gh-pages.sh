#!/bin/bash

# Script to set up the gh-pages branch for GitHub Pages deployment
# This script creates an orphan gh-pages branch and pushes the docs content

set -e

echo "Setting up gh-pages branch for GitHub Pages..."

# Save current branch
CURRENT_BRANCH=$(git branch --show-current)

# Create and switch to orphan gh-pages branch
git checkout --orphan gh-pages

# Remove all files from the working directory
git rm -rf .

# Copy docs content to root
cp -r docs/* .
cp docs/.nojekyll .

# Remove docs directory
rm -rf docs

# Add and commit files
git add .
git commit -m "Initial GitHub Pages deployment"

# Push to remote
git push -u origin gh-pages

# Switch back to original branch
git checkout "$CURRENT_BRANCH"

echo "✅ gh-pages branch set up successfully!"
echo "📝 Next steps:"
echo "   1. Go to your GitHub repository Settings"
echo "   2. Navigate to Pages section"
echo "   3. Set Source to 'Deploy from a branch'"
echo "   4. Select 'gh-pages' branch and '/ (root)' folder"
echo "   5. Save the settings"
echo ""
echo "🌐 Your site will be available at: https://coredds.github.io/ChronoGo/"
