package main

import (
	"fmt"
	"iter"
	"sort"
	"strconv"
)

type CustomMapType map[string]int

func (cm CustomMapType) SortedKeys() iter.Seq2[string, int] {

	return func(yield func(string, int) bool) {

		sortedKeys := SortMapKeysByNumericalOrder(cm)

		for index := range sortedKeys {

			if !yield(sortedKeys[index], cm[sortedKeys[index]]) {
				return
			}
		}
	}
}

func SortMapKeysByNumericalOrder(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		num1, _ := strconv.Atoi(keys[i])
		num2, _ := strconv.Atoi(keys[j])
		return num1 < num2
	})

	return keys
}

func ProcessTotalWithIter(m CustomMapType) int {

	result := int(0)

	for _, v := range m.SortedKeys() {
		result += v
	}

	return result
}

func ProcessTotalWithSepDataStructure(m CustomMapType) int {

	result := int(0)

	sortedKeys := SortMapKeysByNumericalOrder(m)

	for _, v := range sortedKeys {
		result += m[v]
	}

	return result
}

func main() {

	cm := CustomMapType{
		"2":   500,
		"1":   1000,
		"3":   400,
		"6":   300,
		"0.1": 200,
	}

	fmt.Println("result #1: ", ProcessTotalWithIter(cm))
	fmt.Println("result #2: ", ProcessTotalWithSepDataStructure(cm))

}
