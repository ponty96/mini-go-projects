package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Product struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Set struct {
	data map[int]bool
}

func NewSet() Set {
	return Set{data: make(map[int]bool)}
}

func (set *Set) Add(key int) {
	set.data[key] = true
}
func (set *Set) Contains(key int) bool {
	_, ok := set.data[key]
	return ok
}

func (set *Set) Empty() bool {
	return len(set.data) == 0
}

func (set *Set) Intersection(compareWithSet *Set) Set {
	resultSet := NewSet()
	// delete(set.data, key)
	if len(set.data) > len(compareWithSet.data) {
		for key := range set.data {
			if compareWithSet.Contains(key) {
				resultSet.Add(key)
			}
		}
	} else {
		for key := range compareWithSet.data {
			if set.Contains(key) {
				resultSet.Add(key)
			}
		}
	}
	return resultSet
}

func (set *Set) Union(otherSet *Set) Set {
	resultSet := NewSet()
	for key := range set.data {
		resultSet.Add(key)
	}
	for key := range otherSet.data {
		resultSet.Add(key)
	}
	return resultSet
}

type SearchEngine struct {
	Index map[string]Set
}

func (s *SearchEngine) AddProduct(product Product) {
	words := strings.Split(product.Name, " ")
	for _, word := range words {
		// fmt.Sprintf("Word: %s\n", word)
		// fmt.Println(word)
		// we want to add the word to the Index
		if _, ok := s.Index[word]; !ok {
			// fmt.Println("Word does not exist in the index")
			setAtIndex := NewSet()
			setAtIndex.Add(product.Id)
			s.Index[word] = setAtIndex
		} else {
			// the word already exists in the index
			// we want to add the product to the set
			setAtIndex := s.Index[word]
			setAtIndex.Add(product.Id)
			s.Index[word] = setAtIndex
		}
	}
}

func (s *SearchEngine) Search(query string) Set {
	words := strings.Split(query, " ")
	fmt.Println(words)
	resultSet := NewSet()
	for _, word := range words {
		fmt.Println(word)
		if setOfProducts, ok := s.Index[word]; ok {
			// the word does not exist in the index
			// we can return an empty set
			if resultSet.Empty() {
				resultSet = setOfProducts
			} else {
				resultSet = resultSet.Intersection(&setOfProducts)
			}
		}
	}
	return resultSet
}

func main() {
	// Call the function
	file, err := os.Open("/Users/ayomidearegbede/Documents/devland/go-stuff/simple-search/products_small.json")
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file
	reader := bufio.NewReader(file)

	// Read the file line by line
	var content []byte
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				content = append(content, line...)
				break
			}
			// fmt.Println("Error reading file:", err)
			return
		}
		content = append(content, line...)
	}

	var products []Product
	json.Unmarshal(content, &products)

	searchEngine := SearchEngine{
		Index: make(map[string]Set),
	}

	for _, product := range products {
		// fmt.Println(product.Name, product.Id)
		searchEngine.AddProduct(product)
	}

	// fmt.Println(searchEngine.Index)

	result := searchEngine.Search("Linen")

	fmt.Println(result)
}
