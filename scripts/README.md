# Scripts

This directory contains utility scripts for the ChronoGo project.

## GitHub Pages Setup

### `setup-gh-pages.sh` (Bash)
Sets up the `gh-pages` branch for GitHub Pages deployment on Unix-like systems.

```bash
chmod +x scripts/setup-gh-pages.sh
./scripts/setup-gh-pages.sh
```

### `setup-gh-pages.ps1` (PowerShell)
Sets up the `gh-pages` branch for GitHub Pages deployment on Windows.

```powershell
./scripts/setup-gh-pages.ps1
```

## What these scripts do:

1. Create an orphan `gh-pages` branch
2. Copy docs content to the branch root
3. Commit and push the branch
4. Switch back to your original branch

## Automatic Deployment

Once the `gh-pages` branch is set up, the GitHub Action in `.github/workflows/deploy-gh-pages.yml` will automatically update it whenever changes are made to the `docs/` folder on the main branch.
