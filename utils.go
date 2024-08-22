package wtsc

import "reflect"

func GetStrSliceDiff(slice1, slice2 []string) (diff []string) {
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

func GetServiceStakeSliceDiff(slice1, slice2 []ServiceStake) (updated, removed, added []ServiceStake) {
	map1 := make(map[string]ServiceStake)
	map2 := make(map[string]ServiceStake)

	// Populate map1 with slice1 elements
	for _, item := range slice1 {
		map1[item.Service] = item
	}

	// Populate map2 with slice2 elements
	for _, item := range slice2 {
		map2[item.Service] = item
	}

	// Check for updated and removed services
	for service, item1 := range map1 {
		if item2, found := map2[service]; found {
			if !reflect.DeepEqual(item1, item2) {
				updated = append(updated, item2)
			}
		} else {
			removed = append(removed, item1)
		}
	}

	// Check for added services
	for service, item2 := range map2 {
		if _, found := map1[service]; !found {
			added = append(added, item2)
		}
	}

	return
}

// FindStringInSlice checks if a string is present in a slice of strings
func FindStringInSlice(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
