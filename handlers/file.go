package handlers

// Must is an error handler
func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
