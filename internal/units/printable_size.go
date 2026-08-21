package units

import "fmt"

// PrintableSize formats size (in bytes) as a human-readable string using
// binary units (KiB/MiB/GiB/TiB, powers of 1024) — not decimal (KB/MB/GB,
// powers of 1000). The largest unit for which size >= unit is used (so
// size == 1024 prints as "1.00 KiB", not "1024 B"), with two decimal
// places; sizes below 1 KiB print as a plain byte count ("512 B").
//
// Negative values are not specially handled: they fall through to the
// byte format (e.g. "-5 B") without converting to a larger unit.
func PrintableSize(size int64) string {
	switch {
	case size >= TiB:
		return fmt.Sprintf("%.2f TiB", float64(size)/float64(TiB))
	case size >= GiB:
		return fmt.Sprintf("%.2f GiB", float64(size)/float64(GiB))
	case size >= MiB:
		return fmt.Sprintf("%.2f MiB", float64(size)/float64(MiB))
	case size >= KiB:
		return fmt.Sprintf("%.2f KiB", float64(size)/float64(KiB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
