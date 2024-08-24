package wtsc

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
)

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

// IsEmptyString checks if a string is empty by comparing it to an empty string.
func IsEmptyString(str string) bool {
	return str == ""
}

func IsValidHttpURI(uri string) bool {
	if IsEmptyString(uri) {
		return false
	}

	parsedURL, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	// Check if the scheme is HTTP or HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if the host is not empty
	if parsedURL.Host == "" {
		return false
	}

	return true
}

func IsValidDomain(domain string) bool {
	// Check the length of the entire domain
	if len(domain) > 253 {
		return false
	}

	// Define the regex for validating each label
	var labelRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)

	// Split the domain into labels and validate each one
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 || !labelRegex.MatchString(label) {
			return false
		}
	}

	return true
}

func IsWritableDirectory(path string) bool {
	// Check if the path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	// Check if the path is a directory
	if !info.IsDir() {
		return false
	}

	// Check write permission by attempting to create a temporary file
	testFile := fmt.Sprintf("%s/.write_test", path)
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	err = file.Close()
	if err != nil {
		return false
	}
	err = os.Remove(testFile)
	if err != nil {
		return false
	}

	return true
}
