package hw03frequencyanalysis

import (
	"container/heap"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`[\[\],.!?:;'"&$()]`)

func splitText(text string) []string {
	cleanText := reg.ReplaceAllString(text, " ")
	return strings.Fields(cleanText)
}

func calcWordsFreq(data []string) map[string]int {
	result := make(map[string]int)

	for _, word := range data {
		if word == "-" {
			continue
		}
		wordLower := strings.ToLower(word)
		result[wordLower]++
	}

	return result
}

type itemQ struct {
	word  string
	count int
	index int
}

type priorityQueue []*itemQ

func (pq *priorityQueue) Len() int {
	return len(*pq)
}

func (pq *priorityQueue) Less(i, j int) bool {
	pqV := *pq
	if pqV[i].count == pqV[j].count {
		return pqV[i].word < pqV[j].word
	}
	return pqV[i].count > pqV[j].count
}

func (pq *priorityQueue) Swap(i, j int) {
	pqV := *pq
	pqV[i], pqV[j] = pqV[j], pqV[i]
	pqV[i].index = i
	pqV[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*itemQ)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}

func top(s string, t int) []string {
	splittedText := splitText(s)
	counterMap := calcWordsFreq(splittedText)

	pq := make(priorityQueue, 0, len(counterMap))
	index := 0
	for k, v := range counterMap {
		heap.Push(&pq, &itemQ{word: k, count: v, index: index})
		index++
	}

	if len(counterMap) <= t {
		t = len(counterMap)
	}

	res := make([]string, 0, t)
	for t > 0 {
		item := heap.Pop(&pq).(*itemQ)
		res = append(res, item.word)
		t--
	}
	return res
}

func Top10(s string) []string {
	return top(s, 10)
}
