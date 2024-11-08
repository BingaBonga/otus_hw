package hw03frequencyanalysis

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

func Top10(input string) []string {
	splitText := strings.Fields(input)
	wordCountMap := make(map[string]int)

	for _, word := range splitText {
		if word == "" || word == "-" {
			continue
		}

		wordCountMap[onlyLetters(strings.ToLower(word))]++
	}

	worlds := make([]string, 0, len(wordCountMap))
	for word := range wordCountMap {
		worlds = append(worlds, word)
	}

	sort.SliceStable(worlds, func(i, j int) bool {
		if wordCountMap[worlds[i]] == wordCountMap[worlds[j]] {
			return worlds[i] < worlds[j]
		}

		return wordCountMap[worlds[i]] > wordCountMap[worlds[j]]
	})

	returnCount := int(math.Min(float64(len(worlds)), float64(10)))
	return worlds[:returnCount]
}

func onlyLetters(input string) string {
	runes := []rune(input)
	trimSuffixLength := 0
	trimPostfixLength := 0

	for i := 0; i < len(runes); i++ {
		if unicode.IsLetter(runes[i]) {
			break
		}

		trimSuffixLength++
	}

	for i := len(runes) - 1; i >= 0; i-- {
		if unicode.IsLetter(runes[i]) {
			break
		}

		trimPostfixLength++
	}

	return string(runes[trimSuffixLength : len(runes)-trimPostfixLength])
}
