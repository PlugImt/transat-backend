package internal

import (
	"strings"
	"unicode"
)

// NormalizeItemName normalizes menu item names for consistent comparison
func NormalizeItemName(name string) string {
	normalized := strings.TrimSpace(name)
	normalized = strings.ToLower(normalized)
	normalized = strings.Join(strings.Fields(normalized), " ")

	// Remove common punctuation but keep essential characters
	result := strings.Builder{}
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}

	return strings.TrimSpace(result.String())
}

// CapitalizeItemName capitalizes only the first character of the item name
func CapitalizeItemName(name string) string {
	if name == "" {
		return name
	}

	runes := []rune(name)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}

// CalculateSimilarity calculates similarity between two sets of strings
func CalculateSimilarity(set1, set2 []string) float64 {
	if len(set1) == 0 && len(set2) == 0 {
		return 1.0
	}

	map1 := make(map[string]bool)
	map2 := make(map[string]bool)

	for _, item := range set1 {
		map1[item] = true
	}

	for _, item := range set2 {
		map2[item] = true
	}

	intersection := 0
	for item := range map1 {
		if map2[item] {
			intersection++
		}
	}

	union := len(map1) + len(map2) - intersection

	if union == 0 {
		return 1.0
	}

	return float64(intersection) / float64(union)
}
