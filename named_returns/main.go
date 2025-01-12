package main

func FuncTypeA() (string, error) {
	return "Hello", nil
}

func FuncTypeB() (s string, err error) {
	s = "Hello"
	return
}
