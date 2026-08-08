package units

import "fmt"

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
