package roots

import (
	"github.com/berkayersoyy/go-arabic-light-stemmer/arabic/constant"
	"github.com/berkayersoyy/go-arabic-light-stemmer/arabic/utils"
	"strings"
)

type RootsManager interface {
	IsRoot(word string) bool
	NormalizeRoot(word string) string
	MostCommon(lst []string) string
	FilterRootLengthValid(roots []string) []string
	LookupRoots(roots []string) []string
	ChooseRoot(affixationList []map[string]string) string
}

type rootsManager struct {
	roots map[string]bool
}

// NewRootsManager creates a new instance of rootsManager with the provided roots map.
func NewRootsManager() RootsManager {
	roots := make(map[string]bool)
	for _, root := range constant.ROOTS {
		roots[root] = true
	}
	return &rootsManager{roots: roots}
}

// IsRoot checks if a given word exists as a root in the dictionary.
func (r *rootsManager) IsRoot(word string) bool {
	_, exists := r.roots[word]
	return exists
}

// NormalizeRoot normalizes a given root word by replacing or removing specific characters.
func (r *rootsManager) NormalizeRoot(word string) string {
	word = strings.ReplaceAll(word, constant.ALEF_MADDA, constant.HAMZA+constant.ALEF)
	word = strings.ReplaceAll(word, constant.TEH_MARBUTA, "")
	word = strings.ReplaceAll(word, constant.ALEF_MAKSURA, constant.YEH)
	return utils.NormalizeHamza(word)
}

// MostCommon finds and returns the most common string in a given list.
func (r *rootsManager) MostCommon(lst []string) string {
	counts := make(map[string]int)
	for _, item := range lst {
		counts[item]++
	}

	var mostCommon string
	maxCount := 0
	for item, count := range counts {
		if count > maxCount {
			mostCommon = item
			maxCount = count
		}
	}

	return mostCommon
}

// FilterRootLengthValid filters a list of roots, returning only those that have a valid length (3-4 characters)
// and do not contain the ALEF character.
func (r *rootsManager) FilterRootLengthValid(roots []string) []string {
	var validRoots []string
	for _, root := range roots {
		runeRoot := []rune(root)
		if len(runeRoot) >= 3 && len(runeRoot) <= 4 && !strings.Contains(root, constant.ALEF) {
			validRoots = append(validRoots, root)
		}
	}
	return validRoots
}

// LookupRoots checks a list of roots against the dictionary and returns only the roots that exist in the dictionary.
func (r *rootsManager) LookupRoots(roots []string) []string {
	var accepted []string
	for _, root := range roots {
		if r.IsRoot(root) {
			accepted = append(accepted, root)
		}
	}
	return accepted
}

// ChooseRoot selects the most suitable root from a list of affixations.
// It normalizes the roots and stems, filters them, and then checks them against the dictionary.
// If debug is enabled, additional information could be logged or displayed.
func (r *rootsManager) ChooseRoot(affixationList []map[string]string) string {
	var accepted, allAccepted []string

	stems := make([]string, len(affixationList))
	roots := make([]string, len(affixationList))
	prefixes := make([]string, len(affixationList))
	suffixes := make([]string, len(affixationList))

	for i, affixation := range affixationList {
		stems[i] = r.NormalizeRoot(affixation["stem"])
		roots[i] = r.NormalizeRoot(affixation["root"])
		prefixes[i] = affixation["prefix"]
		suffixes[i] = affixation["suffix"]
	}

	// Check roots
	rootsTmp := r.FilterRootLengthValid(roots)
	accepted = r.LookupRoots(rootsTmp)
	allAccepted = append(allAccepted, accepted...)

	// Check stems as roots
	stemsTmp := r.FilterRootLengthValid(stems)
	accepted = r.LookupRoots(stemsTmp)
	allAccepted = append(allAccepted, accepted...)

	if len(allAccepted) > 0 {
		return r.MostCommon(allAccepted)
	} else if len(accepted) > 0 {
		return r.MostCommon(accepted)
	}
	return ""
}
