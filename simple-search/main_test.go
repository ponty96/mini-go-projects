package main

import (
	"bufio"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEventSet(t *testing.T) {
	Convey("Given a event set", t, func() {
		set := NewEventSet()
		Convey("When a event is added to the set", func() {
			eventID := 1
			eventTitle := "Event 1"
			set.Add(eventID, Event{Id: eventID, Title: eventTitle})
			Convey("Then the set should not be empty", func() {
				So(set.Empty(), ShouldBeFalse)
			})
			Convey("Then the set should contain the event", func() {
				event, ok := set.Contains(eventID)
				So(ok, ShouldBeTrue)
				So(event, ShouldResemble, Event{Id: eventID, Title: eventTitle})
			})
			Convey("When the event is removed from the set", func() {
				delete(set.data, eventID)
				Convey("Then the set should be empty", func() {
					So(set.Empty(), ShouldBeTrue)
				})
			})
		})

		Convey("When we compare the set with another set", func() {
			compareWithEventSet := NewEventSet()
			compareWithEventSet.Add(1, Event{Id: 1, Title: "Event 1"})
			compareWithEventSet.Add(2, Event{Id: 2, Title: "Event 2"})
			set.Add(1, Event{Id: 1, Title: "Event 1"})
			set.Add(3, Event{Id: 3, Title: "Event 3"})
			Convey("Then the intersection of the two sets should be a set containing the common events", func() {
				resultEventSet := set.Intersection(&compareWithEventSet)
				So(resultEventSet.Empty(), ShouldBeFalse)
				So(resultEventSet.data[1], ShouldResemble, Event{Id: 1, Title: "Event 1"})
			})

			Convey("Then the union of the two sets should be a set containing all events", func() {
				resultEventSet := set.Union(&compareWithEventSet)
				So(len(resultEventSet.data), ShouldEqual, 3)
				So(resultEventSet.Empty(), ShouldBeFalse)
				So(resultEventSet.data[1], ShouldResemble, Event{Id: 1, Title: "Event 1"})
				So(resultEventSet.data[2], ShouldResemble, Event{Id: 2, Title: "Event 2"})
				So(resultEventSet.data[3], ShouldResemble, Event{Id: 3, Title: "Event 3"})
			})
		})

		Convey("When we check if it is empty", func() {
			Convey("Then it should be empty", func() {
				So(set.Empty(), ShouldBeTrue)
			})
			Convey("It should not be empty if we add an element", func() {
				set.Add(1, Event{Id: 1, Title: "Event 1"})
				So(set.Empty(), ShouldBeFalse)
			})
		})

		Convey("When it check it's elements", func() {
			Convey("Then it should not contain any element", func() {
				_, ok := set.Contains(1)
				So(ok, ShouldBeFalse)
			})
			Convey("Then it should contain the element", func() {
				set.Add(1, Event{Id: 1, Title: "Event 1"})
				_, ok := set.Contains(1)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestSearchEngineDataProcessing(t *testing.T) {
	Convey("Given a file reader", t, func() {
		file, err := os.Open("./realistic_event_data.json")
		if err != nil {
			// fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		searchEngine := NewSearchEngine()

		fileReader := bufio.NewReader(file)

		searchEngine.PreprocessData(fileReader)

		Convey("When the data is preprocessed", func() {
			Convey("Then the search engine should contain the data", func() {
				So(len(searchEngine.Index), ShouldEqual, 21)
			})
		})
	})
}
func TestSearchEngine(t *testing.T) {
	Convey("Given a search engine with data", t, func() {
		file, err := os.Open("./realistic_event_data.json")
		if err != nil {
			// fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		searchEngine := NewSearchEngine()

		fileReader := bufio.NewReader(file)

		searchEngine.PreprocessData(fileReader)

		Convey("When an AND search query is made", func() {
			result, err := searchEngine.SearchQuery(Query{
				LeftTerm:  "Technology",
				RightTerm: "Webinar",
				Operation: "AND",
			})

			So(err, ShouldBeNil)
			Convey("it returns events with names that contains both words", func() {
				So(len(result.data), ShouldBeGreaterThan, 2)

				for _, event := range result.data {
					So(event.Title, ShouldContainSubstring, "Technology")
					So(event.Title, ShouldContainSubstring, "Webinar")
				}
			})
		})

		Convey("When an OR search query is made", func() {
			result, err := searchEngine.SearchQuery(Query{
				LeftTerm:  "Technology",
				RightTerm: "Science",
				Operation: "OR",
			})

			So(err, ShouldBeNil)
			Convey("it returns events with names that contains either of the words", func() {
				So(len(result.data), ShouldBeGreaterThan, 2)
				for _, event := range result.data {
					containsWord1 := strings.Contains(event.Title, "Technology")
					containsWord2 := strings.Contains(event.Title, "Science")

					So(containsWord1 || containsWord2, ShouldBeTrue)
				}
			})
		})
	})
}
