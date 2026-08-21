package finder

import "os"

// Delete removes all files from each duplicate group, leaving the first one (group.Files[0]).
// Returns the number of files deleted.
//
// Delete stops at the first deletion error and returns it along with the
// number of files successfully deleted so far. The remaining files in the
// current group, as well as any subsequent groups, are left untouched.
func Delete(groups []DuplicateGroup) (int, error) {
	count := 0

	for _, v := range groups {
		for _, file := range v.Files[1:] {
			if err := os.Remove(file.Path); err != nil {
				return count, err
			}
			count++
		}
	}

	return count, nil
}
