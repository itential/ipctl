// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package utils

func StringInsert(slice []string, index int, value string) {
	if index < 0 || index > len(slice) {
		panic("Index out of range")
	}

	// Grow the slice by one element.
	slice = append(slice, "")

	// Copy the elements after the insertion point to the right.
	copy(slice[index+1:], slice[index:])

	// Insert the new value.
	slice[index] = value

	//return slice
}

func StringExists(slice []string, name string) bool {
	for _, ele := range slice {
		if ele == name {
			return true
		}
	}
	return false
}
