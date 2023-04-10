package create

type someStruct struct {
	firstField, secondField int
	thirdField              string
}

// useNames tests struct construction non-equivalence where a name-less list is replaced with a list of key-value pairs
func useNames() someStruct {
	return someStruct{firstField: 7, secondField: 42, thirdField: "hello"}
}
