# Releasing the engine module

Gomag’s Go module lives at `engine/`, so published versions use **subdirectory tags**:

```text
engine/v0.1.0
```

Consumers install with:

```bash
go get github.com/just-Bri/gomakeagame/engine@v0.1.0
```

(Go strips the `engine/` prefix from the module version.)

## Workflows

### 1. Create release

**Actions → Create release → Run workflow**

| Input | Example | Notes |
|-------|---------|--------|
| `version` | `0.1.0` | Also accepts `v0.1.0` or `engine/v0.1.0` |
| `target` | `main` | Branch or SHA to tag |
| `prerelease` | false | Optional |

The job:

1. Resolves the previous **latest** GitHub release (else newest `engine/v*` tag)
2. Builds release notes from commits + diffstat since that point
3. Pushes tag `engine/vX.Y.Z` and creates a release with **`latest=false`**

Promotion to “Latest” is intentional and separate.

### 2. Mark release latest

**Actions → Mark release latest → Run workflow**

Enter an existing tag (e.g. `engine/v0.1.0`). That release becomes the repo’s latest.

Prereleases and drafts cannot be marked latest.
