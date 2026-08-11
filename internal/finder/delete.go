package finder

import "os"

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
