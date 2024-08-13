package stamp

type VerbListManager interface {
	IsVerbStamp(stem string) bool
}

// verbListManager manages the list of verbs.
type verbListManager struct {
	verbList       []string
	verbNormalizer VerbNormalizer
}

// NewVerbListManager creates a new instance of VerbListManager with the provided initial verb list and VerbNormalizer.
// It initializes the verb list by normalizing the provided verbs using the VerbNormalizer.
func NewVerbListManager(initialVerbList []string, verbNormalizer VerbNormalizer) VerbListManager {
	vlm := &verbListManager{
		verbNormalizer: verbNormalizer,
	}
	vlm.initializeVerbList(initialVerbList)
	return vlm
}

// initializeVerbList normalizes each verb in the initial verb list and appends it to the internal verb list.
// This method is called during the creation of the VerbListManager instance.
func (vlm *verbListManager) initializeVerbList(initialVerbList []string) {
	for _, verb := range initialVerbList {
		normalizedVerb := vlm.verbNormalizer.Normalize(verb)
		vlm.verbList = append(vlm.verbList, normalizedVerb)
	}
}

// IsVerbStamp checks if the normalized version of the given stem is present in the verb list.
// It returns true if the normalized stem is found in the list, false otherwise.
func (vlm *verbListManager) IsVerbStamp(stem string) bool {
	normalizedStem := vlm.verbNormalizer.Normalize(stem)
	for _, verb := range vlm.verbList {
		if verb == normalizedStem {
			return true
		}
	}
	return false
}
