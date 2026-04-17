# driftwatch

> CLI tool to detect config drift between deployed services and source manifests

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Compare a deployed service against a local manifest:

```bash
driftwatch check --manifest ./deploy/api-service.yaml --env production
```

Watch for drift continuously:

```bash
driftwatch watch --manifest ./deploy/ --interval 60s
```

Example output:

```
[DRIFT DETECTED] api-service
  replicas:     expected=3  actual=1
  image:        expected=app:v1.4.2  actual=app:v1.3.9
  memory limit: expected=512Mi  actual=256Mi

[OK] worker-service — no drift detected
```

### Flags

| Flag | Description |
|------|-------------|
| `--manifest` | Path to source manifest file or directory |
| `--env` | Target environment to check against |
| `--interval` | Polling interval for `watch` mode |
| `--output` | Output format: `text`, `json`, `yaml` |

---

## License

MIT © 2024 yourusername