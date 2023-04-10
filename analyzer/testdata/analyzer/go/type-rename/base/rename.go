package rename

import "fmt"

type test struct {
	intValue int
}

func rename(i int) int {
	// test local type renaming
	t := test{i}
	fmt.Println(t)

	// test local struct field renaming
	t = test{intValue: i}
	fmt.Println(t)

	// test package type renaming
	t2 := testA{i}
	fmt.Println(t2)

	// test package struct field renaming
	ant2 := testC{foo: i}
	fmt.Println(ant2)
	return i
}

type someMap map[int]int

// The following test renaming parameter and an element of a
// constructor. This is a similar renaming to the one tested in
// reject-type-rename-f, but while here it is correct, in the reject
// test (with a struct instead initialized using key/value pairs), it
// is not.
func createMap(i int) map[int]int {
	return map[int]int{i: 42}
}

func createNamedMap(i int) someMap {
	return someMap{i: 42}
}

func createNamedArray(i int) someArray {
	// defined in different file
	return someArray{i}
}

func createNamedStruct(i int) test {
	return test{i}
}
