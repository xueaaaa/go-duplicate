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

func Find(ctx context.Context, dir string, params params.Params) ([]DuplicateGroup, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	bar := util.NewBar(params.Silent, -1, "scanning directory...")

	files := make(chan file.File)
	errors := make(chan error, 1)
	go func() {
		defer close(files)
		errors <- file.Scan(ctx, dir, files)
	}()

	sizes := make(map[int64][]file.File)
	for f := range files {
		if !f.IsHardlink {
			sizes[f.Size] = append(sizes[f.Size], f)
		}

		if bar != nil {
			if err := bar.Add(1); err != nil {
				return nil, err
			}
		}
	}

	if err := <-errors; err != nil {
		return nil, err
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
					return nil, err
				}
			}
			if err != nil {
				if util.IsSkippableFSError(err) {
					continue
				}
				return nil, err
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
					return nil, err
				}
			}
			if err != nil {
				if util.IsSkippableFSError(err) {
					continue
				}
				return nil, err
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

	return groups, nil
}
