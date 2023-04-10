package path

type S1 struct {
	a        int
	aRenamed int
}

// The local variable "a" is renamed to aRenamed, which happens to be the same as a field declaration of S1, in the last
// diff, the "a" in the access path is also changed to "aRenamed", but since they refer to different declarations, the
// "b.a = 1" line actually has different meaning than "b.aRenamed = 1". Hence, the approver should reject this example.
func test() {
	var aRenamed int
	var b S1
	// b refers to the line above, and aRenamed refers to the field declaration of S1
	b.aRenamed = 1
	// this a refers to the variable declaration at the begining
	aRenamed = 2
}
