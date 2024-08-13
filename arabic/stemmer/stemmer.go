package stemmer

import (
	"fmt"
	"go-arabic-stemmer/arabic/constant"
	"go-arabic-stemmer/arabic/roots"
	"go-arabic-stemmer/arabic/stamp"
	"go-arabic-stemmer/arabic/stop_words"
	"go-arabic-stemmer/arabic/utils"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// ArabicLightStemmer defines a stemmer with configurable parameters.
type ArabicLightStemmer struct {
	stopWordManager  stop_words.StopwordManager
	wordProcessor    stop_words.WordProcessor
	tashkeelChecker  stop_words.TashkeelChecker
	verbListManager  stamp.VerbListManager
	verbNormalizer   stamp.VerbNormalizer
	rootsManager     roots.RootsManager
	prefixLetters    string
	suffixLetters    string
	infixLetters     string
	maxPrefixLength  int
	maxSuffixLength  int
	minStemLength    int
	joker            string
	prefixList       []string
	suffixList       []string
	rootList         []string
	validAffixesList []string
	tokenPat         *regexp.Regexp
	prefixesTree     map[string]interface{}
	suffixesTree     map[string]interface{}
}

// NewArabicLightStemmer creates a new instance of ArabicLightStemmer with default values.
func NewArabicLightStemmer() *ArabicLightStemmer {
	affixList := append([]string{}, constant.NOUN_AFFIX_LIST...)
	affixList = append(affixList, constant.VERB_AFFIX_LIST...)

	tashkeelChecker := stop_words.NewTashkeelChecker()
	wordProcessor := stop_words.NewWordProcessor(tashkeelChecker)
	stopWordManager := stop_words.NewStopwordManager(wordProcessor)
	verbNormalizer := stamp.NewVerbNormalizer(wordProcessor)
	verbListManager := stamp.NewVerbListManager(stamp.INITIAL_VERB_LIST, verbNormalizer)
	rootsManager := roots.NewRootsManager()
	stemmer := &ArabicLightStemmer{
		stopWordManager:  stopWordManager,
		wordProcessor:    wordProcessor,
		tashkeelChecker:  tashkeelChecker,
		verbListManager:  verbListManager,
		verbNormalizer:   verbNormalizer,
		rootsManager:     rootsManager,
		prefixLetters:    constant.DEFAULT_PREFIX_LETTERS,
		suffixLetters:    constant.DEFAULT_SUFFIX_LETTERS,
		infixLetters:     constant.DEFAULT_INFIX_LETTERS,
		maxPrefixLength:  constant.DEFAULT_MAX_PREFIX,
		maxSuffixLength:  constant.DEFAULT_MAX_SUFFIX,
		minStemLength:    constant.DEFAULT_MIN_STEM,
		joker:            constant.DEFAULT_JOKER,
		prefixList:       constant.DEFAULT_PREFIX_LIST,
		suffixList:       constant.DEFAULT_SUFFIX_LIST,
		rootList:         constant.ROOTS,
		validAffixesList: affixList,
		tokenPat:         regexp.MustCompile(`[^\w\x{064b}-\x{0652}']+`),
		prefixesTree:     make(map[string]interface{}),
		suffixesTree:     make(map[string]interface{}),
	}

	// Initialize prefix and suffix trees
	stemmer.prefixesTree = stemmer.createPrefixTree()
	stemmer.suffixesTree = stemmer.createSuffixTree()

	return stemmer
}

// SetPrefixLetters sets the prefix letters used in the stemming process.
// The prefix letters define the characters or sequences of characters that may appear at the beginning of words.
func (als *ArabicLightStemmer) SetPrefixLetters(newPrefixLetters string) {
	als.prefixLetters = newPrefixLetters
}

// GetPrefixLetters returns the current prefix letters used in the stemming process.
// These letters are used to identify and remove prefixes from words during the stemming process.
func (als *ArabicLightStemmer) GetPrefixLetters() string {
	return als.prefixLetters
}

// SetSuffixLetters sets the suffix letters used in the stemming process.
// The suffix letters define the characters or sequences of characters that may appear at the end of words.
func (als *ArabicLightStemmer) SetSuffixLetters(newSuffixLetters string) {
	als.suffixLetters = newSuffixLetters
}

// GetSuffixLetters returns the current suffix letters used in the stemming process.
// These letters are used to identify and remove suffixes from words during the stemming process.
func (als *ArabicLightStemmer) GetSuffixLetters() string {
	return als.suffixLetters
}

// SetInfixLetters sets the infix letters used in the stemming process.
// Infix letters are characters or sequences of characters that may appear within the root of a word, not at the edges.
func (als *ArabicLightStemmer) SetInfixLetters(newInfixLetters string) {
	als.infixLetters = newInfixLetters
}

// GetInfixLetters returns the current infix letters used in the stemming process.
// These letters are used to identify and handle infixes within words during the stemming process.
func (als *ArabicLightStemmer) GetInfixLetters() string {
	return als.infixLetters
}

// SetJoker sets the joker character used in the stemming process.
// The joker character is typically used as a wildcard to represent any letter in certain stemming operations.
func (als *ArabicLightStemmer) SetJoker(newJoker string) {
	// Ensure that the joker character is only one character long.
	if len(newJoker) > 1 {
		newJoker = newJoker[:1]
	}
	als.joker = newJoker
}

// GetJoker returns the current joker character used in the stemming process.
// The joker is often used as a placeholder for any character in pattern matching and root extraction.
func (als *ArabicLightStemmer) GetJoker() string {
	return als.joker
}

// SetMaxPrefixLength sets the maximum length for prefixes during the stemming process.
// This value limits how long a prefix can be when identifying and removing prefixes from words.
func (als *ArabicLightStemmer) SetMaxPrefixLength(newMaxPrefixLength int) {
	als.maxPrefixLength = newMaxPrefixLength
}

// GetMaxPrefixLength returns the current maximum length for prefixes used in the stemming process.
// It defines the maximum number of characters that can be considered a prefix in words.
func (als *ArabicLightStemmer) GetMaxPrefixLength() int {
	return als.maxPrefixLength
}

// SetMaxSuffixLength sets the maximum length for suffixes during the stemming process.
// This value limits how long a suffix can be when identifying and removing suffixes from words.
func (als *ArabicLightStemmer) SetMaxSuffixLength(newMaxSuffixLength int) {
	als.maxSuffixLength = newMaxSuffixLength
}

// GetMaxSuffixLength returns the current maximum length for suffixes used in the stemming process.
// It defines the maximum number of characters that can be considered a suffix in words.
func (als *ArabicLightStemmer) GetMaxSuffixLength() int {
	return als.maxSuffixLength
}

// SetMinStemLength sets the minimum length for the stem after removing prefixes and suffixes.
// This value ensures that the resulting stem is not shorter than a certain length, which could lead to incorrect results.
func (als *ArabicLightStemmer) SetMinStemLength(newMinStemLength int) {
	als.minStemLength = newMinStemLength
}

// GetMinStemLength returns the current minimum length for the stem used in the stemming process.
// It ensures that the stemmed word maintains a certain minimum length for accuracy.
func (als *ArabicLightStemmer) GetMinStemLength() int {
	return als.minStemLength
}

// SetPrefixList sets the list of possible prefixes used during the stemming process.
// This list contains the specific prefixes that the stemmer will look for when processing words.
func (als *ArabicLightStemmer) SetPrefixList(newPrefixList []string) {
	als.prefixList = newPrefixList
	// Recreate the prefix tree based on the new prefix list.
	als.createPrefixTree()
}

// GetPrefixList returns the current list of prefixes used in the stemming process.
// The stemmer uses this list to identify and remove prefixes from words.
func (als *ArabicLightStemmer) GetPrefixList() []string {
	return als.prefixList
}

// SetSuffixList sets the list of possible suffixes used during the stemming process.
// This list contains the specific suffixes that the stemmer will look for when processing words.
func (als *ArabicLightStemmer) SetSuffixList(newSuffixList []string) {
	als.suffixList = newSuffixList
	// Recreate the suffix tree based on the new suffix list.
	als.createSuffixTree()
}

// GetSuffixList returns the current list of suffixes used in the stemming process.
// The stemmer uses this list to identify and remove suffixes from words.
func (als *ArabicLightStemmer) GetSuffixList() []string {
	return als.suffixList
}

// SetRootsList sets the list of known roots used during the stemming process.
// This list contains the valid roots that the stemmer will check against when processing words.
func (als *ArabicLightStemmer) SetRootsList(newRootsList []string) {
	als.rootList = newRootsList
}

// GetRootsList returns the current list of known roots used in the stemming process.
// The stemmer uses this list to verify whether a stem is a valid root.
func (als *ArabicLightStemmer) GetRootsList() []string {
	return als.rootList
}

// SetValidAffixesList sets the list of valid affixes (combinations of prefixes and suffixes) used during the stemming process.
// This list defines which combinations of affixes are considered valid when extracting stems.
func (als *ArabicLightStemmer) SetValidAffixesList(newValidAffixesList []string) {
	als.validAffixesList = newValidAffixesList
}

// GetValidAffixesList returns the current list of valid affixes used in the stemming process.
// The stemmer uses this list to ensure that the affix combinations applied to words are valid.
func (als *ArabicLightStemmer) GetValidAffixesList() []string {
	return als.validAffixesList
}

// createPrefixTree creates a prefix tree from the list of prefixes.
// It organizes prefixes into a tree structure to allow efficient prefix lookup during the stemming process.
func (als *ArabicLightStemmer) createPrefixTree() map[string]interface{} {
	prefixTree := make(map[string]interface{})
	for _, prefix := range als.prefixList {
		branch := prefixTree
		for _, char := range prefix {
			charStr := string(char)
			if _, exists := branch[charStr]; !exists {
				branch[charStr] = make(map[string]interface{})
			}
			branch = branch[charStr].(map[string]interface{})
		}
		if _, exists := branch["#"]; exists {
			branch["#"].(map[string]interface{})[prefix] = "#"
		} else {
			branch["#"] = map[string]interface{}{prefix: "#"}
		}
	}
	als.prefixesTree = prefixTree
	return prefixTree
}

// createSuffixTree creates a suffix tree from the list of suffixes.
// It organizes suffixes into a tree structure in reverse order to allow efficient suffix lookup during the stemming process.
func (als *ArabicLightStemmer) createSuffixTree() map[string]interface{} {
	suffixTree := make(map[string]interface{})
	for _, suffix := range als.suffixList {
		branch := suffixTree
		// Iterate over the suffix in reverse order
		for i := len(suffix) - 1; i >= 0; {
			r, size := utf8.DecodeLastRuneInString(suffix[:i+1])
			charStr := string(r)
			if _, exists := branch[charStr]; !exists {
				branch[charStr] = make(map[string]interface{})
			}
			branch = branch[charStr].(map[string]interface{})
			i -= size
		}
		if _, exists := branch["#"]; exists {
			branch["#"].(map[string]interface{})[suffix] = "#"
		} else {
			branch["#"] = map[string]interface{}{suffix: "#"}
		}
	}
	return suffixTree
}

// MostCommon returns the most common string from a list, prioritizing 3-letter roots.
// This method is used to select the most frequent root or stem when multiple options are available.
func (als *ArabicLightStemmer) mostCommon(lst []string) string {
	// Filter for three-letter roots
	var triRoots []string
	for _, item := range lst {
		if len(item) == 3 {
			triRoots = append(triRoots, item)
		}
	}

	// If there are three-letter roots, use them instead of the full list
	if len(triRoots) > 0 {
		lst = triRoots
	}

	// Create a map to count occurrences of each string
	counts := make(map[string]int)
	for _, item := range lst {
		counts[item]++
	}

	// Sort the list to ensure consistent order
	sort.Strings(lst)

	// Find the most common element
	var mostCommon string
	maxCount := 0
	for _, item := range lst {
		if counts[item] > maxCount {
			mostCommon = item
			maxCount = counts[item]
		}
	}

	return mostCommon
}

// IsRootLengthValid checks if the length of a root is valid, ensuring it is between 2 and 4 characters.
// This validation is important to filter out roots that are too short or too long.
func (als *ArabicLightStemmer) isRootLengthValid(root string) bool {
	length := len(root)
	return length >= 2 && length <= 4
}

// LightStem performs a light stemming operation on the given Arabic word and returns the stem.
// This method simplifies the word by removing affixes and reducing it to its core stem.
func (als *ArabicLightStemmer) LightStem(word string) string {
	if word == "" {
		return ""
	}
	_, unvocalized, stemLeft, stemRight := als.transform2Stars(word)
	segmentList, unvocalized, left, right := als.segment(word)
	return als.getStem(word, unvocalized, left, right, stemLeft, stemRight, -1, -1, segmentList)
}

// Transform2Stars transforms all non-affixation letters in a word into a star (joker character, default '*').
// It is used in the stemming process to identify the core components of a word by marking non-essential parts.
func (als *ArabicLightStemmer) transform2Stars(word string) (string, string, int, int) {
	word = als.wordProcessor.StripTashkeel(word)
	unvocalized := word
	word = strings.ReplaceAll(word, "آ", "أا")

	// Replace all non-prefix and non-suffix letters with joker
	nonAffixPattern := fmt.Sprintf("[^%s%s]", als.prefixLetters, als.suffixLetters)
	re := regexp.MustCompile(nonAffixPattern)
	word = re.ReplaceAllString(word, als.joker)

	// Convert word to rune slice for proper character indexing
	runeWord := []rune(word)
	jokerRune := []rune(als.joker)[0]

	// Find the left and right positions of the joker character
	left := -1
	right := -1
	for i, char := range runeWord {
		if char == jokerRune {
			if left == -1 {
				left = i
			}
			right = i
		}
	}

	if left >= 0 {
		left = min(left, als.maxPrefixLength-1)
		right = max(right+1, len(runeWord)-als.maxSuffixLength)

		// Original word segment and make all letters jokers except infixes
		prefix := string(runeWord[:left])
		stem := string([]rune(word)[left:right])
		suffix := string(runeWord[right:])

		prefix = regexp.MustCompile(fmt.Sprintf("[^%s]", als.prefixLetters)).ReplaceAllString(prefix, als.joker)

		if als.infixLetters != "" {
			stem = regexp.MustCompile(fmt.Sprintf("[^%s]", als.infixLetters)).ReplaceAllString(stem, als.joker)
		}
		suffix = regexp.MustCompile(fmt.Sprintf("[^%s]", als.suffixLetters)).ReplaceAllString(suffix, als.joker)
		word = prefix + stem + suffix
	}

	// Re-evaluate left and right positions after transformation
	runeWord = []rune(word)
	left = -1
	right = -1
	for i, char := range runeWord {
		if char == jokerRune {
			if left == -1 {
				left = i
			}
			right = i
		}
	}

	if left < 0 {
		left = min(als.maxPrefixLength, len(runeWord)-2)
	}
	if left >= 0 {
		prefix := string(runeWord[:left])
		for prefix != "" && !utils.Contains(als.prefixList, prefix) {
			prefix = string([]rune(prefix)[:len([]rune(prefix))-1])
		}
		if right < 0 {
			right = max(len([]rune(prefix)), len(runeWord)-als.maxSuffixLength)
		}
		suffix := string(runeWord[right:])

		for suffix != "" && !utils.Contains(als.suffixList, suffix) {
			suffix = string([]rune(suffix)[1:])
		}
		left = len([]rune(prefix))
		right = len(runeWord) - len([]rune(suffix))

		// Get the original word segment and make all letters jokers except infixes
		stem := string([]rune(word)[left:right])
		if als.infixLetters != "" {
			stem = regexp.MustCompile(fmt.Sprintf("[^%s]", als.infixLetters)).ReplaceAllString(stem, als.joker)
		}
		word = prefix + stem + suffix
	}

	// Store result
	//stemLeft := left
	//stemRight := right
	//starword := word

	return word, unvocalized, left, right
}

// Segment segments the given word by identifying prefix and suffix positions.
// It returns a map of segment indices, the unvocalized word, and the left and right positions of the stem.
func (als *ArabicLightStemmer) segment(word string) (map[int][][2]int, string, int, int) {
	unvocalized := als.wordProcessor.StripTashkeel(word)
	word = strings.ReplaceAll(word, constant.ALEF_MADDA, constant.HAMZA+constant.ALEF)

	var left, right int
	// Get all left positions of prefixes
	lefts := als.lookupPrefixes(word)
	// Get all right positions of suffixes
	rights := als.lookupSuffixes(word)

	if len(lefts) > 0 {
		left = utils.MaxFromSlice(lefts)
	} else {
		left = -1
	}

	if len(rights) > 0 {
		right = utils.MinFromSlice(rights)
	} else {
		right = -1
	}

	// Initialize the segment list without the entire word's segment
	segmentList := make(map[int][][2]int)

	// Track seen segments to avoid duplicates
	seenSegments := make(map[int]map[[2]int]struct{})

	// Helper function to check if a segment has been seen
	isSeen := func(left int, segment [2]int) bool {
		if _, ok := seenSegments[left]; !ok {
			seenSegments[left] = make(map[[2]int]struct{})
		}
		if _, exists := seenSegments[left][segment]; exists {
			return true
		}
		seenSegments[left][segment] = struct{}{}
		return false
	}

	// Add segmentation points based on prefix and suffix positions
	for _, i := range lefts {
		for _, j := range rights {
			if j >= i+2 {
				segment := [2]int{i, j}
				if !isSeen(i, segment) {
					segmentList[i] = append(segmentList[i], segment)
				}
			}
		}
	}

	// Filter segments according to valid affixes list
	left, right = als.getLeftRight(segmentList)

	return segmentList, unvocalized, left, right
}

// GetStem returns the stem of the word by slicing it based on identified prefix and suffix positions.
// This method ensures that the correct stem is extracted based on the segmented parts of the word.
func (als *ArabicLightStemmer) getStem(word, unvocalized string, left, right, stemLeft, stemRight, prefixIndex, suffixIndex int, segmentList map[int][][2]int) string {
	// Determine the left (prefix) index
	if prefixIndex >= 0 || suffixIndex >= 0 {
		if prefixIndex < 0 {
			left = stemLeft // Default left position
		} else {
			left = prefixIndex
		}

		// Determine the right (suffix) index
		if suffixIndex < 0 {
			right = stemRight // Default right position
		} else {
			right = suffixIndex
		}

		// Convert unvocalized string to rune slice for proper slicing
		unvocalizedRunes := []rune(unvocalized)

		// Ensure indices are within bounds of the rune slice
		if left < 0 {
			left = 0
		}
		if right > len(unvocalizedRunes) {
			right = len(unvocalizedRunes)
		}

		// Return the substring from unvocalized if indices are valid
		if left <= right && left < len(unvocalizedRunes) {
			return string(unvocalizedRunes[left:right])
		}
	}

	// Default case: return the chosen stem
	return als.chooseStem(word, unvocalized, left, right, stemLeft, stemRight, segmentList)
}

// ChooseStem selects the most appropriate stem from the word by evaluating possible segments.
// It checks for stopwords, validates affixes, and returns the best possible stem.
func (als *ArabicLightStemmer) chooseStem(word, unvocalized string, left, right, stemLeft, stemRight int, segmentList map[int][][2]int) string {
	// Check if the word is a stop word
	if als.stopWordManager.IsStopword(word) {
		return als.stopWordManager.StopStem(word)
	}

	// Segment the word if the segment list is empty
	if len(segmentList) == 0 {
		als.segment(word)
	}
	segList := segmentList

	validSegList := make(map[int][][2]int)
	for leftIndex, segments := range segList {
		for _, segment := range segments {
			rightIndex := segment[1]
			if als.verifyAffix(word, unvocalized, left, right, stemLeft, stemRight, leftIndex, rightIndex, segmentList) {
				validSegList[leftIndex] = append(validSegList[leftIndex], [2]int{leftIndex, rightIndex})
			}
		}
	}

	runeWord := []rune(word)
	runeUnvocalized := []rune(unvocalized)

	if len(validSegList) == 0 {
		// If no valid segments, use the entire word
		left = 0
		right = len(runeWord)
	} else {
		// Otherwise, choose the leftmost and rightmost valid segment
		left, right = als.getLeftRight(validSegList)
	}

	// Ensure left and right are within bounds
	if left < 0 {
		left = 0
	}
	if right > len(runeUnvocalized) {
		right = len(runeUnvocalized)
	}

	// Return the substring from unvocalized based on rune indexing
	return string(runeUnvocalized[left:right])
}

// VerifyAffix checks if the prefix and suffix combination (affix) is valid according to predefined rules.
// It validates the affix against known verb and noun rules to ensure correct stemming.
func (als *ArabicLightStemmer) verifyAffix(word, unvocalized string, left, right, stemLeft, stemRight int, prefixIndex, suffixIndex int, segmentList map[int][][2]int) bool {
	prefix := als.getPrefix(unvocalized, left, prefixIndex)
	suffix := als.getSuffix(unvocalized, right, suffixIndex)

	affix := prefix + "-" + suffix
	stem := als.getStem(word, unvocalized, left, right, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList)

	if utils.AffixInList(affix, constant.VERB_AFFIX_LIST) && als.validStem(stem, "verb", prefix) {
		if utils.AffixInList(affix, constant.NOUN_AFFIX_LIST) && als.validStem(stem, "noun", prefix) {
			return true // Valid as both a verb and a noun
		}
		return true // Valid as a verb
	}
	if utils.AffixInList(affix, constant.NOUN_AFFIX_LIST) && als.validStem(stem, "noun", prefix) {
		return true // Valid as a noun
	}
	return false // Not a valid verb or noun
}

// GetPrefix extracts and returns the prefix of the word based on the given left and prefix indices.
// This function helps in isolating the prefix part of a word for further processing.
func (als *ArabicLightStemmer) getPrefix(unvocalized string, left, prefixIndex int) string {
	unvocalizedRunes := []rune(unvocalized)

	if prefixIndex < 0 {
		if left >= 0 && left <= len(unvocalizedRunes) {
			return string(unvocalizedRunes[:left])
		}
		return string(unvocalizedRunes[:])
	} else {
		if prefixIndex >= 0 && prefixIndex <= len(unvocalizedRunes) {
			return string(unvocalizedRunes[:prefixIndex])
		}
		return string(unvocalizedRunes[:])
	}
}

// GetSuffix extracts and returns the suffix of the word based on the given right and suffix indices.
// This function helps in isolating the suffix part of a word for further processing.
func (als *ArabicLightStemmer) getSuffix(unvocalized string, right, suffixIndex int) string {
	unvocalizedRunes := []rune(unvocalized)

	if suffixIndex < 0 {
		if right >= 0 && right <= len(unvocalizedRunes) {
			return string(unvocalizedRunes[right:])
		}
		return ""
	} else {
		if suffixIndex >= 0 && suffixIndex <= len(unvocalizedRunes) {
			return string(unvocalizedRunes[suffixIndex:])
		}
		return ""
	}
}

// ValidStem checks if the extracted stem is valid based on the type of word (verb or noun) and the prefix.
// It applies specific rules to ensure that the stem follows Arabic language constraints.
func (als *ArabicLightStemmer) validStem(stem string, tag string, prefix string) bool {
	if stem == "" {
		return false
	}

	// Convert the stem and prefix to rune slices
	runeStem := []rune(stem)
	runePrefix := []rune(prefix)

	// Determine the length of the stem in runes
	stemLength := len(runeStem)

	switch tag {
	case "verb":
		// Verb has length <= 6
		if stemLength > 6 || stemLength < 2 {
			return false
		}
		// Forbidden letters in verbs like Teh Marbuta
		if strings.Contains(stem, constant.TEH_MARBUTA) {
			return false
		}
		// 6-letter stems must start with ALEF
		if stemLength == 6 && !strings.HasPrefix(stem, constant.ALEF) {
			return false
		}
		// 5-letter stems must start with ALEF/TEH
		if stemLength == 5 && !strings.HasPrefix(stem, constant.ALEF) && !strings.HasPrefix(stem, constant.TEH) {
			if strings.HasSuffix(prefix, constant.YEH) || strings.HasSuffix(prefix, constant.TEH) ||
				strings.HasSuffix(prefix, constant.NOON) || strings.HasSuffix(prefix, constant.ALEF_HAMZA_ABOVE) {
				return false
			}
		}
		// ALEF not allowed after certain prefix letters
		if strings.HasPrefix(stem, constant.ALEF) && len(runePrefix) > 0 && strings.ContainsAny(string(runePrefix[len(runePrefix)-1]), constant.YEH+constant.NOON+constant.TEH+constant.ALEF_HAMZA_ABOVE+constant.ALEF) {
			return false
		}
		// Lookup for verb stamp
		if !als.verbListManager.IsVerbStamp(stem) {
			return false
		}

	case "noun":
		// Noun length should be less than 8
		if stemLength >= 8 {
			return false
		}
	}

	return true
}

// GetAffixList generates a list of possible affix combinations (prefix and suffix) for the word.
// It uses segment indices to create tuples representing different combinations of prefixes and suffixes.
func (als *ArabicLightStemmer) getAffixList(word, unvocalized, root string, stemLeft, stemRight, prefixIndex, suffixIndex int, segmentList map[int][][2]int) []map[string]string {
	affixList := []map[string]string{}
	for leftIndex, segmentPairs := range segmentList {
		for _, pair := range segmentPairs {
			rightIndex := pair[1]
			affixTuple := als.getAffixTuple(word, unvocalized, root, leftIndex, rightIndex, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList)
			affixList = append(affixList, affixTuple)
		}
	}
	return affixList
}

// GetAffixTuple returns a dictionary representing a single affix tuple, including the prefix, suffix, stem, and root.
// It combines these elements to form a comprehensive affix representation.
func (als *ArabicLightStemmer) getAffixTuple(word, unvocalized, root string, left, right, stemLeft, stemRight, prefixIndex, suffixIndex int, segmentList map[int][][2]int) map[string]string {
	return map[string]string{
		"prefix":   als.getPrefix(unvocalized, left, prefixIndex),
		"suffix":   als.getSuffix(unvocalized, right, suffixIndex),
		"stem":     als.getStem(word, unvocalized, left, right, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList),
		"starstem": als.getStarStem(word, left, right, prefixIndex, suffixIndex),
		"root":     als.getRoot(word, unvocalized, root, left, right, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList),
	}
}

// GetRoot retrieves the root of the word by either extracting it from the stem or choosing from available options.
// This function handles the logic for determining the base root of the word after removing affixes.
func (als *ArabicLightStemmer) getRoot(word, unvocalized, root string, left, right, stemLeft, stemRight, prefixIndex, suffixIndex int, segmentList map[int][][2]int) string {
	if prefixIndex >= 0 || suffixIndex >= 0 {
		als.extractRoot(word, unvocalized, root, left, right, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList)
	} else {
		root = als.chooseRoot(word, unvocalized, root, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList)
	}
	return root
}

// ExtractRoot processes the word to extract its root by analyzing the stem and applying normalization techniques.
// This method is critical for isolating the root form of the word, which is used for further linguistic processing.
func (als *ArabicLightStemmer) extractRoot(word, unvocalized, root string, left, right, stemLeft, stemRight, prefixIndex, suffixIndex int, segmentList map[int][][2]int) string {
	stem := als.getStem(word, unvocalized, left, right, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList)

	// If the stem has 3 letters, it can be the root directly
	if len(stem) == 3 {
		root = als.ajustRoot(root, stem)
		return root
	}

	starStem := als.getStarStem(word, left, right, prefixIndex, suffixIndex)
	root = ""

	if len(starStem) == len(stem) {
		for i, char := range stem {
			if string(starStem[i]) == als.joker {
				root += string(char)
			}
		}
	} else {
		root = stem
	}

	// Normalize root
	root = als.normalizeRoot(root)

	// If the root length is 2, adjust the root
	if len(root) == 2 {
		root = als.ajustRoot(root, starStem)
	}

	return root
}

// ChooseRoot selects the best root from the possible roots extracted from the word.
// It applies length checks, dictionary validations, and frequency analysis to choose the most appropriate root.
func (als *ArabicLightStemmer) chooseRoot(word, unvocalized, root string, stemLeft, stemRight, prefixIndex, suffixIndex int, segmentList map[int][][2]int) string {
	if als.stopWordManager.IsStopword(word) {
		return als.stopWordManager.StopRoot(word)
	}

	if len(segmentList) == 0 {
		als.segment(word)
	}

	affixList := als.getAffixList(word, unvocalized, root, stemLeft, stemRight, prefixIndex, suffixIndex, segmentList)
	var roots []string
	for _, d := range affixList {
		roots = append(roots, d["root"])
	}

	// Filter roots by valid length
	var accepted []string
	for _, root := range roots {
		if als.isRootLengthValid(root) {
			accepted = append(accepted, root)
		}
	}
	if len(accepted) > 0 {
		roots = accepted
	}

	// Filter roots by checking if they are in the dictionary
	accepted = nil // Reset the accepted slice
	for _, root := range roots {
		if als.rootsManager.IsRoot(root) {
			accepted = append(accepted, root)
		}
	}
	if len(accepted) > 0 {
		roots = accepted
	}

	// Choose the most frequent root
	acceptedRoot := als.mostCommon(roots)

	return acceptedRoot
}

// AjustRoot modifies and refines the root based on specific patterns and linguistic rules.
// It adjusts the root, especially in cases where the standard root extraction process needs fine-tuning.
func (als *ArabicLightStemmer) ajustRoot(root, starstem string) string {
	if starstem == "" {
		return root
	}

	if len(starstem) == 3 {
		starstem = strings.ReplaceAll(starstem, constant.ALEF, constant.WAW)
		starstem = strings.ReplaceAll(starstem, constant.ALEF_MAKSURA, constant.YEH)
		return starstem
	}

	first := string(starstem[0])
	last := string(starstem[len(starstem)-1])

	switch {
	case first == constant.ALEF || first == constant.WAW:
		root = constant.WAW + root
	case first == constant.YEH:
		root = constant.YEH + root
	case first == als.joker && (last == constant.ALEF || last == constant.WAW):
		root += constant.WAW
	case first == als.joker && (last == constant.ALEF_MAKSURA || last == constant.YEH):
		root += constant.WAW
	case first == als.joker && last == als.joker:
		if len(starstem) == 2 {
			root += string(root[len(root)-1])
		} else {
			root = string(root[0]) + constant.WAW + string(root[1])
		}
	}

	return root
}

// NormalizeRoot standardizes the root by applying a series of replacements and adjustments.
// It ensures that the root conforms to expected linguistic norms in Arabic, such as handling specific characters.
func (als *ArabicLightStemmer) normalizeRoot(word string) string {
	// Replace ALEF_MADDA with HAMZA + ALEF
	word = strings.ReplaceAll(word, constant.ALEF_MADDA, constant.HAMZA+constant.ALEF)
	// Remove TEH_MARBUTA
	word = strings.ReplaceAll(word, constant.TEH_MARBUTA, "")
	// Replace ALEF_MAKSURA with YEH
	word = strings.ReplaceAll(word, constant.ALEF_MAKSURA, constant.YEH)
	// Normalize Hamza in the word
	return utils.NormalizeHamza(word)
}

// GetStarStem generates a "starred" version of the stem, where non-affix letters are replaced with a joker character.
// This method is used for pattern matching and helps in identifying the structure of the stem.
func (als *ArabicLightStemmer) getStarStem(word string, left, right int, prefixIndex, suffixIndex int) string {
	starword := word
	var tempLeft, tempRight int

	if prefixIndex < 0 && suffixIndex < 0 {
		tempLeft = left
		tempRight = right
	} else {
		tempLeft = left
		tempRight = right
		if prefixIndex >= 0 {
			tempLeft = prefixIndex
		}
		if suffixIndex >= 0 {
			tempRight = suffixIndex
		}
	}

	var newStarstem string
	if als.infixLetters != "" {
		// Convert all non-infix letters to the joker character
		infixPattern := fmt.Sprintf("[^%s%s]", als.infixLetters, constant.TEH_MARBUTA)
		newStarstem = regexp.MustCompile(infixPattern).ReplaceAllString(starword[tempLeft:tempRight], als.joker)
		// Handle specific infix cases
		newStarstem = als.handleTehInfix(word, newStarstem, tempLeft, tempRight)
	} else {
		// If there are no infix letters, convert all characters to jokers
		newStarstem = strings.Repeat(als.joker, len(starword[tempLeft:tempRight]))
	}

	return newStarstem
}

// HandleTehInfix applies special rules for handling the "Teh" infix and its variants within the stem.
// It ensures that certain infixes are correctly managed according to linguistic rules in Arabic.
func (als *ArabicLightStemmer) handleTehInfix(word, starword string, left, right int) string {
	newStarstem := starword

	// Case of Teh Marbuta
	keyStem := strings.ReplaceAll(newStarstem, constant.TEH_MARBUTA, "")
	if len(keyStem) != 4 {
		// Apply teh and variants only if the stem has 4 letters
		newStarstem = regexp.MustCompile(fmt.Sprintf("[%s%s%s]", constant.TEH, constant.TAH, constant.DAL)).ReplaceAllString(newStarstem, als.joker)
		return newStarstem
	}

	// Substitute teh in infixes, the teh must be in the first or second place, all others are converted
	newStarstem = newStarstem[:2] + strings.Replace(newStarstem[2:], constant.TEH, als.joker, -1)

	// Tah طاء is an infix if preceded by DHAD only
	if strings.HasPrefix(word[left:right], "ضط") {
		newStarstem = newStarstem[:2] + strings.Replace(newStarstem[2:], constant.TAH, als.joker, -1)
	} else {
		newStarstem = strings.ReplaceAll(newStarstem, constant.TAH, als.joker)
	}

	// DAL دال is an infix if preceded by ZAY only
	if strings.HasPrefix(word[left:right], "زد") {
		newStarstem = newStarstem[:2] + strings.Replace(newStarstem[2:], constant.DAL, als.joker, -1)
	} else {
		newStarstem = strings.ReplaceAll(newStarstem, constant.DAL, als.joker)
	}

	return newStarstem
}

// GetAffix returns a concatenated string of the prefix and suffix for the word, based on the provided indices.
// This method combines these elements into a single representation, useful for further processing.
func (als *ArabicLightStemmer) getAffix(unvocalized string, left int, right, prefixIndex, suffixIndex int) string {
	return strings.Join([]string{als.getPrefix(unvocalized, left, prefixIndex), als.getSuffix(unvocalized, right, suffixIndex)}, "-")
}

// GetLeftRight determines and returns the maximum left and minimum right values from a list of segments.
// This method helps in isolating the core segment of the word by narrowing down the possible prefixes and suffixes.
func (als *ArabicLightStemmer) getLeftRight(ls map[int][][2]int) (int, int) {
	if len(ls) == 0 {
		return -1, -1
	}

	// Find the maximum left position
	maxLeft := -1
	for left := range ls {
		if left > maxLeft {
			maxLeft = left
		}
	}

	// Find the minimum right position with the maximum left
	minRight := -1
	for _, segmentPairs := range ls {
		for _, pair := range segmentPairs {
			right := pair[1]
			if minRight == -1 || right < minRight {
				minRight = right
			}
		}
	}

	return maxLeft, minRight
}

// LookupPrefixes identifies and returns the positions of valid prefixes in the word by traversing the prefix tree.
// This method is used to locate the starting points of potential prefixes that can be removed from the word.
func (als *ArabicLightStemmer) lookupPrefixes(word string) []int {
	branch := als.prefixesTree
	lefts := []int{0}
	runeWord := []rune(word)
	i := 0

	for i < len(word) {
		char := string(runeWord[i])
		if _, ok := branch[char]; ok {
			if _, hasHash := branch["#"]; hasHash {
				lefts = append(lefts, i)
			}
			branch = branch[char].(map[string]interface{})
		} else {
			break
		}
		i++
	}

	if i < len(word) {
		if _, hasHash := branch["#"]; hasHash {
			lefts = append(lefts, i)
		}
	}

	return lefts
}

// LookupSuffixes identifies and returns the positions of valid suffixes in the word by traversing the suffix tree.
// This method is used to locate the ending points of potential suffixes that can be removed from the word.
func (als *ArabicLightStemmer) lookupSuffixes(word string) []int {
	branch := als.suffixesTree
	suffix := ""
	rights := []int{}
	runeWord := []rune(word)
	i := len(runeWord) - 1
	for i >= 0 {
		char := string(runeWord[i])
		if _, ok := branch[char]; ok {
			suffix = char + suffix
			if _, hasHash := branch["#"]; hasHash {
				rights = append(rights, i+1)
			}
			branch = branch[char].(map[string]interface{})
		} else {
			break
		}
		i--
	}

	if i >= 0 {
		if _, hasHash := branch["#"]; hasHash {
			rights = append(rights, i+1)
		}
	}

	return rights
}
