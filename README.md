# gitfame

`gitfame` is a CLI utility that calculates per-author statistics for a Git repository at a given revision.

## Features

- Line count (`Lines`)
- Unique commit count (`Commits`)
- Unique file count (`Files`)
- Author or committer mode (`--use-committer`)
- Filters: extensions, languages, exclude glob, restrict-to glob
- Output formats: `tabular`, `csv`, `json`, `json-lines`

## Build

```bash
(cd cmd/gitfame && go build .)
```

## Install

```bash
go install ./cmd/gitfame/...
```

## Run

```bash
gitfame --repository=. --revision=HEAD --order-by=lines --format=tabular
```

Example with filters:

```bash
gitfame --repository=. --extensions='.go,.md' --languages='go,markdown' --exclude='vendor/*'
```

## Tests

```bash
go test ./...
```

```bash
go test -v ./gitfame/test/integration/...
```

## Notes

- Unknown languages in `--languages` are ignored and reported as a warning in `stderr`.
- Invalid flag values (including invalid glob patterns) return a non-zero exit code.
- The tool does not modify repository state.
