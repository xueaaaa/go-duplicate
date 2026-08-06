package units

import "go-duplicate/internal/finder"

func CalculatePotentialSaveSpace(groups []finder.DuplicateGroup) int64 {
	calculated := int64(0)

	for _, group := range groups {
		calculated += int64(len(group.Files)-1) * group.FileSize
	}

	return calculated
}
