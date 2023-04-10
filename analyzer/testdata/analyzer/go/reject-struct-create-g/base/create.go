package create

type innerStruct struct {
	someField int
}

type outerStruct struct {
	innerStruct
}

// crateEmbedded tests replacing key-value pairs in struct
// initialization for embedded structs
func crateEmbedded() outerStruct {
	return outerStruct{innerStruct: innerStruct{someField: 42}}
}
