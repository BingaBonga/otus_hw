package hw03frequencyanalysis

import (
	"math"
	"sort"
	"strings"
)

func Top10(input string) []string {
	splitText := strings.Fields(input)
	wordCountMap := make(map[string]int)

	for _, word := range splitText {
		if word == "" {
			continue
		}

		wordCountMap[word]++
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
