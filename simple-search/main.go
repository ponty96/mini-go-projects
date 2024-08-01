package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type EventSet struct {
	data map[int]Event
}

func NewEventSet() EventSet {
	return EventSet{data: make(map[int]Event)}
}

func (set *EventSet) Add(key int, event Event) {
	set.data[key] = event
}
func (set *EventSet) Contains(key int) (Event, bool) {
	event, ok := set.data[key]
	return event, ok
}

func (set *EventSet) Empty() bool {
	return len(set.data) == 0
}

func (set *EventSet) Intersection(compareWithEventSet *EventSet) EventSet {
	resultEventSet := NewEventSet()
	// delete(set.data, key)
	if len(set.data) < len(compareWithEventSet.data) {
		for key := range set.data {
			if event, ok := compareWithEventSet.Contains(key); ok {
				resultEventSet.Add(key, event)
			}
		}
	} else {
		for key := range compareWithEventSet.data {
			if event, ok := set.Contains(key); ok {
				resultEventSet.Add(key, event)
			}
		}
	}
	return resultEventSet
}

func (set *EventSet) Union(otherEventSet *EventSet) EventSet {
	resultEventSet := NewEventSet()
	for key, event := range set.data {
		resultEventSet.Add(key, event)
	}
	for key, event := range otherEventSet.data {
		resultEventSet.Add(key, event)
	}
	return resultEventSet
}

type SearchEngine struct {
	Index map[string]EventSet
}

type Event struct {
	Title string `json:"Title"`
	Id    int    `json:"ID"`
}

type Query struct {
	LeftTerm  string
	RightTerm string
	Operation string
}

func NewSearchEngine() SearchEngine {
	// consider preprocessing the file
	return SearchEngine{
		Index: make(map[string]EventSet),
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

	var events []Event
	json.Unmarshal(content, &events)

	for _, event := range events {
		// fmt.Println(event.Title, event.Id)
		s.addEvent(event)
	}
}

func (s *SearchEngine) SearchQuery(query Query) (EventSet, error) {
	leftTermResult := s.searchWord(query.LeftTerm)
	rightTermResult := s.searchWord(query.RightTerm)

	switch {
	case query.Operation == "OR":
		return leftTermResult.Union(&rightTermResult), nil
	case leftTermResult.Empty() || rightTermResult.Empty():
		return NewEventSet(), nil
	case query.Operation == "AND":
		return leftTermResult.Intersection(&rightTermResult), nil
	default:
		return NewEventSet(), errors.New(fmt.Sprintf("Wrong operation passed %s", query.Operation))
	}
}

func (s *SearchEngine) Search(words string) EventSet {
	var resultEventSet EventSet
	for _, word := range strings.Split(words, " ") {
		searchResult := s.searchWord(word)

		resultEventSet = resultEventSet.Union(&searchResult)
	}
	return resultEventSet
}

func (s *SearchEngine) searchWord(word string) EventSet {
	lowerCaseWord := strings.ToLower(word)
	if setOfEvents, ok := s.Index[lowerCaseWord]; ok {
		return setOfEvents
	}
	return NewEventSet()
}

func (s *SearchEngine) addEvent(event Event) {
	words := strings.Split(event.Title, " ")
	for _, word := range words {
		lowerCaseWord := strings.ToLower(word)
		// fmt.Sprintf("Word: %s\n", lowerCaseWord)
		// fmt.Println(lowerCaseWord)
		// we want to add the lowerCaseWord to the Index
		if _, ok := s.Index[lowerCaseWord]; !ok {
			// fmt.Println("Word does not exist in the index")
			setAtIndex := NewEventSet()
			setAtIndex.Add(event.Id, event)
			s.Index[lowerCaseWord] = setAtIndex
		} else {
			// the lowerCaseWord already exists in the index
			// we want to add the event to the set
			setAtIndex := s.Index[lowerCaseWord]
			setAtIndex.Add(event.Id, event)
			s.Index[lowerCaseWord] = setAtIndex
		}
	}
}

func main() {
	// Call the function
	file, err := os.Open("./realistic_event_data.json")
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file
	reader := bufio.NewReader(file)

	searchEngine := NewSearchEngine()

	searchEngine.PreprocessData(reader)
	fmt.Println(searchEngine.Index)

	result := searchEngine.Search("Linen")
	// we need to convert the result to a list of events

	// for key, event := range result.data {
	// 	fmt.Println(key, event)
	// 	fmt.Println(event.Title)
	// 	fmt.Println(event.Material)
	// }
	fmt.Println(result)
}
