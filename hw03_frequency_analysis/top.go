package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Pair struct {
	Word  string
	Count int
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) Less(i, j int) bool {
	first := p[i]
	second := p[j]

	switch {
	case first.Count == second.Count:
		return p.LexicographicalSort(first.Word, second.Word)
	default:
		return first.Count > second.Count
	}
}

func (p PairList) LexicographicalSort(firstString, secondString string) bool {
	switch strings.Compare(firstString, secondString) {
	case 1:
		return false
	default:
		return true
	}
}

func (p *PairList) Sort() {
	sort.Sort(*p)
}

func (p *PairList) transformMapToSlice(mapOfWords *map[string]int) {
	for word, count := range *mapOfWords {
		*p = append(*p, Pair{
			Word:  word,
			Count: count,
		})
	}
}

func (p *PairList) getTopWords(count int) []string {
	result := make([]string, 0)

	for i, item := range *p {
		if i == count {
			break
		}
		result = append(result, item.Word)
	}

	return result
}

func Top10(text string) []string {
	mapOfWords := make(map[string]int)

	sliceOfWords := strings.Fields(text)

	for _, word := range sliceOfWords {
		if val, ok := mapOfWords[word]; ok {
			val++
			mapOfWords[word] = val
		} else {
			mapOfWords[word] = 1
		}
	}

	pairList := new(PairList)

	pairList.transformMapToSlice(&mapOfWords)
	pairList.Sort()

	result := pairList.getTopWords(10)

	return result
}
