package create

type someStruct struct {
	firstField, secondField int
	thirdField              string
}

// useNames tests struct construction equivalence where a name-less list is replaced with a list of key-value pairs
func useNames() someStruct {
	return someStruct{7, 42, "hello"}
}
