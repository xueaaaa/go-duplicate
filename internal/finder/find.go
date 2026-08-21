package finder

import (
	"cmp"
	"context"
	"fmt"
	"go-duplicate/internal/file"
	hash2 "go-duplicate/internal/hash"
	"go-duplicate/internal/params"
	"go-duplicate/internal/util"
	"os"
	"slices"
)

// Find scans dir for duplicate files and groups them by content.
//
// Find works in three narrowing passes: files are first grouped by size,
// then candidate groups (size > 1) are grouped by a fingerprint computed
// from a sample of each file ([hash2.Fingerprint], params.SampleSize bytes),
// then remaining candidates are grouped by a full-content hash
// ([hash2.Hash]). Only files that share size, fingerprint, and hash end up
// in the same [DuplicateGroup]. Hardlinks (file.File.IsHardlink) are
// excluded before the first pass and never appear in the result.
//
// Unless params.Silent is set, Find prints progress bars for each pass to
// os.Stderr.
//
// Find returns the total number of files scanned (regardless of outcome)
// and, on success, the duplicate groups sorted by descending file size.
//
// Errors from the initial directory scan are always fatal. During the
// fingerprint and hash passes, errors satisfying util.IsSkippableFSError
// cause the affected file to be skipped rather than aborting the scan; any
// other error is fatal. On a fatal error, Find returns the files-scanned
// count seen so far, a nil group slice, and the error.
func Find(ctx context.Context, dir string, params params.Params) (int64, []DuplicateGroup, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	bar := util.NewBar(params.Silent, -1, "scanning directory...")

	files := make(chan file.File)
	errors := make(chan error, 1)
	go func() {
		defer close(files)
		errors <- file.Scan(ctx, dir, files)
	}()

	filesCount := int64(0)
	sizes := make(map[int64][]file.File)
	for f := range files {
		if !f.IsHardlink {
			sizes[f.Size] = append(sizes[f.Size], f)
		}
		filesCount++

		if bar != nil {
			if err := bar.Add(1); err != nil {
				return filesCount, nil, err
			}
		}
	}

	if err := <-errors; err != nil {
		return filesCount, nil, err
	}

	if bar != nil {
		bar.Finish()
		fmt.Fprintln(os.Stderr)
	}

	fingerprintTotal := 0
	for _, v := range sizes {
		if len(v) > 1 {
			fingerprintTotal += len(v)
		}
	}

	bar = util.NewBar(params.Silent || fingerprintTotal == 0, fingerprintTotal, "calculating fingerprints...")

	fingerprints := make(map[[32]byte][]file.File)
	for _, v := range sizes {
		if len(v) <= 1 {
			continue
		}

		for _, f := range v {
			fingerprint, err := hash2.Fingerprint(f, params.SampleSize)
			if bar != nil {
				if err = bar.Add(1); err != nil {
					return filesCount, nil, err
				}
			}
			if err != nil {
				if util.IsSkippableFSError(err) {
					continue
				}
				return filesCount, nil, err
			}

			var key [32]byte
			copy(key[:], fingerprint)
			fingerprints[key] = append(fingerprints[key], f)
		}
	}

	if bar != nil {
		bar.Finish()
		fmt.Fprintln(os.Stderr)
	}

	hashesTotal := 0
	for _, v := range fingerprints {
		if len(v) > 1 {
			hashesTotal += len(v)
		}
	}

	bar = util.NewBar(params.Silent || hashesTotal == 0, hashesTotal, "calculating hashes...")

	hashes := make(map[[32]byte][]file.File)
	for _, v := range fingerprints {
		if len(v) <= 1 {
			continue
		}

		for _, f := range v {
			hash, err := hash2.Hash(f)
			if bar != nil {
				if err = bar.Add(1); err != nil {
					return filesCount, nil, err
				}
			}
			if err != nil {
				if util.IsSkippableFSError(err) {
					continue
				}
				return filesCount, nil, err
			}

			var key [32]byte
			copy(key[:], hash)
			hashes[key] = append(hashes[key], f)
		}
	}

	if bar != nil {
		bar.Finish()
		fmt.Fprintln(os.Stderr)
	}

	groups := make([]DuplicateGroup, 0)
	for k, v := range hashes {
		if len(v) <= 1 {
			continue
		}

		group := DuplicateGroup{
			Hash:     k,
			FileSize: v[0].Size,
			Files:    v,
		}
		groups = append(groups, group)
	}

	slices.SortFunc(groups, func(a, b DuplicateGroup) int {
		return cmp.Compare(b.FileSize, a.FileSize)
	})

	return filesCount, groups, nil
}
