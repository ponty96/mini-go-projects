package main

import (
	"bufio"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductSet(t *testing.T) {
	Convey("Given a product set", t, func() {
		set := NewProductSet()
		Convey("When a product is added to the set", func() {
			productID := 1
			productName := "Product 1"
			set.Add(productID, Product{Id: productID, Name: productName})
			Convey("Then the set should not be empty", func() {
				So(set.Empty(), ShouldBeFalse)
			})
			Convey("Then the set should contain the product", func() {
				product, ok := set.Contains(productID)
				So(ok, ShouldBeTrue)
				So(product, ShouldResemble, Product{Id: productID, Name: productName})
			})
			Convey("When the product is removed from the set", func() {
				delete(set.data, productID)
				Convey("Then the set should be empty", func() {
					So(set.Empty(), ShouldBeTrue)
				})
			})
		})

		Convey("When we compare the set with another set", func() {
			compareWithProductSet := NewProductSet()
			compareWithProductSet.Add(1, Product{Id: 1, Name: "Product 1"})
			compareWithProductSet.Add(2, Product{Id: 2, Name: "Product 2"})
			set.Add(1, Product{Id: 1, Name: "Product 1"})
			set.Add(3, Product{Id: 3, Name: "Product 3"})
			Convey("Then the intersection of the two sets should be a set containing the common products", func() {
				resultProductSet := set.Intersection(&compareWithProductSet)
				So(resultProductSet.Empty(), ShouldBeFalse)
				So(resultProductSet.data[1], ShouldResemble, Product{Id: 1, Name: "Product 1"})
			})

			Convey("Then the union of the two sets should be a set containing all products", func() {
				resultProductSet := set.Union(&compareWithProductSet)
				So(len(resultProductSet.data), ShouldEqual, 3)
				So(resultProductSet.Empty(), ShouldBeFalse)
				So(resultProductSet.data[1], ShouldResemble, Product{Id: 1, Name: "Product 1"})
				So(resultProductSet.data[2], ShouldResemble, Product{Id: 2, Name: "Product 2"})
				So(resultProductSet.data[3], ShouldResemble, Product{Id: 3, Name: "Product 3"})
			})
		})

		Convey("When we check if it is empty", func() {
			Convey("Then it should be empty", func() {
				So(set.Empty(), ShouldBeTrue)
			})
			Convey("It should not be empty if we add an element", func() {
				set.Add(1, Product{Id: 1, Name: "Product 1"})
				So(set.Empty(), ShouldBeFalse)
			})
		})

		Convey("When it check it's elements", func() {
			Convey("Then it should not contain any element", func() {
				_, ok := set.Contains(1)
				So(ok, ShouldBeFalse)
			})
			Convey("Then it should contain the element", func() {
				set.Add(1, Product{Id: 1, Name: "Product 1"})
				_, ok := set.Contains(1)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestSearchEngineDataProcessing(t *testing.T) {
	Convey("Given a file reader", t, func() {
		file, err := os.Open("/Users/ayomidearegbede/Documents/devland/go-stuff/simple-search/products_small.json")
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
				So(len(searchEngine.Index), ShouldEqual, 55)
			})
		})
	})
}
func TestSearchEngine(t *testing.T) {
	Convey("Given a search engine with data", t, func() {
		file, err := os.Open("/Users/ayomidearegbede/Documents/devland/go-stuff/simple-search/products_small.json")
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
				LeftTerm:  "Bronze",
				RightTerm: "Watch",
				Operation: "AND",
			})

			So(err, ShouldBeNil)
			Convey("it returns products with names that contains both words", func() {
				// So(len(result.data), ShouldEqual, 2)

				for _, product := range result.data {
					So(product.Name, ShouldContainSubstring, "Watch")
					So(product.Name, ShouldContainSubstring, "Bronze")
				}
			})
		})

		Convey("When an OR search query is made", func() {
			result, err := searchEngine.SearchQuery(Query{
				LeftTerm:  "Bronze",
				RightTerm: "Watch",
				Operation: "OR",
			})

			So(err, ShouldBeNil)
			Convey("it returns products with names that contains either of the words", func() {
				// So(len(result.data), ShouldEqual, 15)
				for _, product := range result.data {
					containsWord1 := strings.Contains(product.Name, "Bronze")
					containsWord2 := strings.Contains(product.Name, "Watch")

					So(containsWord1 || containsWord2, ShouldBeTrue)
				}
			})
		})
	})
}
