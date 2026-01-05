# Railway Deployment Configuration

This directory contains the canonical Railway deployment configuration files.

## Files

- `railway.json` - Railway deployment configuration
- `nixpacks.toml` - Nixpacks build configuration

## Important Notes

**Railway requires these configuration files to be in the repository root directory.**

The files in this directory are the source of truth and should be copied to the root
when making changes:

```bash
cp deployments/railway/railway.json .
cp deployments/railway/nixpacks.toml .
```

### Build Command

The correct build command is:
```bash
go build -o main ./cmd/api
```

This builds the application from `cmd/api/main.go`, not from the repository root.

### Configuration Changes

When modifying Railway configuration:
1. Edit the files in this directory (`deployments/railway/`)
2. Copy the updated files to the repository root
3. Commit both versions to keep them in sync
