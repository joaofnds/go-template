---
name: update-deps
description: Update every dependency in this project to its latest stable version — Go toolchain, mise tools, go modules, Docker images, GitHub Actions — then make sure it still builds and tests pass. Use when the user says "update deps", "upgrade everything", "bring deps up to date", or any variant.
user-invocable: true
allowed-tools:
  - Bash
  - Read
  - Edit
---

# /update-deps

Bring every version pin in the project to latest stable. The Go version appears in three places (`mise.toml`, `go.mod`, `Dockerfile`) and must stay identical — bump it first, propagate, then everything else. CI uses `jdx/mise-action`, so it picks up the new Go version from `mise.toml` automatically.

## Steps

1. **Bump Go via mise, propagate to `go.mod` and `Dockerfile`.**

   ```bash
   mise upgrade go --bump --local
   ```

   Then edit `go.mod` (`go` directive, and `toolchain` line if present) and `Dockerfile` (`FROM golang:...`) to the same version mise just wrote into `mise.toml`.

2. **Bump the rest of the mise tools.**

   ```bash
   mise upgrade --bump --local --exclude go
   ```

3. **Update Go modules.**

   ```bash
   go get -u ./...
   go mod tidy
   ```

4. **Bump Docker image tags.** For each pinned image in `Dockerfile` and `compose*.yaml`, look up the latest stable tag and edit. `gcr.io/distroless/static:nonroot` is a rolling tag — leave it alone.

   - Docker Hub: `curl -s 'https://registry.hub.docker.com/v2/repositories/<image>/tags/?page_size=100' | jq -r '.results[].name'` and pick the highest version matching the existing tag shape (e.g., `N-alpine`). Official images live under `library/` — use `library/postgres`, `library/redis`.
   - GitHub-released images (e.g. `casbin/casdoor`): `curl -s https://api.github.com/repos/<owner>/<repo>/releases/latest | jq -r .tag_name`.

5. **Bump GitHub Actions.** For each `uses: owner/repo@vN` line in `.github/workflows/*.yaml`:

   ```bash
   curl -s https://api.github.com/repos/<owner>/<repo>/releases/latest | jq -r .tag_name
   ```

   Keep the existing pin shape (major-only `@v4`, full semver `@v1.4.0`, or full SHA).

6. **Build, fix, test.**

   ```bash
   go build ./... && go vet ./...
   docker compose up -d && task migrate-up && task test
   ```

   `task test` needs the services up and `CONFIG_PATH` from the mise env, so run it inside `mise`. If anything fails because of an upgrade, fix the call sites — don't revert the bump. Library breakages are an expected part of this skill; adapt the code to the new API.

Never downgrade or pick a pre-release: if a `releases/latest` lookup returns a version behind the current pin, leave it and note it.

Report what changed and what (if anything) needed code changes to accommodate.
