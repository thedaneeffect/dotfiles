# Bootstrap Checksum Verification Setup

This repository uses GitHub Actions to automatically update the bootstrap checksum in Cloudflare KV whenever `bootstrap.sh` changes.

## Architecture

```
bootstrap.sh modified → GitHub Actions → Calculate SHA256 → Update Cloudflare KV
                                                                      ↓
bootstrap alias → Download script → Fetch checksum → Verify → Execute
```

## Initial Setup

### 1. Create Cloudflare API Token

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. Create Token → Use "Edit Cloudflare Workers" template
3. Add additional permission: **Account.Workers KV Storage** (Edit)
4. Copy the token

### 2. Get Cloudflare Account & Namespace IDs

```bash
# Get Account ID (from Cloudflare dashboard URL)
# https://dash.cloudflare.com/<ACCOUNT_ID>/workers

# Get KV Namespace ID
wrangler kv:namespace list
# Copy the ID for your SECRETS namespace
```

### 3. Add GitHub Secrets

Add these secrets to your GitHub repository:
- Settings → Secrets and variables → Actions → New repository secret

Required secrets:
- `CLOUDFLARE_API_TOKEN` - The API token from step 1
- `CF_ACCOUNT_ID` - Your Cloudflare account ID
- `CF_NAMESPACE_ID` - Your KV namespace ID (from wrangler)

### 4. Initial Checksum Setup

Run this manually to set the initial checksum:

```bash
# Calculate current checksum
CHECKSUM=$(sha256sum bootstrap.sh | awk '{print $1}')

# Upload to Cloudflare KV
curl -X PUT "https://api.cloudflare.com/client/v4/accounts/$CF_ACCOUNT_ID/storage/kv/namespaces/$CF_NAMESPACE_ID/values/bootstrap:checksum" \
  -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  -H "Content-Type: text/plain" \
  --data "$CHECKSUM"
```

### 5. Deploy Worker

Deploy the updated worker with bootstrap endpoints:

```bash
cd worker
wrangler deploy
```

## How It Works

### Automatic Updates

When you push changes to `bootstrap.sh`:

1. GitHub Actions triggers on push to `main` branch
2. Calculates SHA256 of `bootstrap.sh`
3. Updates `bootstrap:checksum` in Cloudflare KV
4. Next time `bootstrap` alias runs, it gets the new checksum

### Verification Process

When you run `bootstrap` alias:

1. Downloads `bootstrap.sh` from worker
2. Fetches expected checksum from worker
3. Calculates actual checksum of downloaded script
4. Compares checksums
5. Only executes if match (or if checksum is "unknown")

### Endpoints

- `GET /bootstrap` - Returns bootstrap.sh with `X-Expected-SHA256` header
- `GET /bootstrap/checksum` - Returns just the checksum string

## Security Properties

✅ **Integrity verification** - Detects tampering with bootstrap script
✅ **Auto-updates** - Checksum updates automatically via CI/CD
✅ **Graceful degradation** - Works even if checksum not set (with warning)
✅ **No manual maintenance** - Just push changes, CI handles the rest

## Troubleshooting

### Checksum shows "unknown"

The initial checksum hasn't been set yet. Run the manual setup in step 4 above.

### GitHub Actions failing

Check that all three secrets are set correctly in your repository settings.

### Bootstrap verification failing

The script has been tampered with or there's a race condition between GitHub Actions updating the checksum and your download. Wait a minute and try again.

## Manual Verification

You can manually verify the bootstrap script:

```bash
# Download and check
curl -fsSL https://coffer.medieval.software/bootstrap > /tmp/bootstrap.sh
curl -fsSL https://coffer.medieval.software/bootstrap/checksum > /tmp/expected

# Calculate and compare
sha256sum /tmp/bootstrap.sh
cat /tmp/expected
```
