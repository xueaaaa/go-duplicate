package params

// Params holds user-configurable options controlling how a scan is run
// and how its results are used.
type Params struct {
	// Format selects the output format and should be one of the format
	// identifiers defined alongside it.
	Format Format
	// SampleSize is the sample size in bytes passed to [finder.Find]
	// (and from there to [hash.Fingerprint]); it should be positive —
	// see [hash.Fingerprint] for the consequences of zero or negative
	// values.
	SampleSize int64
	// Confirm, if true, requires interactive confirmation (see
	// [PlainConfirm]) before Delete or Hardlink is carried out.
	Confirm bool
	// Delete, if true, removes redundant duplicate files after a scan
	// (see [finder.Delete]).
	Delete bool
	// Hardlink, if true, replaces redundant duplicate files with hard
	// links to the kept copy after a scan (see [finder.Hardlink]).
	Hardlink bool
	// Silent suppresses progress bars during scanning (see [finder.Find]).
	Silent bool
}
