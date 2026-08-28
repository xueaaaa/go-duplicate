# go-duplicate

![Go Version](https://img.shields.io/github/go-mod/go-version/xueaaaa/go-duplicate)
[![Go Report Card](https://goreportcard.com/badge/github.com/xueaaaa/go-duplicate)](https://goreportcard.com/report/github.com/xueaaaa/go-duplicate)
[![Go Reference](https://pkg.go.dev/badge/github.com/xueaaaa/go-duplicate.svg)](https://pkg.go.dev/github.com/xueaaaa/go-duplicate)
![License](https://img.shields.io/github/license/xueaaaa/go-duplicate)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-lightgrey)

A CLI tool for finding and cleaning up duplicate files by content.

It scans a directory tree and groups files by size. Then confirms real 
duplicates using a fast partial fingerprint first, then a full SHA-256 
comparison, so two different files are never reported as duplicates.

## Features

- Recursive directory scan with size-based pre-grouping
- Two-phase hashing to avoid false positives
- Plain-text or JSON output
- Delete redundant copies, or replace them with hard links
- Detects files that are already hard-linked

## Installation

```
go install github.com/xueaaaa/go-duplicate@latest
```

Or build from source:

```
git clone https://github.com/xueaaaa/go-duplicate.git
cd go-duplicate
go build -o .
```

## Requirements

- Go 1.26+ (see `go.mod`)
- Any Unix-like OS

## Usage

```
go-duplicate <directory> [flags]
```

### Basic scan

#### Plain output

```
go-duplicate find /home/user/photos
```

Prints a human-readable report: how many files were scanned, how 
many duplicate groups were found, and how much space could be reclaimed.

**Output example:**

```
Scanning complete (directory: /home/user/photos)

Scanning complete at: 2026-06-07 14:58:40.112233445 +0300 MSK m=+0.652412869

Files scanned: X
Duplicate groups: X
Duplicate files: X
Potential space savings: X KiB

[1]
File size: X KiB
Summary size: X KiB
Hash: X
Files:
[1] - /home/xueaaaa/photos/photo1.jpg
[2] - /home/xueaaaa/photos/some_photos/photoX.jpg

[2]
...
```

#### JSON output

```
go-duplicate find /home/user/photos --format json
```

**Output example:**

```
{
  "schema_version": 1,
  "stats": {
    "scanned_dir": "/home/user/photos",
    "scanned_at": "2026-06-07T14:58:40+03:00",
    "files_scanned": X,
    "duplicate_groups": X,
    "duplicate_files": X,
    "potential_saving_bytes": X
  },
  "groups": [
    {
      "hash": "X",
      "file_size_bytes": X,
      "files": [
        { "path": "/home/xueaaaa/photos/photo1.jpg" },
        { "path": "/home/xueaaaa/photos/some_photos/photoX.jpg" }
      ]
    }
  ]
}
```

## Deleting duplicates

```
godup find /home/user/photos --delete
```

Keeps one file per duplicate group (the one with the lexicographically 
first path) and removes the rest. Without -y, you'll be asked to confirm 
before anything is deleted.

```
godup find /home/user/photos --delete -y
```

Skips the confirmation prompt — use this only once you're confident in the 
result, or in non-interactive scripts.

## Replacing duplicate with hard links

```
godup find /home/user/photos --hardlink
```

Instead of deleting redundant copies, replaces them with hard links to the 
kept file. This reclaims disk space without actually removing any path — every 
duplicate path still exists and still works, but they all point to the same data on disk.

> After `--hardlink`, all linked paths share the same underlying data. Editing or truncating 
> any one of them will affect the others. Do not use `--hardlink` on files you may want to 
> independently modify later.

## Flags

| **Flag**      | **Short** | **Default value** | **Description**                                                                             |
|---------------|-----------|-------------------|---------------------------------------------------------------------------------------------|
| --format      | -f        | plain             | Output format (`plain` or `json`)                                                           |
| --sample-size | -s        | 8192              | Size in bytes of the initial and final parts used to calculate the partial hash/fingerprint |
| --delete      | -d        | false             | Delete duplicate files, keeping only the first file in each group                           |
| --hardlink    | -h        | false             | Converts duplicate files into hardlinks to the first file in each group                     |
| --yes         | -y        | false             | Skip confirmation prompt (only applies with `--delete` or `--hardlink`)                     |
| --silent      |           | false             | Disable the progress bar                                                                    |

## Known limitations

- Hashing is currently single-threaded
- Only one directory can be scanned per invocation
- There is no `--exclude`/`--include` filtering yet
- go-duplicate does not currently follow symlinks; symlinked files are skipped

## Roadmap

- `--workers` for parallel hashing for large scans
- Multiple directories per invocation
- `--exclude`/`--min-size` filters

## Contributing

Issues and pull requests are welcome. Please run `go vet ./...` and `go test ./...` before submitting.

## License

See [LICENSE](LICENSE) for details.

![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?logo=go&logoColor=white)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)