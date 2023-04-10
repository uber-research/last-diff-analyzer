package rename

import "fmt"

type testRenamed struct {
	intValueRenamed int
}

func rename(i int) int {
	// test local type renaming
	t := testRenamed{i}
	fmt.Println(t)

	// test local struct field renaming
	t = testRenamed{intValueRenamed: i}
	fmt.Println(t)

	// test package type renaming
	t2 := testARenamed{i}
	fmt.Println(t2)

	// test package struct field renaming
	ant2 := testC{fooRenamed: i}
	fmt.Println(ant2)
	return i
}

type someMap map[int]int

// The following test renaming parameter and an element of a
// constructor. This is a similar renaming to the one tested in
// reject-type-rename-f, but while here it is correct, in the reject
// test (with a struct instead initialized using key/value pairs), it
// is not.
func createMap(j int) map[int]int {
	return map[int]int{j: 42}
}

func createNamedMap(j int) someMap {
	return someMap{j: 42}
}

func createNamedArray(j int) someArray {
	// defined in different file
	return someArray{j}
}

func createNamedStruct(j int) testRenamed {
	return testRenamed{j}
}
