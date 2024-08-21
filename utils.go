package main

func GetSliceDiff(slice1, slice2 []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range slice2 {
		m[item] = true
	}

	for _, item := range slice1 {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}
	return
}
