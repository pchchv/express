package optionext

type testStruct struct{}

func returnTypedNoneOption() Option[testStruct] {
	return None[testStruct]()
}

func returnTypedSomeOption() Option[testStruct] {
	return Some(testStruct{})
}

func returnTypedNoneOptionPtr() Option[*testStruct] {
	return None[*testStruct]()
}

func returnTypedSomeOptionPtr() Option[*testStruct] {
	return Some(new(testStruct))
}
