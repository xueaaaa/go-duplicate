package units

import "go-duplicate/internal/finder"

// CalculatePotentialSaveSpace returns the total disk space, in bytes, that
// could be reclaimed by keeping only one file per group and removing the
// rest. For each group this is (len(Files)-1) * FileSize — i.e. it
// excludes the one file that would be kept, unlike [finder.DuplicateGroup]
// totals computed elsewhere (e.g. [GroupJSON.TotalSizeBytes]) which include
// it.
//
// Groups with fewer than 2 files contribute 0 (though [finder.Find] never
// produces such groups). An empty groups slice returns 0.
func CalculatePotentialSaveSpace(groups []finder.DuplicateGroup) int64 {
	calculated := int64(0)

	for _, group := range groups {
		calculated += int64(len(group.Files)-1) * group.FileSize
	}

	return calculated
}
