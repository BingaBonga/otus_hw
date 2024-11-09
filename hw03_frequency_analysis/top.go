package hw03frequencyanalysis

import (
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

	if len(worlds) < 10 {
		return worlds
	}

	return worlds[:10]
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

	if (trimSuffixLength + trimPostfixLength) > len(runes) {
		return input
	}

	return string(runes[trimSuffixLength : len(runes)-trimPostfixLength])
}
