package stop_words

import (
	"encoding/json"
	"log"
	"os"
)

type StopwordManager interface {
	IsStopword(word string) bool
	StopStem(word string) string
	StopRoot(word string) string
}

// stopwordManager manages stopwords.
type stopwordManager struct {
	stopwords map[string]map[string]string
	processor WordProcessor
}

// NewStopwordManager creates a new instance of StopwordManager with the provided WordProcessor.
// It initializes the stopwords map by loading stopwords from a JSON file. If the file cannot be loaded,
// the function logs a fatal error and terminates the program.
func NewStopwordManager(processor WordProcessor) StopwordManager {
	stopWordManager := stopwordManager{processor: processor, stopwords: make(map[string]map[string]string)}

	err := stopWordManager.loadStopwords("./arabic/stop_words/stopwords.json")
	if err != nil {
		log.Fatal(err)
	}

	return &stopWordManager
}

// IsStopword checks if the given word is in the stopwords list.
// It returns true if the word is a stopword, false otherwise.
func (sm *stopwordManager) IsStopword(word string) bool {
	_, exists := sm.stopwords[word]
	return exists
}

// StopStem returns the stem of the given word if it is in the stopwords list.
// The stem is stripped of Tashkeel characters before being returned.
func (sm *stopwordManager) StopStem(word string) string {
	stem := ""
	if stopWord, exists := sm.stopwords[word]; exists {
		stem = stopWord["stem"]
		stem = sm.processor.StripTashkeel(stem)
	}
	return stem
}

// StopRoot returns the root of the given word, which in this case is the same as the stem.
// It calls StopStem to retrieve the root.
func (sm *stopwordManager) StopRoot(word string) string {
	return sm.StopStem(word)
}

// loadStopwords loads the stopwords from a JSON file specified by the filename.
// It returns an error if the file cannot be read or the JSON cannot be unmarshaled.
func (sm *stopwordManager) loadStopwords(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &sm.stopwords)
}
