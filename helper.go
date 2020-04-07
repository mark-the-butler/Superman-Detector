package main

// CheckErr is a helping function for checking for erros
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
