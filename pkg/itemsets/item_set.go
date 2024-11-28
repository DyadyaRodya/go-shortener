package itemsets

import "slices"

// AddItems adds each item of items to provided slice if item not in slice.
func AddItems(slice, items []string) []string {
	for _, item := range items {
		if !slices.Contains(slice, item) {
			slice = append(slice, item)
		}
	}
	return slice
}

// RemoveItems keeps in result slice only those items of slice that are not presented in items.
func RemoveItems(slice, items []string) []string {
	resultItems := []string{}
	for _, item := range slice {
		if !slices.Contains(items, item) {
			resultItems = append(resultItems, item)
		}
	}
	return resultItems
}

// Intersection keeps in result slice only those items that are presented in both slice and items.
func Intersection(slice1, slice2 []string) []string {
	result := make([]string, 0)
	for _, s1 := range slice1 {
		if slices.Contains(slice2, s1) {
			result = append(result, s1)
		}
	}
	return result
}
