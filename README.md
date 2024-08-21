# Go Arabic Light Stemmer

Go Arabic Light Stemmer is a Go implementation of an Arabic light stemming algorithm. It is designed to provide efficient and accurate stemming for Arabic text, focusing on reducing words to their root form while handling the complexities of Arabic morphology. This project is inspired by and based on the [Tashaphyne](https://github.com/linuxscout/tashaphyne) Python library, with modifications and enhancements tailored for the Go programming language.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [How It Works](#how-it-works)
    - [Stemming Process](#stemming-process)
    - [Stopword Filtering](#stopword-filtering)
- [Customization](#customization)
    - [Adding Custom Stopwords](#adding-custom-stopwords)
    - [Extending the Affix List](#extending-the-affix-list)
- [Contributing](#contributing)
- [Acknowledgements](#acknowledgements)

## Features

- Light stemming for Arabic text.
- Optimized for performance using Go's concurrency model.

## Installation

```
go get github.com/berkayersoyy/go-arabic-light-stemmer
```

## Usage

```
package main

import (
    "fmt"
    "github.com/berkayersoyy/go-arabic-light-stemmer"
)

func main() {
    // Example usage
    text := "النص العربي هنا"
    stemmed := github.com/berkayersoyy/go-arabic-light-stemmer.LightStem(text)
    fmt.Println("Stemmed Text:", stemmed)
}
```

## How It Works
### Stemming Process
The Arabic Light Stemmer follows these basic steps:

- Normalization: Diacritical marks are removed, and letters are normalized.
- Affix Removal: The algorithm removes recognized prefixes, suffixes, and infixes.
- Root Extraction: The remaining part of the word is checked against known root patterns, and the root is extracted.
- Verb Stamping: The verb stamping mechanism checks the word against known verb forms and validates its structure based on predefined rules.

### Stopword Filtering
The stopword filtering function compares each word in the input text against a list of common Arabic stopwords and removes them from the output.

### Customization
#### Adding Custom Stopwords
To add custom stopwords, modify the stopwords.json file located in the /resources directory. Add new stopwords to the list in this file, and they will be automatically included in the filtering process.

#### Extending the Affix List
To extend or modify the list of prefixes, suffixes, or infixes, update the respective constants in the affix_constant.go file.

## Contributing
Contributions are welcome! If you find a bug or have an idea for an improvement, feel free to open an issue or submit a pull request.

## Acknowledgements
This project is inspired by and based on the [Tashaphyne](https://github.com/linuxscout/tashaphyne) Python library.
