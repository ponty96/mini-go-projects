package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ProductSet struct {
	data map[int]Product
}

func NewProductSet() ProductSet {
	return ProductSet{data: make(map[int]Product)}
}

func (set *ProductSet) Add(key int, product Product) {
	set.data[key] = product
}
func (set *ProductSet) Contains(key int) (Product, bool) {
	product, ok := set.data[key]
	return product, ok
}

func (set *ProductSet) Empty() bool {
	return len(set.data) == 0
}

func (set *ProductSet) Intersection(compareWithProductSet *ProductSet) ProductSet {
	resultProductSet := NewProductSet()
	// delete(set.data, key)
	if len(set.data) < len(compareWithProductSet.data) {
		for key := range set.data {
			if product, ok := compareWithProductSet.Contains(key); ok {
				resultProductSet.Add(key, product)
			}
		}
	} else {
		for key := range compareWithProductSet.data {
			if product, ok := set.Contains(key); ok {
				resultProductSet.Add(key, product)
			}
		}
	}
	return resultProductSet
}

func (set *ProductSet) Union(otherProductSet *ProductSet) ProductSet {
	resultProductSet := NewProductSet()
	for key, product := range set.data {
		resultProductSet.Add(key, product)
	}
	for key, product := range otherProductSet.data {
		resultProductSet.Add(key, product)
	}
	return resultProductSet
}

type SearchEngine struct {
	Index map[string]ProductSet
}

type Product struct {
	Name     string `json:"name"`
	Id       int    `json:"id"`
	Material string `json:"material"`
}

type Query struct {
	LeftTerm  string
	RightTerm string
	Operation string
}

func NewSearchEngine() SearchEngine {
	// consider preprocessing the file
	return SearchEngine{
		Index: make(map[string]ProductSet),
	}
}

func (s *SearchEngine) PreprocessData(dataStream *bufio.Reader) {
	// Read the file line by line
	var content []byte
	for {
		line, err := dataStream.ReadBytes('\n')
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

	for _, product := range products {
		// fmt.Println(product.Name, product.Id)
		s.addProduct(product)
	}
}

func (s *SearchEngine) SearchQuery(query Query) (ProductSet, error) {
	leftTermResult := s.searchWord(query.LeftTerm)
	rightTermResult := s.searchWord(query.RightTerm)

	switch {
	case query.Operation == "OR":
		return leftTermResult.Union(&rightTermResult), nil
	case leftTermResult.Empty() || rightTermResult.Empty():
		return NewProductSet(), nil
	case query.Operation == "AND":
		return leftTermResult.Intersection(&rightTermResult), nil
	default:
		return NewProductSet(), errors.New(fmt.Sprintf("Wrong operation passed %s", query.Operation))
	}
}

func (s *SearchEngine) Search(words string) ProductSet {
	var resultProductSet ProductSet
	for _, word := range strings.Split(words, " ") {
		searchResult := s.searchWord(word)

		resultProductSet = resultProductSet.Union(&searchResult)
	}
	return resultProductSet
}

func (s *SearchEngine) searchWord(word string) ProductSet {
	if setOfProducts, ok := s.Index[word]; ok {
		return setOfProducts
	}
	return NewProductSet()
}

func (s *SearchEngine) addProduct(product Product) {
	words := strings.Split(product.Name, " ")
	for _, word := range words {
		// fmt.Sprintf("Word: %s\n", word)
		// fmt.Println(word)
		// we want to add the word to the Index
		if _, ok := s.Index[word]; !ok {
			// fmt.Println("Word does not exist in the index")
			setAtIndex := NewProductSet()
			setAtIndex.Add(product.Id, product)
			s.Index[word] = setAtIndex
		} else {
			// the word already exists in the index
			// we want to add the product to the set
			setAtIndex := s.Index[word]
			setAtIndex.Add(product.Id, product)
			s.Index[word] = setAtIndex
		}
	}
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

	searchEngine := NewSearchEngine()

	searchEngine.PreprocessData(reader)
	// fmt.Println(searchEngine.Index)

	result := searchEngine.Search("Linen")
	// we need to convert the result to a list of products

	// for key, product := range result.data {
	// 	fmt.Println(key, product)
	// 	fmt.Println(product.Name)
	// 	fmt.Println(product.Material)
	// }
	fmt.Println(result)
}
