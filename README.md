# envdiff

A CLI tool to compare `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <file1> <file2>
```

### Example

```bash
envdiff .env.development .env.production
```

**Output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED keys:
  - APP_ENV: "development" vs "production"
  - DEBUG: "true" vs "false"
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Only report missing keys, skip value comparison |
| `--quiet` | Exit with non-zero status if differences found (useful for CI) |
| `--format json` | Output results as JSON |

### CI Integration

```bash
envdiff --quiet .env.example .env && echo "Env files are in sync"
```

---

## License

MIT © [yourusername](https://github.com/yourusername)