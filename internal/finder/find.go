package finder

import (
	"cmp"
	"context"
	"go-duplicate/internal/file"
	hash2 "go-duplicate/internal/hash"
	"go-duplicate/internal/util"
	"slices"
)

func Find(ctx context.Context, dir string, params int64 /** TODO REPLACE WITH PARAMS STRUCT **/) ([]DuplicateGroup, error) {
	files := make(chan file.File)
	errors := make(chan error, 1)
	go func() {
		defer close(files)
		errors <- file.Scan(ctx, dir, files)
	}()

	sizes := make(map[int64][]file.File)
	for f := range files {
		sizes[f.Size] = append(sizes[f.Size], f)
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	fingerprints := make(map[[32]byte][]file.File)
	for _, v := range sizes {
		if len(v) <= 1 {
			continue
		}

		for _, f := range v {
			hash, err := hash2.Fingerprint(f, params)
			if err != nil {
				if util.IsSkippableFSError(err) {
					continue
				}
				return nil, err
			}

			fingerprints[[32]byte(hash)] = append(fingerprints[[32]byte(hash)], f)
		}
	}

	groups := make([]DuplicateGroup, 0)
	for k, v := range fingerprints {
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
