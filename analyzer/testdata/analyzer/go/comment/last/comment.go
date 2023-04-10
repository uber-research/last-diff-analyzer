package comment

// Foo is a dummy function
func Foo() bool {
	return true
}

// SimpleComment test the simplest comment version
func SimpleComment(b bool) int {
	// simple comment
	return 42
}

// EndingComment tests a comment ending the line
func EndingComment(b bool) int {
	return 42 // ending comment
}

// LongComment tests a comment that would span many lines
func LongComment(b bool) int {
	/*
	   l
	   o
	   n
	   g

	   c
	   o
	   m
	   m
	   e
	   n
	   t
	*/
	return 42
}
